package main

import (
	_ "bookstore/docs"
	"bookstore/pkg/routes"
	"bookstore/pkg/utils"
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
	httpSwagger "github.com/swaggo/http-swagger/v2"
)

func main() {
	// Load .env file
	_ = godotenv.Load()

	// Initialize and validate JWT secret (fail-fast if missing or insecure)
	utils.InitJWT()

	port := os.Getenv("SERVER_PORT")
	if port == "" {
		port = "8080"
	}

	r := mux.NewRouter()
	routes.RegisterBookStoreRoutes(r)

	// Swagger UI
	r.PathPrefix("/swagger/").Handler(httpSwagger.WrapHandler)

	addr := fmt.Sprintf(":%s", port)

	// Create HTTP server
	server := &http.Server{
		Addr:         addr,
		Handler:      r,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	// Start server in a goroutine
	go func() {
		log.Printf("Server starting on port %s", port)
		log.Printf("Swagger docs available at http://localhost:%s/swagger/", port)
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Server error: %v", err)
		}
	}()

	// Setup graceful shutdown on SIGTERM or SIGINT
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGTERM, syscall.SIGINT)

	// Block until a signal is received
	sig := <-sigChan
	log.Printf("Received signal: %v. Shutting down gracefully...", sig)

	// Graceful shutdown with 30 second timeout
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		log.Printf("Server forced to shutdown: %v", err)
		os.Exit(1)
	}

	log.Println("Server stopped")
}
