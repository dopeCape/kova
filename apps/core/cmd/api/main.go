package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/dopeCape/kova/internal/config"
	"github.com/dopeCape/kova/internal/services"
	"github.com/dopeCape/kova/internal/shared/database"
	"github.com/dopeCape/kova/internal/store/postgres/repository"
)

func main() {
	cfg := config.Load()
	ctx := context.Background()
	dbConfig := &database.Config{
		Host:     cfg.Database.Host,
		Port:     cfg.Database.Port,
		User:     cfg.Database.User,
		Password: cfg.Database.Password,
		Database: cfg.Database.Database,
		SSLMode:  cfg.Database.SSLMode,
	}

	db, err := database.NewConnection(ctx, dbConfig)
	if err != nil {
		log.Fatal("❌ Failed to connect to database:", err)
	}
	defer db.Close()
	store := repository.NewStore(db)
	if err := store.Ping(ctx); err != nil {
		log.Fatal("❌ Database ping failed:", err)
	}
	log.Println("✅ Database health check passed")
	log.Println("✅ Services initialized")
	serverErrors := make(chan error, 1)
	useService := services.NewUserService(store)
	authService := services.NewAuthService(useService, store, cfg.Auth.JWTSecret)
	app := GetApp(store, useService, authService)
	go func() {
		port := ":" + cfg.Server.Port
		log.Printf("🚀 Server starting on port %s", cfg.Server.Port)
		log.Printf("📋 Available endpoints:")
		log.Printf("   Health: http://localhost%s/health", port)
		log.Printf("   API:    http://localhost%s/api/v1", port)
		log.Printf("   Users:  http://localhost%s/api/v1/users", port)
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

func init() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)
}
