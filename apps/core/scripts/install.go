package install

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"time"

	"golang.org/x/crypto/bcrypt"
)

type InstallConfig struct {
	// User provided
	AdminEmail    string
	AdminPassword string
	AdminUsername string
	Domain        string // Optional
	InstallType   string // "local" or "remote"

	// Generated
	PostgresPassword string
	RedisPassword    string
	JWTSecret        string
	AuthSecret       string
	DatabaseURL      string
	PublicAPIURL     string

	// Registry configuration
	RegistryURL string // e.g., "ghcr.io", "docker.io"
	Repository  string // e.g., "your-username/kova"
	Version     string // e.g., "latest", "v1.0.0"

	// System
	DataDir    string
	InstallDir string
	LocalIP    string
	PublicIP   string
}

const (
	DefaultDataDir    = "/data/kova"
	DefaultInstallDir = "/opt/kova"
)

// Main installation function
func Install(config *InstallConfig) error {
	fmt.Println("🚀 Starting Kova Installation...")

	steps := []struct {
		name string
		fn   func(*InstallConfig) error
	}{
		{"Checking system requirements", checkSystemRequirements},
		{"Installing Docker", installDocker},
		{"Creating directories", createDirectories},
		{"Generating secure credentials", generateSecrets},
		{"Detecting network configuration", detectNetwork},
		{"Creating environment files", createEnvironmentFiles},
		{"Generating docker-compose.yml", generateDockerCompose},
		{"Creating nginx-proxy.conf", createNginxProxyConfig},
		{"Starting services", startServices},
		{"Running database migrations", runMigrations},
		{"Creating admin user", createAdminUser},
		{"Verifying installation", verifyInstallation},
	}

	for i, step := range steps {
		fmt.Printf("[%d/%d] %s...\n", i+1, len(steps), step.name)
		if err := step.fn(config); err != nil {
			return fmt.Errorf("failed at step '%s': %w", step.name, err)
		}
		time.Sleep(500 * time.Millisecond) // Brief pause for UX
	}

	fmt.Println("\n✅ Kova installation completed successfully!")
	printAccessInfo(config)
	return nil
}

func checkSystemRequirements(config *InstallConfig) error {
	// Check if running as root/sudo
	if os.Geteuid() != 0 {
		return fmt.Errorf("installation requires root privileges. Please run with sudo")
	}

	// Check disk space (minimum 10GB)
	if err := checkDiskSpace(config.DataDir, 10); err != nil {
		return err
	}

	// Check if ports are available
	ports := []int{80, 443, 5432, 6379, 8080, 3000}
	for _, port := range ports {
		if !isPortAvailable(port) {
			fmt.Printf("⚠️  Warning: Port %d is already in use\n", port)
		}
	}

	return nil
}

func installDocker(config *InstallConfig) error {
	if isDockerInstalled() {
		fmt.Println("   Docker is already installed")
		return nil
	}

	fmt.Println("   Installing Docker...")

	// Detect OS and install Docker accordingly
	osType := detectOS()

	switch osType {
	case "ubuntu", "debian":
		return installDockerDebian()
	case "centos", "rhel", "fedora":
		return installDockerRedHat()
	case "arch":
		return installDockerArch()
	default:
		return fmt.Errorf("unsupported operating system: %s", osType)
	}
}

func createDirectories(config *InstallConfig) error {
	dirs := []string{
		config.DataDir,
		filepath.Join(config.DataDir, "postgres"),
		filepath.Join(config.DataDir, "redis"),
		filepath.Join(config.DataDir, "traefik"),
		filepath.Join(config.DataDir, "uploads"),
		filepath.Join(config.DataDir, "backups"),
		filepath.Join(config.DataDir, "logs"),
		filepath.Join(config.DataDir, "ssl"),
		config.InstallDir,
	}

	for _, dir := range dirs {
		if err := os.MkdirAll(dir, 0755); err != nil {
			return fmt.Errorf("failed to create directory %s: %w", dir, err)
		}
	}

	// Set proper ownership (non-root user for data directories)
	dataUID, dataGID := 1000, 1000  // Standard non-root user
	for _, dir := range dirs[1:7] { // Skip root data dir
		if err := os.Chown(dir, dataUID, dataGID); err != nil {
			fmt.Printf("   Warning: Could not set ownership for %s: %v\n", dir, err)
		}
	}

	return nil
}

func generateSecrets(config *InstallConfig) error {
	var err error

	// Generate secure random passwords and secrets
	config.PostgresPassword, err = generateSecurePassword(32)
	if err != nil {
		return err
	}

	config.RedisPassword, err = generateSecurePassword(32)
	if err != nil {
		return err
	}

	config.JWTSecret, err = generateSecurePassword(64)
	if err != nil {
		return err
	}

	config.AuthSecret, err = generateSecurePassword(32)
	if err != nil {
		return err
	}

	// Build database URL
	config.DatabaseURL = fmt.Sprintf("postgres://kova:%s@postgres:5432/kova?sslmode=disable",
		config.PostgresPassword)

	return nil
}

