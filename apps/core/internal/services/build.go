package services

import (
	"bufio"
	"context"
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"sync"
	"text/template"

	"github.com/dopeCape/kova/internal/models"
	"github.com/dopeCape/kova/internal/store"
)

const (
	BUILDKIT_ENV_MISSING       = "BUILDKIT_HOST environment variable is not set."
	BUILDKIT_UNABLE_TO_CONNECT = "ERRO failed to get buildkit information."
	REPO_BASE_PATH             = "/data/kova/repo"
	SERVICES_BASE_PATH         = "/data/kova/services"
	NETWORK_NAME               = "proxy"
)

type BuildJob struct {
	ProjectID string
	UserID    string
}

type BuildService struct {
	store        store.Store
	accountStore store.AccountStore
	queue        chan BuildJob
	wg           sync.WaitGroup
	wsHub        *WebSocketHub
	ctx          context.Context
	cancel       context.CancelFunc
}

func NewBuildService(store store.Store, accountStore store.AccountStore, wsHub *WebSocketHub) *BuildService {
	ctx, cancel := context.WithCancel(context.Background())
	bs := &BuildService{
		store:        store,
		accountStore: accountStore,
		queue:        make(chan BuildJob, 100), // Buffer of 100 jobs
		wsHub:        wsHub,
		ctx:          ctx,
		cancel:       cancel,
	}

	// Start single worker
	bs.wg.Add(1)
	go bs.worker()

	log.Println("✅ Build service initialized with 1 worker")
	return bs
}

func (bs *BuildService) Enqueue(projectID, userID string) {
	bs.queue <- BuildJob{
		ProjectID: projectID,
		UserID:    userID,
	}
	log.Printf("📦 Build job enqueued for project: %s", projectID)
}

func (bs *BuildService) worker() {
	defer bs.wg.Done()

	for {
		select {
		case <-bs.ctx.Done():
			log.Println("🛑 Build worker shutting down...")
			return
		case job := <-bs.queue:
			log.Printf("🔨 Processing build job for project: %s", job.ProjectID)
			if err := bs.processBuild(job); err != nil {
				log.Printf("❌ Build failed for project %s: %v", job.ProjectID, err)
				// Update to failed status
				bs.updateDeploymentStatus(job.ProjectID, "failed")
				bs.broadcastStatus(job.ProjectID, "failed")
			}
		}
	}
}

