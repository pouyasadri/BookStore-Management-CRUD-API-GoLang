package main

import (
	_ "bookstore/docs"
	"bookstore/pkg/routes"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
	httpSwagger "github.com/swaggo/http-swagger/v2"
)

func main() {
	// Load .env file
	_ = godotenv.Load()

	port := os.Getenv("SERVER_PORT")
	if port == "" {
		port = "8080"
	}

	r := mux.NewRouter()
	routes.RegisterBookStoreRoutes(r)

	// Swagger UI
	r.PathPrefix("/swagger/").Handler(httpSwagger.WrapHandler)

	addr := fmt.Sprintf(":%s", port)
	log.Printf("Server starting on port %s", port)
	log.Printf("Swagger docs available at http://localhost:%s/swagger/", port)
	log.Fatal(http.ListenAndServe(addr, r))
}