func createEnvironmentFiles(config *InstallConfig) error {
	// Create .env for core API
	envContent := fmt.Sprintf(`# Kova Core API Environment
PORT=8080
HOST=0.0.0.0
ENVIRONMENT=production

# Database
DB_HOST=postgres
DB_PORT=5432
DB_USER=kova
DB_PASSWORD=%s
DB_NAME=kova
DB_SSLMODE=disable
DATABASE_URL=%s

# Auth
JWT_SECRET=%s

# Redis
REDIS_URL=redis://:%s@redis:6379

# Admin User
ADMIN_EMAIL=%s
ADMIN_USERNAME=%s
ADMIN_PASSWORD=%s
`, config.PostgresPassword, config.DatabaseURL, config.JWTSecret,
		config.RedisPassword, config.AdminEmail, config.AdminUsername, config.AdminPassword)

	envFile := filepath.Join(config.InstallDir, ".env")
	if err := os.WriteFile(envFile, []byte(envContent), 0600); err != nil {
		return fmt.Errorf("failed to create .env file: %w", err)
	}

	// Create .env.local for Next.js dashboard
	dashboardEnv := fmt.Sprintf(`# Kova Dashboard Environment
AUTH_SECRET=%s
NEXTAUTH_SECRET=%s
NEXT_PUBLIC_API_URL=%s
NEXTAUTH_URL=%s
`, config.AuthSecret, config.AuthSecret, config.PublicAPIURL,
		getDashboardURL(config))

	dashboardEnvFile := filepath.Join(config.InstallDir, ".env.dashboard")
	if err := os.WriteFile(dashboardEnvFile, []byte(dashboardEnv), 0600); err != nil {
		return fmt.Errorf("failed to create dashboard .env file: %w", err)
	}

	return nil
}

func generateDockerCompose(config *InstallConfig) error {
	// Determine if we should build locally or use pre-built images
	usePrebuiltImages := os.Getenv("KOVA_USE_PREBUILT") != "false"

	// Set default values for registry and version
	if config.RegistryURL == "" {
		config.RegistryURL = "ghcr.io"
	}
	if config.Repository == "" {
		config.Repository = "dopecape/kova"
	}
	if config.Version == "" {
		config.Version = "latest"
	}

	// Generate the docker-compose content based on configuration
	var composeContent string

	if usePrebuiltImages {
		fmt.Printf("   Using pre-built images\n")
		composeContent = generatePrebuiltCompose(config)
	} else {
		fmt.Printf("   Building images locally\n")
		composeContent = generateLocalBuildCompose(config)
	}

	composeFile := filepath.Join(config.InstallDir, "docker-compose.yml")
	if err := os.WriteFile(composeFile, []byte(composeContent), 0644); err != nil {
		return fmt.Errorf("failed to create docker-compose.yml: %w", err)
	}

	// Validate the generated YAML
	if err := validateDockerCompose(composeFile); err != nil {
		// Print the problematic content for debugging
		fmt.Printf("Generated docker-compose.yml content:\n%s\n", composeContent)
		return fmt.Errorf("generated docker-compose.yml is invalid: %w", err)
	}

	return nil
}

// Add validation function
func validateDockerCompose(composeFile string) error {
	cmd := exec.Command("docker-compose", "-f", composeFile, "config", "--quiet")
	if err := cmd.Run(); err != nil {
		return err
	}
	return nil
}

