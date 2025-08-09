package main

import (
	"log"
	"time"

	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/middleware/cors"
	"github.com/gofiber/fiber/v3/middleware/logger"
	"github.com/gofiber/fiber/v3/middleware/recover"
	"github.com/gofiber/fiber/v3/middleware/requestid"

	"github.com/dopeCape/kova/internal/api"
	"github.com/dopeCape/kova/internal/services"
	"github.com/dopeCape/kova/internal/store"
)

func GetApp(store store.Store, userService *services.UserService, authService *services.AuthService) *fiber.App {
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

			// Log error
			log.Printf("❌ Error: %v", err)

			return c.Status(code).JSON(fiber.Map{
				"error":     message,
				"timestamp": time.Now().Unix(),
				"path":      c.Path(),
				"method":    c.Method(),
			})
		},
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  60 * time.Second,
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
		// Check database connection
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
			"uptime":    time.Since(time.Now()).Seconds(),
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
				"projects": {
					"Coming soon...",
				},
				"deployments": {
					"Coming soon...",
				},
			},
		})
	})

	userHandler := api.NewUserHandler(userService)
	authHandler := api.NewAuthHandler(authService, userService)

	authHandler.RegisterRoutes(apiV1.Group("/auth"))

	apiV1.Post("/register", userHandler.CreateUser)

	userHandler.RegisterRoutes(apiV1.Group("/users", authHandler.RequireAuthMiddleware()))

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
