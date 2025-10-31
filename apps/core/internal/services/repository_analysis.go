package services

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/dopeCape/kova/internal/models"
	"github.com/google/uuid"
)

type RepositoryAnalyzerService struct {
	tempDir       string
	githubService *GitHubService
}

func NewRepositoryAnalyzerService(githubService *GitHubService) *RepositoryAnalyzerService {
	// Use system temp directory
	tempDir := filepath.Join(os.TempDir(), "kova-analyzer")
	return &RepositoryAnalyzerService{
		tempDir:       tempDir,
		githubService: githubService,
	}
}

var skippableCommands = []string{"mkdir -p /app/node_modules/.cache"}

// AnalyzeRepository clones a repository and analyzes it using railpack
func (s *RepositoryAnalyzerService) AnalyzeRepository(ctx context.Context, accessToken string, req *models.AnalyzeRepositoryRequest) (*models.RepositoryAnalysis, error) {

	analysisID := uuid.New().String()
	cloneDir := filepath.Join(s.tempDir, analysisID)

	defer s.githubService.cleanup(cloneDir)

	// Create temp directory
	if err := os.MkdirAll(cloneDir, 0755); err != nil {
		return nil, fmt.Errorf("failed to create temp directory: %w", err)
	}

	// Clone repository
	if err := s.githubService.cloneRepository(ctx, accessToken, GithubCloneRequest(*req), cloneDir); err != nil {
		return nil, fmt.Errorf("failed to clone repository: %w", err)
	}

	// Run railpack analysis
	railpackOutput, err := s.runRailpack(ctx, cloneDir)
	if err != nil {
		return nil, fmt.Errorf("failed to run railpack: %w", err)
	}

	if !s.checkIfSupported(railpackOutput) {
		return &models.RepositoryAnalysis{
			Success: false,
			Install: []string{},
			Build:   []string{},
			Deploy:  "",
		}, nil
	}

	// Parse commands
	analysis := s.parseCommands(railpackOutput)
	return analysis, nil
}

// runRailpack executes railpack info command and returns parsed output
func (s *RepositoryAnalyzerService) runRailpack(ctx context.Context, repoDir string) (*models.RailpackOutput, error) {
	cmd := exec.CommandContext(ctx, "railpack", "info", ".", "--format", "json")
	cmd.Dir = repoDir

	output, err := cmd.Output()
	if err != nil {
		return nil, fmt.Errorf("railpack command failed: %w", err)
	}

	var parsed models.RailpackOutput
	if err := json.Unmarshal(output, &parsed); err != nil {
		return nil, fmt.Errorf("failed to parse railpack output: %w", err)
	}

	return &parsed, nil
}

// checkIfSupported checks if railpack could determine how to build the app
func (s *RepositoryAnalyzerService) checkIfSupported(output *models.RailpackOutput) bool {
	if len(output.Logs) > 0 {
		if strings.Contains(output.Logs[0].Msg, "Railpack could not determine how to build the app") {
			return false
		}
	}
	return true
}

// parseCommands extracts install, build, and deploy commands from railpack output
func (s *RepositoryAnalyzerService) parseCommands(output *models.RailpackOutput) *models.RepositoryAnalysis {
	analysis := &models.RepositoryAnalysis{
		Install: make([]string, 0),
		Build:   make([]string, 0),
		Deploy:  "",
		Success: true,
	}

	for _, step := range output.Plan.Steps {
		if step.Name == "install" || step.Name == "build" {
			for _, cmd := range step.Commands {
				var command string

				if cmd.CustomName != "" {
					command = cmd.CustomName
				} else if cmd.Cmd != "" && !s.isSkippableCommand(cmd.Cmd) {
					command = cmd.Cmd
				}

				if command != "" {
					if step.Name == "install" {
						analysis.Install = append(analysis.Install, command)
					} else if step.Name == "build" {
						analysis.Build = append(analysis.Build, command)
					}
				}
			}
		}
	}

	analysis.Deploy = output.Plan.Deploy.StartCommand
	if analysis.Deploy == "" {
		analysis.Success = false
	}
	return analysis
}

// isSkippableCommand checks if a command should be skipped
func (s *RepositoryAnalyzerService) isSkippableCommand(cmd string) bool {
	for _, skippable := range skippableCommands {
		if cmd == skippable {
			return true
		}
	}
	return false
}