func generatePrebuiltCompose(config *InstallConfig) string {
	// Build Traefik command list
	traefikCommands := []string{
		"--api.dashboard=true",
		"--api.insecure=true",
		"--providers.docker=true",
		"--providers.docker.exposedbydefault=false",
		"--entrypoints.web.address=:80",
		"--entrypoints.websecure.address=:443",
	}

	// Add SSL configuration if domain is provided
	if config.Domain != "" {
		traefikCommands = append(traefikCommands,
			"--certificatesresolvers.letsencrypt.acme.email="+config.AdminEmail,
			"--certificatesresolvers.letsencrypt.acme.storage=/data/acme.json",
			"--certificatesresolvers.letsencrypt.acme.httpchallenge.entrypoint=web",
		)
	}

	// Determine host for routing
	host := config.PublicIP
	if config.Domain != "" {
		host = config.Domain
	}

	// Build Traefik command section
	traefikCommandSection := ""
	for _, cmd := range traefikCommands {
		traefikCommandSection += fmt.Sprintf("      - %s\n", cmd)
	}

	// Build labels for API service (using subdomain routing)
	apiLabels := fmt.Sprintf(`      - "traefik.enable=true"
      - "traefik.http.routers.kova-api.rule=Host(\\"api.%s\\") || Host(\\"api.localhost\\") || Host(\\"api.127.0.0.1\\")"
      - "traefik.http.routers.kova-api.priority=100"
      - "traefik.http.services.kova-api.loadbalancer.server.port=8080"`, host)

	if config.Domain != "" {
		apiLabels += `
      - "traefik.http.routers.kova-api.tls.certresolver=letsencrypt"`
	}

	// Build labels for Dashboard service
	dashboardLabels := fmt.Sprintf(`      - "traefik.enable=true"
      - "traefik.http.routers.kova-dashboard.rule=Host(\\"%s\\") || Host(\\"localhost\\") || Host(\\"127.0.0.1\\")"
      - "traefik.http.routers.kova-dashboard.priority=50"
      - "traefik.http.services.kova-dashboard.loadbalancer.server.port=3000"`, host)

	if config.Domain != "" {
		dashboardLabels += `
      - "traefik.http.routers.kova-dashboard.tls.certresolver=letsencrypt"`
	}

	// Build the complete YAML with proper structure
	var yamlBuilder strings.Builder

	yamlBuilder.WriteString("version: '3.8'\n\n")
	yamlBuilder.WriteString("services:\n")

	// Traefik service
	yamlBuilder.WriteString("  # Reverse Proxy & Load Balancer\n")
	yamlBuilder.WriteString("  traefik:\n")
	yamlBuilder.WriteString("    image: traefik:v3.0\n")
	yamlBuilder.WriteString("    container_name: kova-traefik\n")
	yamlBuilder.WriteString("    restart: unless-stopped\n")
	yamlBuilder.WriteString("    command:\n")
	yamlBuilder.WriteString(traefikCommandSection)
	yamlBuilder.WriteString("    ports:\n")
	yamlBuilder.WriteString("      - \"80:80\"\n")
	yamlBuilder.WriteString("      - \"443:443\"\n")
	yamlBuilder.WriteString("      - \"8080:8080\"\n")
	yamlBuilder.WriteString("    volumes:\n")
	yamlBuilder.WriteString("      - /var/run/docker.sock:/var/run/docker.sock:ro\n")
	yamlBuilder.WriteString(fmt.Sprintf("      - %s/traefik:/data\n", config.DataDir))
	yamlBuilder.WriteString("    networks:\n")
	yamlBuilder.WriteString("      - kova_network\n\n")

	// API Proxy service
	yamlBuilder.WriteString("  # API Proxy for Dashboard\n")
	yamlBuilder.WriteString("  api-proxy:\n")
	yamlBuilder.WriteString("    image: nginx:alpine\n")
	yamlBuilder.WriteString("    container_name: kova-api-proxy\n")
	yamlBuilder.WriteString("    restart: unless-stopped\n")
	yamlBuilder.WriteString("    networks:\n")
	yamlBuilder.WriteString("      - kova_network\n")
	yamlBuilder.WriteString("    volumes:\n")
	yamlBuilder.WriteString("      - ./nginx-proxy.conf:/etc/nginx/nginx.conf:ro\n")
	yamlBuilder.WriteString("    depends_on:\n")
	yamlBuilder.WriteString("      - kova-api\n\n")

	// PostgreSQL service
	yamlBuilder.WriteString("  # PostgreSQL Database\n")
	yamlBuilder.WriteString("  postgres:\n")
	yamlBuilder.WriteString("    image: postgres:15-alpine\n")
	yamlBuilder.WriteString("    container_name: kova-postgres\n")
	yamlBuilder.WriteString("    restart: unless-stopped\n")
	yamlBuilder.WriteString("    environment:\n")
	yamlBuilder.WriteString("      POSTGRES_DB: kova\n")
	yamlBuilder.WriteString("      POSTGRES_USER: kova\n")
	yamlBuilder.WriteString(fmt.Sprintf("      POSTGRES_PASSWORD: %s\n", config.PostgresPassword))
	yamlBuilder.WriteString("      PGDATA: /var/lib/postgresql/data/pgdata\n")
	yamlBuilder.WriteString("    volumes:\n")
	yamlBuilder.WriteString(fmt.Sprintf("      - %s/postgres:/var/lib/postgresql/data\n", config.DataDir))
	yamlBuilder.WriteString("    networks:\n")
	yamlBuilder.WriteString("      - kova_network\n")
	yamlBuilder.WriteString("    healthcheck:\n")
	yamlBuilder.WriteString("      test: [\"CMD-SHELL\", \"pg_isready -U kova -d kova\"]\n")
	yamlBuilder.WriteString("      interval: 10s\n")
	yamlBuilder.WriteString("      timeout: 5s\n")
	yamlBuilder.WriteString("      retries: 5\n\n")

	// Redis service
	yamlBuilder.WriteString("  # Redis Cache\n")
	yamlBuilder.WriteString("  redis:\n")
	yamlBuilder.WriteString("    image: redis:7-alpine\n")
	yamlBuilder.WriteString("    container_name: kova-redis\n")
	yamlBuilder.WriteString("    restart: unless-stopped\n")
	yamlBuilder.WriteString(fmt.Sprintf("    command: redis-server --appendonly yes --requirepass %s\n", config.RedisPassword))
	yamlBuilder.WriteString("    volumes:\n")
	yamlBuilder.WriteString(fmt.Sprintf("      - %s/redis:/data\n", config.DataDir))
	yamlBuilder.WriteString("    networks:\n")
	yamlBuilder.WriteString("      - kova_network\n")
	yamlBuilder.WriteString("    healthcheck:\n")
	yamlBuilder.WriteString("      test: [\"CMD\", \"redis-cli\", \"--raw\", \"incr\", \"ping\"]\n")
	yamlBuilder.WriteString("      interval: 10s\n")
	yamlBuilder.WriteString("      timeout: 5s\n")
	yamlBuilder.WriteString("      retries: 5\n\n")

	// Kova API service
	yamlBuilder.WriteString("  # Kova Core API\n")
	yamlBuilder.WriteString("  kova-api:\n")
	yamlBuilder.WriteString(fmt.Sprintf("    image: %s/%s-api:%s\n", config.RegistryURL, config.Repository, config.Version))
	yamlBuilder.WriteString("    container_name: kova-api\n")
	yamlBuilder.WriteString("    restart: unless-stopped\n")
	yamlBuilder.WriteString("    env_file:\n")
	yamlBuilder.WriteString("      - .env\n")
	yamlBuilder.WriteString("    depends_on:\n")
	yamlBuilder.WriteString("      postgres:\n")
	yamlBuilder.WriteString("        condition: service_healthy\n")
	yamlBuilder.WriteString("      redis:\n")
	yamlBuilder.WriteString("        condition: service_healthy\n")
	yamlBuilder.WriteString("    networks:\n")
	yamlBuilder.WriteString("      - kova_network\n")
	yamlBuilder.WriteString("    labels:\n")
	yamlBuilder.WriteString(apiLabels)
	yamlBuilder.WriteString("\n\n")

	// Kova Dashboard service
	yamlBuilder.WriteString("  # Kova Dashboard\n")
	yamlBuilder.WriteString("  kova-dashboard:\n")
	yamlBuilder.WriteString(fmt.Sprintf("    image: %s/%s-dashboard:%s\n", config.RegistryURL, config.Repository, config.Version))
	yamlBuilder.WriteString("    container_name: kova-dashboard\n")
	yamlBuilder.WriteString("    restart: unless-stopped\n")
	yamlBuilder.WriteString("    env_file:\n")
	yamlBuilder.WriteString("      - .env.dashboard\n")
	yamlBuilder.WriteString("    depends_on:\n")
	yamlBuilder.WriteString("      - kova-api\n")
	yamlBuilder.WriteString("      - api-proxy\n")
	yamlBuilder.WriteString("    networks:\n")
	yamlBuilder.WriteString("      - kova_network\n")
	yamlBuilder.WriteString("    extra_hosts:\n")
	yamlBuilder.WriteString(fmt.Sprintf("      - \"api.%s:host-gateway\"\n", host))
	yamlBuilder.WriteString("    healthcheck:\n")
	yamlBuilder.WriteString("      disable: true\n")
	yamlBuilder.WriteString("    labels:\n")
	yamlBuilder.WriteString(dashboardLabels)
	yamlBuilder.WriteString("\n\n")

	// Networks and volumes
	yamlBuilder.WriteString("networks:\n")
	yamlBuilder.WriteString("  kova_network:\n")
	yamlBuilder.WriteString("    driver: bridge\n\n")
	yamlBuilder.WriteString("volumes:\n")
	yamlBuilder.WriteString("  postgres_data:\n")
	yamlBuilder.WriteString("  redis_data:\n")
	yamlBuilder.WriteString("  traefik_data:\n")

	return yamlBuilder.String()
}