func (bs *BuildService) processBuild(job BuildJob) error {
	ctx := context.Background()

	log.Printf("🔨 ============================================")
	log.Printf("🔨 Starting build process for project: %s", job.ProjectID)
	log.Printf("🔨 User ID: %s", job.UserID)
	log.Printf("🔨 ============================================")

	// Get project
	log.Printf("🔨 [1/8] Fetching project details...")
	project, err := bs.store.GetProjectByID(ctx, job.ProjectID)
	if err != nil {
		log.Printf("❌ Failed to get project: %v", err)
		return fmt.Errorf("failed to get project: %w", err)
	}
	log.Printf("✅ Project fetched: %s (Repo: %s, Branch: %s)", project.Name, project.RepoFullName, project.RepoBranch)

	// Get user's first account for GitHub token
	log.Printf("🔨 [2/8] Fetching user accounts...")
	accounts, err := bs.accountStore.GetAccountsByUserID(ctx, job.UserID)
	if err != nil || len(accounts) == 0 {
		log.Printf("❌ Failed to get user accounts: %v (count: %d)", err, len(accounts))
		return fmt.Errorf("failed to get user account: %w", err)
	}
	log.Printf("✅ Found %d account(s) for user", len(accounts))

	// Get account with token
	log.Printf("🔨 [3/8] Fetching account tokens...")
	accountsWithTokens, err := bs.accountStore.GetAccountsByUserIDWithTokens(ctx, job.UserID)
	if err != nil || len(accountsWithTokens) == 0 {
		log.Printf("❌ Failed to get account tokens: %v (count: %d)", err, len(accountsWithTokens))
		return fmt.Errorf("failed to get account tokens: %w", err)
	}
	token := accountsWithTokens[0].AccessToken
	log.Printf("✅ Retrieved access token for account: %s", accountsWithTokens[0].GithubUsername)

	// Stage 1: Clone repository
	log.Printf("🔨 [4/8] Starting repository clone stage...")
	bs.updateDeploymentStatus(job.ProjectID, "building")
	bs.broadcastStatus(job.ProjectID, "building")
	log.Printf("📡 Status updated to: building")

	repoPath := filepath.Join(REPO_BASE_PATH, job.ProjectID)
	log.Printf("🔨 Repository will be cloned to: %s", repoPath)

	if err := bs.cloneRepository(project, token, repoPath); err != nil {
		log.Printf("❌ Clone failed: %v", err)
		bs.cleanup(job.ProjectID)
		return fmt.Errorf("clone failed: %w", err)
	}
	log.Printf("✅ Repository cloned successfully")

	// Stage 2: Build with railpack
	log.Printf("🔨 [5/8] Starting railpack build stage...")
	if err := bs.buildWithRailpack(project, repoPath); err != nil {
		log.Printf("❌ Build failed: %v", err)
		bs.cleanup(job.ProjectID)
		return fmt.Errorf("build failed: %w", err)
	}
	log.Printf("✅ Railpack build completed successfully")

	// Stage 3: Generate docker-compose and deploy
	log.Printf("🔨 [6/8] Starting deployment preparation...")
	bs.updateDeploymentStatus(job.ProjectID, "deploying")
	bs.broadcastStatus(job.ProjectID, "deploying")
	log.Printf("📡 Status updated to: deploying")

	if err := bs.generateDockerCompose(project); err != nil {
		log.Printf("❌ Docker-compose generation failed: %v", err)
		bs.cleanup(job.ProjectID)
		return fmt.Errorf("docker-compose generation failed: %w", err)
	}
	log.Printf("✅ Docker-compose file generated successfully")

	// Stage 4: Deploy with docker swarm
	log.Printf("🔨 [7/8] Deploying to Docker Swarm...")
	if err := bs.deployWithSwarm(project); err != nil {
		log.Printf("❌ Deployment failed: %v", err)
		bs.cleanup(job.ProjectID)
		return fmt.Errorf("deployment failed: %w", err)
	}
	log.Printf("✅ Deployed to Docker Swarm successfully")

	// Success!
	log.Printf("🔨 [8/8] Finalizing deployment...")
	bs.updateDeploymentStatus(job.ProjectID, "deployed")
	bs.broadcastStatus(job.ProjectID, "deployed")
	log.Printf("📡 Status updated to: deployed")

	log.Printf("🎉 ============================================")
	log.Printf("🎉 Build completed successfully for project: %s", job.ProjectID)
	log.Printf("🎉 Domain: %s", project.Domain)
	log.Printf("🎉 ============================================")

	return nil
}

func (bs *BuildService) cloneRepository(project *models.Project, token, repoPath string) error {
	log.Printf("📥 ============================================")
	log.Printf("📥 Cloning repository: %s", project.RepoURL)
	log.Printf("📥 Target path: %s", repoPath)
	log.Printf("📥 Branch: %s", project.RepoBranch)
	log.Printf("📥 ============================================")

	// Create directory
	log.Printf("📥 Creating directory: %s", repoPath)
	if err := os.MkdirAll(repoPath, 0755); err != nil {
		log.Printf("❌ Failed to create directory: %v", err)
		return fmt.Errorf("failed to create repo directory: %w", err)
	}
	log.Printf("✅ Directory created successfully")

	// Build authenticated URL
	authURL := strings.Replace(project.RepoURL, "https://", fmt.Sprintf("https://%s@", token), 1)
	log.Printf("📥 Constructed authenticated URL (token hidden)")

	// Clone with specific branch
	log.Printf("📥 Executing: git clone -b %s --single-branch [REPO_URL] %s", project.RepoBranch, repoPath)
	cmd := exec.Command("git", "clone", "-b", project.RepoBranch, "--single-branch", authURL, repoPath)
	cmd.Env = append(os.Environ(), "GIT_TERMINAL_PROMPT=0") // Disable interactive prompts

	output, err := cmd.CombinedOutput()
	if err != nil {
		log.Printf("❌ Git clone failed")
		log.Printf("❌ Error: %v", err)
		log.Printf("❌ Output: %s", string(output))
		return fmt.Errorf("git clone failed: %w, output: %s", err, string(output))
	}

	log.Printf("✅ Git clone output: %s", string(output))
	log.Printf("✅ Repository cloned successfully")
	return nil
}

