package utils

import (
	"os"
	"testing"
)

func setupJWTTest(t *testing.T) {
	// Set JWT secret for testing
	os.Setenv("JWT_SECRET", "test_secret_key_for_testing")
	InitJWT()
}

func TestGenerateAndValidateToken(t *testing.T) {
	setupJWTTest(t)

	token, err := GenerateToken(1, "testuser", "test@example.com")
	if err != nil {
		t.Fatalf("Failed to generate token: %v", err)
	}

	if token == "" {
		t.Error("Generated token should not be empty")
	}

	claims, err := ValidateToken(token)
	if err != nil {
		t.Fatalf("Failed to validate token: %v", err)
	}

	if claims.UserID != 1 {
		t.Errorf("Expected UserID 1, got %d", claims.UserID)
	}

	if claims.Username != "testuser" {
		t.Errorf("Expected username 'testuser', got '%s'", claims.Username)
	}

	if claims.Email != "test@example.com" {
		t.Errorf("Expected email 'test@example.com', got '%s'", claims.Email)
	}
}

func TestInvalidToken(t *testing.T) {
	setupJWTTest(t)

	_, err := ValidateToken("invalid_token")
	if err == nil {
		t.Error("Should return error for invalid token")
	}
}

func TestExpiredToken(t *testing.T) {
	setupJWTTest(t)

	// This test would require setting up an expired token
	// For now, just test that ValidateToken handles errors gracefully
	token := "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiIxMjM0NTY3ODkwIiwibmFtZSI6IkpvaG4gRG9lIiwiaWF0IjoxNTE2MjM5MDIyfQ.SflKxwRJSMeKKF2QT4fwpMeJf36POk6yJV_adQssw5c"

	_, err := ValidateToken(token)
	if err == nil {
		t.Error("Should return error for invalid/expired token")
	}
}