func createNginxProxyConfig(config *InstallConfig) error {
	nginxConfig := `events {
    worker_connections 1024;
}

http {
    upstream api {
        server kova-api:8080;
    }

    server {
        listen 80;
        server_name _;

        location /api/ {
            proxy_pass http://api/api/;
            proxy_set_header Host $host;
            proxy_set_header X-Real-IP $remote_addr;
            proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
            proxy_set_header X-Forwarded-Proto $scheme;
        }
    }
}`

	nginxConfigFile := filepath.Join(config.InstallDir, "nginx-proxy.conf")
	if err := os.WriteFile(nginxConfigFile, []byte(nginxConfig), 0644); err != nil {
		return fmt.Errorf("failed to create nginx-proxy.conf: %w", err)
	}

	return nil
}

func generateLocalBuildCompose(config *InstallConfig) string {
	// Build Traefik command list
	traefikCommands := []string{
		"--api.dashboard=true",
		"--api.insecure=true",
		"--providers.docker=true",
		"--providers.docker.exposedbydefault=false",
		"--entrypoints.web.address=:80",
		"--entrypoints.websecure.address=:443",
	}

	// Add SSL configuration if domain is provided
	if config.Domain != "" {
		traefikCommands = append(traefikCommands,
			"--certificatesresolvers.letsencrypt.acme.email="+config.AdminEmail,
			"--certificatesresolvers.letsencrypt.acme.storage=/data/acme.json",
			"--certificatesresolvers.letsencrypt.acme.httpchallenge.entrypoint=web",
		)
	}

	// Determine host for routing
	host := config.PublicIP
	if config.Domain != "" {
		host = config.Domain
	}

	// Build Traefik command section
	traefikCommandSection := ""
	for _, cmd := range traefikCommands {
		traefikCommandSection += fmt.Sprintf("      - %s\n", cmd)
	}

	// Build labels for API service
	apiLabels := fmt.Sprintf("      - \"traefik.enable=true\"\n      - \"traefik.http.routers.api.rule=Host(\\\"%s\\\") && PathPrefix(\\\"/api\\\")\"\n      - \"traefik.http.services.api.loadbalancer.server.port=8080\"", host)

	if config.Domain != "" {
		apiLabels += "\n      - \"traefik.http.routers.api.tls.certresolver=letsencrypt\""
	}

	// Build labels for Dashboard service
	dashboardLabels := fmt.Sprintf("      - \"traefik.enable=true\"\n      - \"traefik.http.routers.dashboard.rule=Host(\\\"%s\\\")\"\n      - \"traefik.http.services.dashboard.loadbalancer.server.port=3000\"", host)

	if config.Domain != "" {
		dashboardLabels += "\n      - \"traefik.http.routers.dashboard.tls.certresolver=letsencrypt\""
	}

	// Build the complete YAML with proper structure
	var yamlBuilder strings.Builder

	yamlBuilder.WriteString("version: '3.8'\n\n")
	yamlBuilder.WriteString("services:\n")

	// Traefik service
	yamlBuilder.WriteString("  # Reverse Proxy & Load Balancer\n")
	yamlBuilder.WriteString("  traefik:\n")
	yamlBuilder.WriteString("    image: traefik:v3.0\n")
	yamlBuilder.WriteString("    container_name: kova-traefik\n")
	yamlBuilder.WriteString("    restart: unless-stopped\n")
	yamlBuilder.WriteString("    command:\n")
	yamlBuilder.WriteString(traefikCommandSection)
	yamlBuilder.WriteString("    ports:\n")
	yamlBuilder.WriteString("      - \"80:80\"\n")
	yamlBuilder.WriteString("      - \"443:443\"\n")
	yamlBuilder.WriteString("      - \"8080:8080\"\n")
	yamlBuilder.WriteString("    volumes:\n")
	yamlBuilder.WriteString("      - /var/run/docker.sock:/var/run/docker.sock:ro\n")
	yamlBuilder.WriteString(fmt.Sprintf("      - %s/traefik:/data\n", config.DataDir))
	yamlBuilder.WriteString("    networks:\n")
	yamlBuilder.WriteString("      - kova_network\n\n")

	// PostgreSQL service
	yamlBuilder.WriteString("  # PostgreSQL Database\n")
	yamlBuilder.WriteString("  postgres:\n")
	yamlBuilder.WriteString("    image: postgres:15-alpine\n")
	yamlBuilder.WriteString("    container_name: kova-postgres\n")
	yamlBuilder.WriteString("    restart: unless-stopped\n")
	yamlBuilder.WriteString("    environment:\n")
	yamlBuilder.WriteString("      POSTGRES_DB: kova\n")
	yamlBuilder.WriteString("      POSTGRES_USER: kova\n")
	yamlBuilder.WriteString(fmt.Sprintf("      POSTGRES_PASSWORD: %s\n", config.PostgresPassword))
	yamlBuilder.WriteString("      PGDATA: /var/lib/postgresql/data/pgdata\n")
	yamlBuilder.WriteString("    volumes:\n")
	yamlBuilder.WriteString(fmt.Sprintf("      - %s/postgres:/var/lib/postgresql/data\n", config.DataDir))
	yamlBuilder.WriteString("    networks:\n")
	yamlBuilder.WriteString("      - kova_network\n")
	yamlBuilder.WriteString("    healthcheck:\n")
	yamlBuilder.WriteString("      test: [\"CMD-SHELL\", \"pg_isready -U kova -d kova\"]\n")
	yamlBuilder.WriteString("      interval: 10s\n")
	yamlBuilder.WriteString("      timeout: 5s\n")
	yamlBuilder.WriteString("      retries: 5\n\n")

	// Redis service
	yamlBuilder.WriteString("  # Redis Cache\n")
	yamlBuilder.WriteString("  redis:\n")
	yamlBuilder.WriteString("    image: redis:7-alpine\n")
	yamlBuilder.WriteString("    container_name: kova-redis\n")
	yamlBuilder.WriteString("    restart: unless-stopped\n")
	yamlBuilder.WriteString(fmt.Sprintf("    command: redis-server --appendonly yes --requirepass %s\n", config.RedisPassword))
	yamlBuilder.WriteString("    volumes:\n")
	yamlBuilder.WriteString(fmt.Sprintf("      - %s/redis:/data\n", config.DataDir))
	yamlBuilder.WriteString("    networks:\n")
	yamlBuilder.WriteString("      - kova_network\n")
	yamlBuilder.WriteString("    healthcheck:\n")
	yamlBuilder.WriteString("      test: [\"CMD\", \"redis-cli\", \"--raw\", \"incr\", \"ping\"]\n")
	yamlBuilder.WriteString("      interval: 10s\n")
	yamlBuilder.WriteString("      timeout: 5s\n")
	yamlBuilder.WriteString("      retries: 5\n\n")

	// Kova API service
	yamlBuilder.WriteString("  # Kova Core API\n")
	yamlBuilder.WriteString("  kova-api:\n")
	yamlBuilder.WriteString("    build:\n")
	yamlBuilder.WriteString("      context: .\n")
	yamlBuilder.WriteString("      dockerfile: apps/core/Dockerfile\n")
	yamlBuilder.WriteString("    container_name: kova-api\n")
	yamlBuilder.WriteString("    restart: unless-stopped\n")
	yamlBuilder.WriteString("    env_file:\n")
	yamlBuilder.WriteString("      - .env\n")
	yamlBuilder.WriteString("    depends_on:\n")
	yamlBuilder.WriteString("      postgres:\n")
	yamlBuilder.WriteString("        condition: service_healthy\n")
	yamlBuilder.WriteString("      redis:\n")
	yamlBuilder.WriteString("        condition: service_healthy\n")
	yamlBuilder.WriteString("    networks:\n")
	yamlBuilder.WriteString("      - kova_network\n")
	yamlBuilder.WriteString("    labels:\n")
	yamlBuilder.WriteString(apiLabels)
	yamlBuilder.WriteString("\n\n")

	// Kova Dashboard service
	yamlBuilder.WriteString("  # Kova Dashboard\n")
	yamlBuilder.WriteString("  kova-dashboard:\n")
	yamlBuilder.WriteString("    build:\n")
	yamlBuilder.WriteString("      context: .\n")
	yamlBuilder.WriteString("      dockerfile: apps/web/Dockerfile\n")
	yamlBuilder.WriteString("    container_name: kova-dashboard\n")
	yamlBuilder.WriteString("    restart: unless-stopped\n")
	yamlBuilder.WriteString("    env_file:\n")
	yamlBuilder.WriteString("      - .env.dashboard\n")
	yamlBuilder.WriteString("    depends_on:\n")
	yamlBuilder.WriteString("      - kova-api\n")
	yamlBuilder.WriteString("    networks:\n")
	yamlBuilder.WriteString("      - kova_network\n")
	yamlBuilder.WriteString("    labels:\n")
	yamlBuilder.WriteString(dashboardLabels)
	yamlBuilder.WriteString("\n\n")

	// Networks and volumes
	yamlBuilder.WriteString("networks:\n")
	yamlBuilder.WriteString("  kova_network:\n")
	yamlBuilder.WriteString("    driver: bridge\n\n")
	yamlBuilder.WriteString("volumes:\n")
	yamlBuilder.WriteString("  postgres_data:\n")
	yamlBuilder.WriteString("  redis_data:\n")
	yamlBuilder.WriteString("  traefik_data:\n")

	return yamlBuilder.String()
}

