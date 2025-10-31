package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/dopeCape/kova/internal/config"
	"github.com/dopeCape/kova/internal/store/postgres/repository"

	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/middleware/cors"
	"github.com/gofiber/fiber/v3/middleware/logger"
	"github.com/gofiber/fiber/v3/middleware/recover"
	"github.com/gofiber/fiber/v3/middleware/requestid"

	"github.com/dopeCape/kova/internal/api"
	"github.com/dopeCape/kova/internal/services"
	"github.com/dopeCape/kova/internal/store"
)

func main() {
	cfg := config.Load()
	ctx := context.Background()
	db, store := repository.NewDefaultStore(ctx, cfg)
	defer db.Close()

	if err := store.Ping(ctx); err != nil {
		log.Fatal("❌ Database ping failed:", err)
	}
	log.Println("✅ Database health check passed")

	serverErrors := make(chan error, 1)

	// Initialize WebSocket hub
	wsHub := services.NewWebSocketHub()

	// Initialize services
	userService := services.NewUserService(store)
	githubService := services.NewGitHubService()
	accountService := services.NewAccountService(store, githubService)
	analyzerService := services.NewRepositoryAnalyzerService(githubService)
	authService := services.NewAuthService(userService, store, cfg.Auth.JWTSecret)

	// Initialize build service (needs store and account store)
	buildService := services.NewBuildService(store, store, wsHub)
	defer buildService.Shutdown()

	// Initialize project service with build service
	projectService := services.NewProjectService(store, buildService)

	log.Println("✅ Services initialized")

	app := GetApp(store, userService, accountService, projectService, authService, analyzerService, wsHub, cfg)

	go func() {
		port := ":" + cfg.Server.Port
		log.Printf("🚀 Server starting on port %s", cfg.Server.Port)
		log.Printf("📋 Available endpoints:")
		log.Printf("   Health: http://localhost%s/health", port)
		log.Printf("   API:    http://localhost%s/api/v1", port)
		log.Printf("   Users:  http://localhost%s/api/v1/users", port)
		log.Printf("   Accounts: http://localhost%s/api/v1/users/:id/accounts", port)
		log.Printf("   Projects: http://localhost%s/api/v1/users/:id/projects", port)
		log.Printf("   WebSocket: ws://localhost%s/api/v1/users/:id/projects/:projectId/ws", port)
		if err := app.Listen(port); err != nil {
			serverErrors <- err
		}
	}()

	shutdown := make(chan os.Signal, 1)
	signal.Notify(shutdown, os.Interrupt, syscall.SIGTERM)

	select {
	case err := <-serverErrors:
		log.Fatal("❌ Server failed to start:", err)
	case sig := <-shutdown:
		log.Printf("🔄 Shutting down server due to signal: %v", sig)

		// Shutdown build service first
		buildService.Shutdown()

		shutdownCtx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer cancel()

		if err := app.ShutdownWithContext(shutdownCtx); err != nil {
			log.Printf("❌ Server forced to shutdown: %v", err)
		}

		db.Close()
		log.Println("✅ Database connection closed")
		log.Println("✅ Server shutdown complete")
	}
}