func (bs *BuildService) buildWithRailpack(project *models.Project, repoPath string) error {
	log.Printf("🏗️  ============================================")
	log.Printf("🏗️  Building with railpack")
	log.Printf("🏗️  Project ID: %s", project.ID)
	log.Printf("🏗️  Working directory: %s", repoPath)
	log.Printf("🏗️  Environment variables: %d", len(project.EnvVariables))
	log.Printf("🏗️  ============================================")

	// Build env flags
	envFlags := []string{"build", "."}
	envFlags = append(envFlags, "--name", project.ID)
	log.Printf("🏗️  Image name: %s:latest", project.ID)

	if len(project.EnvVariables) > 0 {
		log.Printf("🏗️  Adding environment variables:")
		for _, env := range project.EnvVariables {
			if env.Key != "" && env.Value != "" {
				log.Printf("🏗️    - %s=%s", env.Key, strings.Repeat("*", min(len(env.Value), 8)))
				envFlags = append(envFlags, "--env", fmt.Sprintf("%s=%s", env.Key, env.Value))
			}
		}
	} else {
		log.Printf("🏗️  No environment variables to add")
	}

	log.Printf("🏗️  Executing: railpack %v", envFlags)
	cmd := exec.Command("railpack", envFlags...)
	cmd.Dir = repoPath

	// Capture output
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		log.Printf("❌ Failed to create stdout pipe: %v", err)
		return fmt.Errorf("failed to create stdout pipe: %w", err)
	}

	stderr, err := cmd.StderrPipe()
	if err != nil {
		log.Printf("❌ Failed to create stderr pipe: %v", err)
		return fmt.Errorf("failed to create stderr pipe: %w", err)
	}

	if err := cmd.Start(); err != nil {
		log.Printf("❌ Failed to start railpack: %v", err)
		return fmt.Errorf("failed to start railpack: %w", err)
	}
	log.Printf("🏗️  Railpack process started (PID: %d)", cmd.Process.Pid)

	// Monitor output
	go func() {
		scanner := bufio.NewScanner(stdout)
		for scanner.Scan() {
			line := scanner.Text()
			log.Printf("🏗️  [STDOUT] %s", line)
		}
		if err := scanner.Err(); err != nil {
			log.Printf("⚠️  [STDOUT] Scanner error: %v", err)
		}
	}()

	// Monitor errors
	errScanner := bufio.NewScanner(stderr)
	buildkitError := false
	for errScanner.Scan() {
		line := errScanner.Text()
		log.Printf("🏗️  [STDERR] %s", line)

		// Check for buildkit issues
		if strings.Contains(line, BUILDKIT_ENV_MISSING) || strings.Contains(line, BUILDKIT_UNABLE_TO_CONNECT) {
			log.Printf("❌ BUILDKIT CONNECTION ISSUE DETECTED!")
			log.Printf("❌ Line: %s", line)
			buildkitError = true
			cmd.Process.Kill()
			break
		}
	}
	if err := errScanner.Err(); err != nil {
		log.Printf("⚠️  [STDERR] Scanner error: %v", err)
	}

	if buildkitError {
		return fmt.Errorf("buildkit connection issue")
	}

	if err := cmd.Wait(); err != nil {
		log.Printf("❌ Railpack build failed: %v", err)
		return fmt.Errorf("railpack build failed: %w", err)
	}

	log.Printf("✅ Build completed successfully")
	log.Printf("✅ Image created: %s:latest", project.ID)
	return nil
}