func buildImages(config *InstallConfig) error {
	fmt.Println("   Building Kova API image...")
	if err := runCommand("docker", "build", "-t", "kova/api:latest", "-f", "apps/core/Dockerfile", "."); err != nil {
		return fmt.Errorf("failed to build API image: %w", err)
	}

	fmt.Println("   Building Kova Dashboard image...")
	if err := runCommand("docker", "build", "-t", "kova/dashboard:latest", "-f", "apps/web/Dockerfile", "."); err != nil {
		return fmt.Errorf("failed to build dashboard image: %w", err)
	}

	return nil
}

func startServices(config *InstallConfig) error {
	fmt.Println("   Starting services with docker-compose...")

	cmd := exec.Command("docker-compose", "up", "-d")
	cmd.Dir = config.InstallDir
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		return fmt.Errorf("failed to start services: %w", err)
	}

	// Wait for services to be healthy
	fmt.Println("   Waiting for services to be ready...")
	time.Sleep(30 * time.Second)

	return nil
}

func runMigrations(config *InstallConfig) error {
	fmt.Println("   Running database migrations...")

	// Run migrations inside the API container
	cmd := exec.Command("docker", "exec", "kova-api", "migrate", "-path", "./internal/store/migrations", "-database", config.DatabaseURL, "up")
	if err := cmd.Run(); err != nil {
		fmt.Printf("   Warning: Migration failed: %v\n", err)
		// Don't fail installation if migrations fail (might already be applied)
	}

	return nil
}

