package utils

import (
	"log"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

const DefaultJWTSecret = "your_secret_key_here_change_in_production"

var jwtSecret string

// InitJWT validates and initializes the JWT secret at startup
// Fails fast if JWT_SECRET is missing or using the insecure default
func InitJWT() {
	secret := os.Getenv("JWT_SECRET")

	if secret == "" {
		log.Fatal("FATAL: JWT_SECRET environment variable is not set. Aborting startup. " +
			"Please set JWT_SECRET to a secure random string.")
	}

	if secret == DefaultJWTSecret {
		log.Fatal("FATAL: JWT_SECRET is set to the default insecure value. " +
			"Please change it to a secure random string in your environment configuration.")
	}

	jwtSecret = secret
	log.Println("JWT secret initialized successfully")
}

type Claims struct {
	UserID   uint   `json:"userId"`
	Username string `json:"username"`
	Email    string `json:"email"`
	jwt.RegisteredClaims
}

func GenerateToken(userID uint, username, email string) (string, error) {
	claims := &Claims{
		UserID:   userID,
		Username: username,
		Email:    email,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(jwtSecret))
}

func ValidateToken(tokenString string) (*Claims, error) {
	claims := &Claims{}
	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		return []byte(jwtSecret), nil
	})

	if err != nil {
		return nil, err
	}

	if !token.Valid {
		return nil, jwt.ErrSignatureInvalid
	}

	return claims, nil
}
