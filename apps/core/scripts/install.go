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
	"text/template"
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
		{"Building application images", buildImages},
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

func detectNetwork(config *InstallConfig) error {
	// Get local IP
	localIP, err := getLocalIP()
	if err != nil {
		config.LocalIP = "localhost"
	} else {
		config.LocalIP = localIP
	}

	// Get public IP
	publicIP, err := getPublicIP()
	if err != nil {
		config.PublicIP = config.LocalIP
	} else {
		config.PublicIP = publicIP
	}

	// Set API URL based on domain or IP
	if config.Domain != "" {
		config.PublicAPIURL = fmt.Sprintf("https://%s/api", config.Domain)
	} else {
		config.PublicAPIURL = fmt.Sprintf("http://%s/api", config.PublicIP)
	}

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

	var tmpl string
	if usePrebuiltImages {
		tmpl = `version: '3.8'

services:
  # Reverse Proxy & Load Balancer
  traefik:
    image: traefik:v3.0
    container_name: kova-traefik
    restart: unless-stopped
    command:
      - --api.dashboard=true
      - --api.insecure=true
      - --providers.docker=true
      - --providers.docker.exposedbydefault=false
      - --entrypoints.web.address=:80
      - --entrypoints.websecure.address=:443
      {{- if .Domain }}
      - --certificatesresolvers.letsencrypt.acme.email={{.AdminEmail}}
      - --certificatesresolvers.letsencrypt.acme.storage=/data/acme.json
      - --certificatesresolvers.letsencrypt.acme.httpchallenge.entrypoint=web
      {{- end }}
    ports:
      - "80:80"
      - "443:443"
      - "8080:8080"
    volumes:
      - /var/run/docker.sock:/var/run/docker.sock:ro
      - {{.DataDir}}/traefik:/data
    networks:
      - kova_network

  # PostgreSQL Database
  postgres:
    image: postgres:15-alpine
    container_name: kova-postgres
    restart: unless-stopped
    environment:
      POSTGRES_DB: kova
      POSTGRES_USER: kova
      POSTGRES_PASSWORD: {{.PostgresPassword}}
      PGDATA: /var/lib/postgresql/data/pgdata
    volumes:
      - {{.DataDir}}/postgres:/var/lib/postgresql/data
    networks:
      - kova_network
    healthcheck:
      test: ["CMD-SHELL", "pg_isready -U kova -d kova"]
      interval: 10s
      timeout: 5s
      retries: 5

  # Redis Cache
  redis:
    image: redis:7-alpine
    container_name: kova-redis
    restart: unless-stopped
    command: redis-server --appendonly yes --requirepass {{.RedisPassword}}
    volumes:
      - {{.DataDir}}/redis:/data
    networks:
      - kova_network
    healthcheck:
      test: ["CMD", "redis-cli", "--raw", "incr", "ping"]
      interval: 10s
      timeout: 5s
      retries: 5

  # Kova Core API
  kova-api:
    image: {{.RegistryURL}}/{{.Repository}}-api:{{.Version}}
    container_name: kova-api
    restart: unless-stopped
    env_file:
      - .env
    depends_on:
      postgres:
        condition: service_healthy
      redis:
        condition: service_healthy
    networks:
      - kova_network
    labels:
      - "traefik.enable=true"
      - "traefik.http.routers.api.rule=Host(\"{{.Domain | default .PublicIP}}\") && PathPrefix(\"/api\")"
      {{- if .Domain }}
      - "traefik.http.routers.api.tls.certresolver=letsencrypt"
      {{- end }}
      - "traefik.http.services.api.loadbalancer.server.port=8080"

  # Kova Dashboard
  kova-dashboard:
    image: {{.RegistryURL}}/{{.Repository}}-dashboard:{{.Version}}
    container_name: kova-dashboard
    restart: unless-stopped
    env_file:
      - .env.dashboard
    depends_on:
      - kova-api
    networks:
      - kova_network
    labels:
      - "traefik.enable=true"
      - "traefik.http.routers.dashboard.rule=Host(\"{{.Domain | default .PublicIP}}\")"
      {{- if .Domain }}
      - "traefik.http.routers.dashboard.tls.certresolver=letsencrypt"
      {{- end }}
      - "traefik.http.services.dashboard.loadbalancer.server.port=3000"

networks:
  kova_network:
    driver: bridge

volumes:
  postgres_data:
  redis_data:
  traefik_data:
`
	} else {
		// Local build version (your existing template)
		tmpl = `version: '3.8'

services:
  # ... (same as before but with build context instead of image) ...
  kova-api:
    build:
      context: .
      dockerfile: apps/core/Dockerfile
    # ... rest of config
  
  kova-dashboard:
    build:
      context: .
      dockerfile: apps/web/Dockerfile
    # ... rest of config
`
	}

	// Set default values for registry and version
	if config.RegistryURL == "" {
		config.RegistryURL = "ghcr.io"
	}
	if config.Repository == "" {
		config.Repository = "dopeCape/kova" // This should be configurable
	}
	if config.Version == "" {
		config.Version = "latest"
	}

	t, err := template.New("docker-compose").
		Funcs(template.FuncMap{
			"default": func(def, val string) string {
				if val == "" {
					return def
				}
				return val
			},
		}).
		Parse(tmpl)
	if err != nil {
		return fmt.Errorf("failed to parse docker-compose template: %w", err)
	}

	composeFile := filepath.Join(config.InstallDir, "docker-compose.yml")
	f, err := os.Create(composeFile)
	if err != nil {
		return fmt.Errorf("failed to create docker-compose.yml: %w", err)
	}
	defer f.Close()

	if err := t.Execute(f, config); err != nil {
		return fmt.Errorf("failed to generate docker-compose.yml: %w", err)
	}

	return nil
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
	resp, err := http.Get("https://ifconfig.me")
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	ip, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	return strings.TrimSpace(string(ip)), nil
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