// Helper function for min
func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func (bs *BuildService) generateDockerCompose(project *models.Project) error {
	log.Printf("📝 ============================================")
	log.Printf("📝 Generating docker-compose")
	log.Printf("📝 Project ID: %s", project.ID)
	log.Printf("📝 Domain: %s", project.Domain)
	log.Printf("📝 ============================================")

	// Create services directory
	servicePath := filepath.Join(SERVICES_BASE_PATH, project.ID)
	log.Printf("📝 Creating service directory: %s", servicePath)
	if err := os.MkdirAll(servicePath, 0755); err != nil {
		log.Printf("❌ Failed to create service directory: %v", err)
		return fmt.Errorf("failed to create service directory: %w", err)
	}
	log.Printf("✅ Service directory created")

	// Docker compose template for Traefik v3 Swarm
	tmpl := `version: '3.8'
services:
  app:
    image: {{.ProjectID}}:latest
    networks:
      - proxy
    deploy:
      restart_policy:
        condition: any
      labels:
        - "traefik.enable=true"
        - "traefik.http.routers.{{.ProjectID}}.rule=Host(` + "`{{.Domain}}`" + `)"
        - "traefik.http.routers.{{.ProjectID}}.entrypoints=web"
        - "traefik.http.services.{{.ProjectID}}.loadbalancer.server.port=3000"

networks:
  proxy:
    external: true
    name: proxy
`

	log.Printf("📝 Parsing docker-compose template...")
	t, err := template.New("docker-compose").Parse(tmpl)
	if err != nil {
		log.Printf("❌ Failed to parse template: %v", err)
		return fmt.Errorf("failed to parse template: %w", err)
	}
	log.Printf("✅ Template parsed successfully")

	// Create docker-compose file
	composePath := filepath.Join(servicePath, "docker-compose.yml")
	log.Printf("📝 Creating docker-compose file: %s", composePath)
	f, err := os.Create(composePath)
	if err != nil {
		log.Printf("❌ Failed to create docker-compose file: %v", err)
		return fmt.Errorf("failed to create docker-compose file: %w", err)
	}
	defer f.Close()

	data := map[string]interface{}{
		"ProjectID": project.ID,
		"Domain":    project.Domain,
	}

	log.Printf("📝 Writing docker-compose with data:")
	log.Printf("📝   - ProjectID: %s", project.ID)
	log.Printf("📝   - Domain: %s", project.Domain)

	if err := t.Execute(f, data); err != nil {
		log.Printf("❌ Failed to write docker-compose: %v", err)
		return fmt.Errorf("failed to write docker-compose: %w", err)
	}

	log.Printf("✅ Docker-compose generated successfully at: %s", composePath)
	return nil
}

func (bs *BuildService) deployWithSwarm(project *models.Project) error {
	log.Printf("🚀 ============================================")
	log.Printf("🚀 Deploying to Docker Swarm")
	log.Printf("🚀 Project ID: %s", project.ID)
	log.Printf("🚀 Stack name: %s", project.ID)
	log.Printf("🚀 ============================================")

	// Check and create network if it doesn't exist
	log.Printf("🚀 Checking if %s network exists...", NETWORK_NAME)
	if err := bs.ensureNetworkExists(); err != nil {
		log.Printf("❌ Failed to ensure network exists: %v", err)
		return fmt.Errorf("failed to ensure network exists: %w", err)
	}
	log.Printf("✅ Network %s is ready", NETWORK_NAME)

	composePath := filepath.Join(SERVICES_BASE_PATH, project.ID, "docker-compose.yml")
	log.Printf("🚀 Docker-compose path: %s", composePath)

	// Check if file exists
	if _, err := os.Stat(composePath); os.IsNotExist(err) {
		log.Printf("❌ Docker-compose file does not exist: %s", composePath)
		return fmt.Errorf("docker-compose file not found: %s", composePath)
	}
	log.Printf("✅ Docker-compose file exists")

	log.Printf("🚀 Executing: docker stack deploy -c %s %s", composePath, project.ID)
	cmd := exec.Command("docker", "stack", "deploy", "-c", composePath, project.ID)
	output, err := cmd.CombinedOutput()

	if err != nil {
		log.Printf("❌ Docker stack deploy failed")
		log.Printf("❌ Error: %v", err)
		log.Printf("❌ Output: %s", string(output))
		return fmt.Errorf("docker stack deploy failed: %w, output: %s", err, string(output))
	}

	log.Printf("✅ Docker stack deploy output:")
	log.Printf("%s", string(output))
	log.Printf("✅ Deployed successfully to Docker Swarm")

	return nil
}

