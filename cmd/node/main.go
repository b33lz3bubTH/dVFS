package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/dvfs/storage-node/pkg/api"
	"github.com/dvfs/storage-node/pkg/config"
	"github.com/dvfs/storage-node/pkg/storage"
)

func main() {
	// Load configuration
	cfg := config.Load()

	// Initialize storage
	fileStorage, err := storage.NewFileStorage(cfg.StoragePath)
	if err != nil {
		log.Fatalf("Failed to initialize storage: %v", err)
	}

	// Initialize API router with instance ID
	router := api.NewRouter(fileStorage, cfg.InstanceID)

	// Setup HTTP server
	server := &http.Server{
		Addr:         fmt.Sprintf(":%d", cfg.Port),
		Handler:      router.Routes(),
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	// Start server in a goroutine
	go func() {
		log.Printf("🚀 Storage node starting on port %d", cfg.Port)
		log.Printf("🆔 Instance ID: %s", cfg.InstanceID)
		log.Printf("📁 Storage path: %s", cfg.StoragePath)
		log.Printf("🌐 API endpoints available at: http://localhost:%d/api/v1/files", cfg.Port)
		log.Printf("❤️  Health check: http://localhost:%d/health", cfg.Port)
		log.Printf("🔍 Instance info: http://localhost:%d/api/v1/instance", cfg.Port)
		
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Server failed to start: %v", err)
		}
	}()

	// Wait for interrupt signal to gracefully shutdown the server
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("🛑 Shutting down server...")

	// Give outstanding requests 30 seconds to complete
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		log.Fatalf("Server forced to shutdown: %v", err)
	}

	log.Println("✅ Server exited gracefully")
}
