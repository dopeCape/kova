package services

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"os/exec"
	"time"
)

type GitHubUser struct {
	ID        int64  `json:"id"`
	Login     string `json:"login"`
	AvatarURL string `json:"avatar_url"`
	Name      string `json:"name"`
	Email     string `json:"email"`
	Type      string `json:"type"`
}

type GitHubRepository struct {
	ID            int64  `json:"id"`
	Name          string `json:"name"`
	FullName      string `json:"full_name"`
	Private       bool   `json:"private"`
	Description   string `json:"description"`
	Language      string `json:"language"`
	StarCount     int    `json:"stargazers_count"`
	UpdatedAt     string `json:"updated_at"`
	DefaultBranch string `json:"default_branch"`
	HTMLURL       string `json:"html_url"`
}

type GithubCloneRequest struct {
	RepoURL   string `json:"repo_url" validate:"required,url"`
	Branch    string `json:"branch" validate:"omitempty"`
	RepoID    int64  `json:"repo_id" validate:"required"`
	RepoName  string `json:"repo_name" validate:"required"`
	RepoOwner string `json:"repo_owner" validate:"required"`
}

type GitHubService struct {
	client  *http.Client
	baseURL string
}

func NewGitHubService() *GitHubService {
	return &GitHubService{
		client: &http.Client{
			Timeout: 30 * time.Second,
		},
		baseURL: "https://api.github.com",
	}
}

// GetUser retrieves GitHub user information using an access token
func (s *GitHubService) GetUser(ctx context.Context, accessToken string) (*GitHubUser, error) {
	if accessToken == "" {
		return nil, fmt.Errorf("access token is required")
	}

	req, err := http.NewRequestWithContext(ctx, "GET", s.baseURL+"/user", nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	// Set authorization header
	req.Header.Set("Authorization", "Bearer "+accessToken)
	req.Header.Set("Accept", "application/vnd.github+json")
	req.Header.Set("X-GitHub-Api-Version", "2022-11-28")
	req.Header.Set("User-Agent", "Kova-App")

	resp, err := s.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to make request: %w", err)
	}
	defer resp.Body.Close()

	// Handle different HTTP status codes
	switch resp.StatusCode {
	case http.StatusOK:
		// Success, continue to parse response
	case http.StatusUnauthorized:
		return nil, fmt.Errorf("invalid or expired access token")
	case http.StatusForbidden:
		return nil, fmt.Errorf("access token lacks required permissions")
	case http.StatusNotFound:
		return nil, fmt.Errorf("user not found")
	default:
		return nil, fmt.Errorf("github API returned status %d", resp.StatusCode)
	}

	var githubUser GitHubUser
	if err := json.NewDecoder(resp.Body).Decode(&githubUser); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	// Validate required fields
	if githubUser.ID == 0 {
		return nil, fmt.Errorf("invalid GitHub user: missing ID")
	}
	if githubUser.Login == "" {
		return nil, fmt.Errorf("invalid GitHub user: missing login")
	}

	return &githubUser, nil
}

// ValidateToken checks if the provided access token is valid
func (s *GitHubService) ValidateToken(ctx context.Context, accessToken string) error {
	_, err := s.GetUser(ctx, accessToken)
	return err
}

// GetUserByUsername retrieves public GitHub user information by username (no token required)
func (s *GitHubService) GetUserByUsername(ctx context.Context, username string) (*GitHubUser, error) {
	if username == "" {
		return nil, fmt.Errorf("username is required")
	}

	req, err := http.NewRequestWithContext(ctx, "GET", s.baseURL+"/users/"+username, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Accept", "application/vnd.github+json")
	req.Header.Set("X-GitHub-Api-Version", "2022-11-28")
	req.Header.Set("User-Agent", "Kova-App")

	resp, err := s.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to make request: %w", err)
	}
	defer resp.Body.Close()

	switch resp.StatusCode {
	case http.StatusOK:
		// Success, continue to parse response
	case http.StatusNotFound:
		return nil, fmt.Errorf("user not found")
	default:
		return nil, fmt.Errorf("github API returned status %d", resp.StatusCode)
	}

	var githubUser GitHubUser
	if err := json.NewDecoder(resp.Body).Decode(&githubUser); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	return &githubUser, nil
}

func (s *GitHubService) GetRepositories(ctx context.Context, accessToken string, page, perPage int) ([]*GitHubRepository, error) {
	if accessToken == "" {
		return nil, fmt.Errorf("access token is required")
	}

	// Default pagination values
	if page < 1 {
		page = 1
	}
	if perPage < 1 || perPage > 100 {
		perPage = 30
	}

	url := fmt.Sprintf("%s/user/repos?page=%d&per_page=%d&sort=updated&direction=desc", s.baseURL, page, perPage)
	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	// Set authorization header
	req.Header.Set("Authorization", "Bearer "+accessToken)
	req.Header.Set("Accept", "application/vnd.github+json")
	req.Header.Set("X-GitHub-Api-Version", "2022-11-28")
	req.Header.Set("User-Agent", "Kova-App")

	resp, err := s.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to make request: %w", err)
	}
	defer resp.Body.Close()

	// Handle different HTTP status codes
	switch resp.StatusCode {
	case http.StatusOK:
		// Success, continue to parse response
	case http.StatusUnauthorized:
		return nil, fmt.Errorf("invalid or expired access token")
	case http.StatusForbidden:
		return nil, fmt.Errorf("access token lacks required permissions")
	default:
		return nil, fmt.Errorf("github API returned status %d", resp.StatusCode)
	}

	var repositories []*GitHubRepository
	if err := json.NewDecoder(resp.Body).Decode(&repositories); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	return repositories, nil
}

func (s *GitHubService) cloneRepository(ctx context.Context, accessToken string, req GithubCloneRequest, destDir string) error {
	// Build clone URL with token
	// Format: https://<token>@github.com/<owner>/<repo>.git
	cloneURL := fmt.Sprintf("https://%s@github.com/%s/%s.git", accessToken, req.RepoOwner, req.RepoName)

	// Prepare git clone command with shallow clone for performance
	args := []string{"clone", "--depth", "1"}

	// Add branch if specified
	if req.Branch != "" {
		args = append(args, "--branch", req.Branch)
	}

	args = append(args, cloneURL, destDir)

	cmd := exec.CommandContext(ctx, "git", args...)

	// Capture output for debugging
	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("git clone failed: %w, output: %s", err, string(output))
	}

	return nil
}

func (s *GitHubService) cleanup(dir string) {
	if err := os.RemoveAll(dir); err != nil {
		// Log error but don't fail the request
		fmt.Printf("Warning: failed to cleanup directory %s: %v\n", dir, err)
	}
}