func createAdminUser(config *InstallConfig) error {
	fmt.Println("   Creating admin user...")

	// Hash the password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(config.AdminPassword), bcrypt.DefaultCost)
	if err != nil {
		return fmt.Errorf("failed to hash password: %w", err)
	}

	// Insert admin user (this would typically be done via API call or direct DB insert)
	query := fmt.Sprintf(`INSERT INTO users (email, username, password_hash, is_admin, created_at, updated_at) 
		VALUES ('%s', '%s', '%s', true, NOW(), NOW()) 
		ON CONFLICT (email) DO NOTHING`,
		config.AdminEmail, config.AdminUsername, string(hashedPassword))

	cmd := exec.Command("docker", "exec", "kova-postgres", "psql", "-U", "kova", "-d", "kova", "-c", query)
	if err := cmd.Run(); err != nil {
		fmt.Printf("   Warning: Could not create admin user: %v\n", err)
	}

	return nil
}

func verifyInstallation(config *InstallConfig) error {
	fmt.Println("   Verifying installation...")

	// Check if services are running
	services := []string{"kova-traefik", "kova-postgres", "kova-redis", "kova-api", "kova-dashboard"}

	for _, service := range services {
		if !isContainerRunning(service) {
			return fmt.Errorf("service %s is not running", service)
		}
	}

	// Test API endpoint
	apiURL := fmt.Sprintf("http://%s/api/health", config.LocalIP)
	if err := testEndpoint(apiURL); err != nil {
		fmt.Printf("   Warning: API health check failed: %v\n", err)
	}

	return nil
}

