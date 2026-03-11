package routes

import (
	"bookstore/pkg/controllers"
	"bookstore/pkg/middleware"
	"net/http"

	"github.com/gorilla/mux"
)

var RegisterBookStoreRoutes = func(router *mux.Router) {
	// Apply global middleware
	router.Use(middleware.RecoveryMiddleware)
	router.Use(middleware.CORSMiddleware)
	router.Use(middleware.LoggingMiddleware)

	// Auth routes (no JWT required)
	router.HandleFunc("/auth/register", controllers.Register).Methods("POST")
	router.HandleFunc("/auth/login", controllers.Login).Methods("POST")

	// Health check endpoint
	router.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(200)
		w.Write([]byte(`{"status":"healthy"}`))
	}).Methods("GET")

	// Protected routes - require JWT token
	protected := router.PathPrefix("").Subrouter()
	protected.Use(middleware.JWTMiddleware)

	// Book routes (protected)
	protected.HandleFunc("/book/", controllers.CreateBook).Methods("POST")
	protected.HandleFunc("/book/", controllers.GetBook).Methods("GET")
	protected.HandleFunc("/book/{bookId}", controllers.GetBookById).Methods("GET")
	protected.HandleFunc("/book/{bookId}", controllers.UpdateBook).Methods("PUT")
	protected.HandleFunc("/book/{bookId}", controllers.DeleteBook).Methods("DELETE")

	// Author routes (protected)
	protected.HandleFunc("/author/", controllers.CreateAuthor).Methods("POST")
	protected.HandleFunc("/author/", controllers.GetAuthors).Methods("GET")
	protected.HandleFunc("/author/{authorId}", controllers.GetAuthorById).Methods("GET")
	protected.HandleFunc("/author/{authorId}", controllers.UpdateAuthor).Methods("PUT")
	protected.HandleFunc("/author/{authorId}", controllers.DeleteAuthor).Methods("DELETE")
}