func GetApp(store store.Store, userService *services.UserService, accountService *services.AccountService, projectService *services.ProjectService, authService *services.AuthService, analyzerService *services.RepositoryAnalyzerService, wsHub *services.WebSocketHub, cfg *config.Config) *fiber.App {
	app := fiber.New(fiber.Config{
		AppName:      "Kova",
		ServerHeader: "Kova",
		ErrorHandler: func(c fiber.Ctx, err error) error {
			code := fiber.StatusInternalServerError
			message := "Internal Server Error"

			if e, ok := err.(*fiber.Error); ok {
				code = e.Code
				message = e.Message
			}

			log.Printf("❌ Error: %v", err)

			return c.Status(code).JSON(fiber.Map{
				"error":     message,
				"timestamp": time.Now().Unix(),
				"path":      c.Path(),
				"method":    c.Method(),
			})
		},
		ReadTimeout:     10 * time.Second,
		WriteTimeout:    10 * time.Second,
		IdleTimeout:     60 * time.Second,
		ReadBufferSize:  16384, // 16KB
		WriteBufferSize: 16384, // 16KB
	})

	app.Use(requestid.New())
	app.Use(logger.New(logger.Config{
		Format: "[${time}] ${status} - ${method} ${path} - ${latency}\n",
	}))
	app.Use(recover.New())
	app.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Accept", "Authorization", "X-Requested-With"},
		AllowCredentials: false,
	}))

	// Health check endpoint
	app.Get("/health", func(c fiber.Ctx) error {
		if err := store.Ping(c.RequestCtx()); err != nil {
			return c.Status(503).JSON(fiber.Map{
				"status":    "unhealthy",
				"database":  "disconnected",
				"error":     err.Error(),
				"timestamp": time.Now().Unix(),
			})
		}

		return c.JSON(fiber.Map{
			"status":    "healthy",
			"database":  "connected",
			"version":   "1.0.0",
			"timestamp": time.Now().Unix(),
		})
	})

	app.Get("/", func(c fiber.Ctx) error {
		return c.JSON(fiber.Map{
			"name":        "Deployment Manager API",
			"version":     "1.0.0",
			"description": "Open source deployment manager API",
			"endpoints": map[string]string{
				"health": "/health",
				"api":    "/api/v1",
				"docs":   "/docs",
			},
			"timestamp": time.Now().Unix(),
		})
	})

	apiV1 := app.Group("/api/v1")

	apiV1.Get("/", func(c fiber.Ctx) error {
		return c.JSON(fiber.Map{
			"message": "Deployment Manager API v1",
			"endpoints": map[string][]string{
				"auth": {
					"POST /auth/login - Login user",
					"POST /auth/logout - Logout user",
					"POST /auth/refresh - Refresh access token",
					"GET /auth/me - Get current user (requires auth)",
				},
				"registration": {
					"POST /register - Register new user (public)",
				},
				"users": {
					"GET /users - List users (requires auth)",
					"GET /users/:id - Get user (requires auth)",
					"PUT /users/:id - Update user (requires auth)",
					"DELETE /users/:id - Delete user (requires auth)",
					"POST /users/:id/change-password - Change password (requires auth)",
					"GET /users/search?q=query - Search users (requires auth)",
				},
				"accounts": {
					"GET /users/:id/accounts - Get user's GitHub accounts (requires auth)",
					"POST /users/:id/accounts - Link new GitHub account (requires auth)",
					"GET /users/:id/accounts/:accountId/repositories - Get repositories (requires auth)",
				},
				"projects": {
					"GET /users/:id/projects - Get user's projects (requires auth)",
					"POST /users/:id/projects - Create new project (requires auth)",
					"GET /users/:id/projects/:projectId - Get project (requires auth)",
					"PUT /users/:id/projects/:projectId - Update project (requires auth)",
					"DELETE /users/:id/projects/:projectId - Delete project (requires auth)",
					"PUT /users/:id/projects/:projectId/archive - Archive project (requires auth)",
					"PUT /users/:id/projects/:projectId/activate - Activate project (requires auth)",
					"GET /users/:id/projects/search?q=query - Search projects (requires auth)",
					"GET /users/:id/projects/active - Get active projects (requires auth)",
					"WS /users/:id/projects/:projectId/ws - WebSocket deployment updates (requires auth)",
				},
				"deployments": {
					"Coming soon...",
				},
			},
		})
	})

	userHandler := api.NewUserHandler(userService)
	accountHandler := api.NewAccountHandler(accountService)
	projectHandler := api.NewProjectHandler(projectService)
	repositoryHandler := api.NewRepositoryHandler(accountService)
	analyzerHandler := api.NewAnalyzerHandler(analyzerService, accountService)
	authHandler := api.NewAuthHandler(authService, userService)
	authHandler.RegisterRoutes(apiV1.Group("/auth"))
	if cfg.Server.Env == config.DEVELOPMENT {
		apiV1.Post("/register", userHandler.CreateUser)
	}

	// IMPORTANT: Register WebSocket route BEFORE authenticated group
	// This ensures it doesn't go through auth middleware
	log.Println("📡 Registering WebSocket route (no auth)")
	apiV1.Get("/users/:id/projects/:projectId/ws", func(c fiber.Ctx) error {
		log.Printf("📡 WebSocket request received: User=%s, Project=%s", c.Params("id"), c.Params("projectId"))
		return wsHub.HandleWebSocket(c)
	})
	authenticatedGroup := apiV1.Group("/users", authHandler.RequireAuthMiddleware())
	userHandler.RegisterRoutes(authenticatedGroup)
	accountHandler.RegisterRoutes(authenticatedGroup)
	repositoryHandler.RegisterRoutes(authenticatedGroup)
	analyzerHandler.RegisterRoutes(authenticatedGroup)
	projectHandler.RegisterRoutes(authenticatedGroup)

	log.Println("✅ Routes registered")

	app.Use(func(c fiber.Ctx) error {
		return c.Status(404).JSON(fiber.Map{
			"error":   "Not Found",
			"message": "The requested endpoint does not exist",
			"path":    c.Path(),
			"method":  c.Method(),
		})
	})

	return app
}

func init() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)
}