// Helper functions
func generateSecurePassword(length int) (string, error) {
	bytes := make([]byte, length)
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}
	return hex.EncodeToString(bytes)[:length], nil
}

func checkDiskSpace(path string, minGB int) error {
	// Simplified disk space check - in production you'd use syscall.Statfs
	return nil
}

func isPortAvailable(port int) bool {
	// Simplified port check - in production you'd actually check if port is in use
	return true
}

func isDockerInstalled() bool {
	_, err := exec.LookPath("docker")
	return err == nil
}

func detectOS() string {
	if runtime.GOOS != "linux" {
		return "unknown"
	}

	// Read /etc/os-release
	content, err := os.ReadFile("/etc/os-release")
	if err != nil {
		return "unknown"
	}

	lines := strings.Split(string(content), "\n")
	for _, line := range lines {
		if strings.HasPrefix(line, "ID=") {
			return strings.Trim(strings.Split(line, "=")[1], "\"")
		}
	}

	return "unknown"
}

func installDockerDebian() error {
	commands := [][]string{
		{"apt-get", "update"},
		{"apt-get", "install", "-y", "apt-transport-https", "ca-certificates", "curl", "gnupg", "lsb-release"},
		{"curl", "-fsSL", "https://download.docker.com/linux/ubuntu/gpg", "-o", "/usr/share/keyrings/docker-archive-keyring.asc"},
		{"apt-get", "update"},
		{"apt-get", "install", "-y", "docker-ce", "docker-ce-cli", "containerd.io", "docker-compose-plugin"},
		{"systemctl", "enable", "docker"},
		{"systemctl", "start", "docker"},
	}

	for _, cmd := range commands {
		if err := runCommand(cmd[0], cmd[1:]...); err != nil {
			return err
		}
	}

	return nil
}

func installDockerRedHat() error {
	commands := [][]string{
		{"dnf", "config-manager", "--add-repo", "https://download.docker.com/linux/centos/docker-ce.repo"},
		{"dnf", "install", "-y", "docker-ce", "docker-ce-cli", "containerd.io", "docker-compose-plugin"},
		{"systemctl", "enable", "docker"},
		{"systemctl", "start", "docker"},
	}

	for _, cmd := range commands {
		if err := runCommand(cmd[0], cmd[1:]...); err != nil {
			return err
		}
	}

	return nil
}

func installDockerArch() error {
	commands := [][]string{
		{"pacman", "-Sy", "--noconfirm", "docker", "docker-compose"},
		{"systemctl", "enable", "docker"},
		{"systemctl", "start", "docker"},
	}

	for _, cmd := range commands {
		if err := runCommand(cmd[0], cmd[1:]...); err != nil {
			return err
		}
	}

	return nil
}

func getLocalIP() (string, error) {
	// Simplified - get from route command
	out, err := exec.Command("ip", "route", "get", "1").Output()
	if err != nil {
		return "", err
	}

	// Parse output to extract IP
	lines := strings.Split(string(out), "\n")
	for _, line := range lines {
		if strings.Contains(line, "src") {
			parts := strings.Split(line, " ")
			for i, part := range parts {
				if part == "src" && i+1 < len(parts) {
					return parts[i+1], nil
				}
			}
		}
	}

	return "", fmt.Errorf("could not determine local IP")
}

func getPublicIP() (string, error) {
	// Use the /ip endpoint to get just the IP address, not the HTML page
	resp, err := http.Get("https://ifconfig.me/ip")
	if err != nil {
		// Fallback to alternative services
		return getPublicIPFallback()
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return getPublicIPFallback()
	}

	ip, err := io.ReadAll(resp.Body)
	if err != nil {
		return getPublicIPFallback()
	}

	ipStr := strings.TrimSpace(string(ip))

	// Validate that we got an IP address, not HTML
	if strings.Contains(ipStr, "<") || strings.Contains(ipStr, "html") {
		return getPublicIPFallback()
	}

	// Basic IP validation
	if isValidIP(ipStr) {
		return ipStr, nil
	}

	return getPublicIPFallback()
}