// ensureNetworkExists checks if proxy network exists with correct scope and creates it if needed
func (bs *BuildService) ensureNetworkExists() error {
	log.Printf("🔍 Checking if network '%s' exists with correct scope...", NETWORK_NAME)

	// Check if network exists and get its details
	checkCmd := exec.Command("docker", "network", "inspect", NETWORK_NAME, "--format", "{{.Scope}}")
	output, err := checkCmd.CombinedOutput()

	if err != nil {
		// Network doesn't exist, create it
		log.Printf("⚠️  Network '%s' does not exist, creating it...", NETWORK_NAME)
		return bs.createSwarmNetwork()
	}

	scope := strings.TrimSpace(string(output))
	log.Printf("🔍 Found network '%s' with scope: %s", NETWORK_NAME, scope)

	if scope != "swarm" {
		log.Printf("⚠️  Network '%s' has wrong scope '%s' (need 'swarm')", NETWORK_NAME, scope)
		log.Printf("🗑️  Removing local network '%s'...", NETWORK_NAME)

		// Remove the local network
		removeCmd := exec.Command("docker", "network", "rm", NETWORK_NAME)
		removeOutput, removeErr := removeCmd.CombinedOutput()
		if removeErr != nil {
			log.Printf("❌ Failed to remove local network: %v", removeErr)
			log.Printf("❌ Output: %s", string(removeOutput))
			return fmt.Errorf("failed to remove local network: %w, output: %s", removeErr, string(removeOutput))
		}
		log.Printf("✅ Local network removed")

		// Create swarm network
		log.Printf("🔧 Creating swarm network '%s'...", NETWORK_NAME)
		return bs.createSwarmNetwork()
	}

	log.Printf("✅ Network '%s' exists with correct scope (swarm)", NETWORK_NAME)
	return nil
}

// createSwarmNetwork creates an overlay network for Docker Swarm
func (bs *BuildService) createSwarmNetwork() error {
	log.Printf("🔧 Creating swarm overlay network '%s'...", NETWORK_NAME)

	createCmd := exec.Command("docker", "network", "create",
		"--driver", "overlay",
		"--attachable",
		"--scope", "swarm",
		NETWORK_NAME,
	)

	output, err := createCmd.CombinedOutput()
	if err != nil {
		log.Printf("❌ Failed to create network: %v", err)
		log.Printf("❌ Output: %s", string(output))
		return fmt.Errorf("failed to create network: %w, output: %s", err, string(output))
	}

	log.Printf("✅ Swarm network '%s' created successfully", NETWORK_NAME)
	log.Printf("✅ Network ID: %s", strings.TrimSpace(string(output)))

	// Verify the network was created with correct scope
	verifyCmd := exec.Command("docker", "network", "inspect", NETWORK_NAME, "--format", "{{.Scope}}")
	verifyOutput, verifyErr := verifyCmd.CombinedOutput()
	if verifyErr == nil {
		scope := strings.TrimSpace(string(verifyOutput))
		log.Printf("✅ Verified network scope: %s", scope)
	}

	return nil
}

func (bs *BuildService) updateDeploymentStatus(projectID, status string) {
	ctx := context.Background()
	if _, err := bs.store.UpdateProjectDeploymentStatus(ctx, projectID, status); err != nil {
		log.Printf("❌ Failed to update deployment status: %v", err)
	}
}

func (bs *BuildService) broadcastStatus(projectID, status string) {
	if bs.wsHub != nil {
		bs.wsHub.BroadcastToProject(projectID, map[string]string{
			"type":   "deployment_status",
			"status": status,
		})
	}
}

func (bs *BuildService) cleanup(projectID string) {
	log.Printf("🧹 Cleaning up failed build: %s", projectID)

	// Remove repo directory
	repoPath := filepath.Join(REPO_BASE_PATH, projectID)
	if err := os.RemoveAll(repoPath); err != nil {
		log.Printf("⚠️  Failed to remove repo directory: %v", err)
	}

	// Remove service directory
	servicePath := filepath.Join(SERVICES_BASE_PATH, projectID)
	if err := os.RemoveAll(servicePath); err != nil {
		log.Printf("⚠️  Failed to remove service directory: %v", err)
	}
}

func (bs *BuildService) Shutdown() {
	log.Println("🛑 Shutting down build service...")
	bs.cancel()
	close(bs.queue)
	bs.wg.Wait()
	log.Println("✅ Build service shut down complete")
}
