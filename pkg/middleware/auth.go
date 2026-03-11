package middleware

import (
	"bookstore/pkg/utils"
	"net/http"
	"strings"
)

func JWTMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			utils.RespondWithError(w, http.StatusUnauthorized, "Missing authorization header", "")
			return
		}

		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			utils.RespondWithError(w, http.StatusUnauthorized, "Invalid authorization header format", "Expected 'Bearer <token>'")
			return
		}

		tokenString := parts[1]
		claims, err := utils.ValidateToken(tokenString)
		if err != nil {
			utils.RespondWithError(w, http.StatusUnauthorized, "Invalid token", err.Error())
			return
		}

		// Store claims in context for handlers to access
		r.Header.Set("UserID", string(rune(claims.UserID)))
		r.Header.Set("Username", claims.Username)

		next.ServeHTTP(w, r)
	})
}