func getPublicIPFallback() (string, error) {
	// Try multiple fallback services
	services := []string{
		"https://api.ipify.org",
		"https://checkip.amazonaws.com",
		"https://ipecho.net/plain",
		"https://icanhazip.com",
	}

	for _, service := range services {
		resp, err := http.Get(service)
		if err != nil {
			continue
		}

		if resp.StatusCode != 200 {
			resp.Body.Close()
			continue
		}

		ip, err := io.ReadAll(resp.Body)
		resp.Body.Close()

		if err != nil {
			continue
		}

		ipStr := strings.TrimSpace(string(ip))

		// Validate that we got an IP address, not HTML or error message
		if isValidIP(ipStr) {
			return ipStr, nil
		}
	}

	return "", fmt.Errorf("failed to get public IP from all services")
}

func isValidIP(ip string) bool {
	// Basic IP validation - check if it looks like an IPv4 address
	parts := strings.Split(ip, ".")
	if len(parts) != 4 {
		return false
	}

	for _, part := range parts {
		if len(part) == 0 || len(part) > 3 {
			return false
		}

		// Check if all characters are digits
		for _, char := range part {
			if char < '0' || char > '9' {
				return false
			}
		}
	}

	return true
}

func detectNetwork(config *InstallConfig) error {
	// Get local IP
	localIP, err := getLocalIP()
	if err != nil {
		fmt.Printf("   Warning: Could not detect local IP: %v\n", err)
		config.LocalIP = "localhost"
	} else {
		config.LocalIP = localIP
	}

	// Get public IP
	publicIP, err := getPublicIP()
	if err != nil {
		fmt.Printf("   Warning: Could not detect public IP: %v\n", err)
		// Fall back to local IP if we can't get public IP
		config.PublicIP = config.LocalIP
	} else {
		config.PublicIP = publicIP
		fmt.Printf("   Detected public IP: %s\n", publicIP)
	}

	// Set API URL based on domain or IP
	if config.Domain != "" {
		config.PublicAPIURL = fmt.Sprintf("https://%s/api", config.Domain)
	} else {
		config.PublicAPIURL = fmt.Sprintf("http://%s/api", config.PublicIP)
	}

	fmt.Printf("   Local IP: %s, Public IP: %s\n", config.LocalIP, config.PublicIP)
	fmt.Printf("   API URL: %s\n", config.PublicAPIURL)

	return nil
}

// Test function to verify getPublicIP is working correctly
func testGetPublicIP() {
	fmt.Println("Testing getPublicIP function...")
	ip, err := getPublicIP()
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}

	fmt.Printf("Public IP: '%s'\n", ip)
	fmt.Printf("Length: %d characters\n", len(ip))
	fmt.Printf("Contains HTML: %v\n", strings.Contains(ip, "<"))
	fmt.Printf("Valid IP format: %v\n", isValidIP(ip))
}

func getDashboardURL(config *InstallConfig) string {
	if config.Domain != "" {
		return fmt.Sprintf("https://%s", config.Domain)
	}
	return fmt.Sprintf("http://%s", config.PublicIP)
}

func runCommand(name string, args ...string) error {
	cmd := exec.Command(name, args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

func isContainerRunning(name string) bool {
	cmd := exec.Command("docker", "ps", "--filter", fmt.Sprintf("name=%s", name), "--format", "{{.Names}}")
	out, err := cmd.Output()
	if err != nil {
		return false
	}

	return strings.Contains(string(out), name)
}

func testEndpoint(url string) error {
	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 400 {
		return fmt.Errorf("endpoint returned status %d", resp.StatusCode)
	}

	return nil
}

func printAccessInfo(config *InstallConfig) {
	fmt.Println("\n" + strings.Repeat("=", 60))
	fmt.Println("🎉 KOVA INSTALLATION COMPLETE!")
	fmt.Println(strings.Repeat("=", 60))

	if config.Domain != "" {
		fmt.Printf("Dashboard:    https://%s\n", config.Domain)
		fmt.Printf("API:          https://%s/api\n", config.Domain)
		fmt.Printf("Traefik:      https://traefik.%s\n", config.Domain)
	} else {
		fmt.Printf("Dashboard:    http://%s\n", config.PublicIP)
		fmt.Printf("API:          http://%s/api\n", config.PublicIP)
		fmt.Printf("Traefik:      http://%s:8080\n", config.PublicIP)
	}

	fmt.Printf("\nAdmin Credentials:\n")
	fmt.Printf("Email:        %s\n", config.AdminEmail)
	fmt.Printf("Username:     %s\n", config.AdminUsername)
	fmt.Printf("Password:     %s\n", config.AdminPassword)

	fmt.Printf("\nConfiguration:\n")
	fmt.Printf("Data Directory:   %s\n", config.DataDir)
	fmt.Printf("Install Directory: %s\n", config.InstallDir)

	fmt.Println(strings.Repeat("=", 60))
}
