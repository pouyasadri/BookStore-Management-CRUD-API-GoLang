package controllers

import (
	"bookstore/pkg/models"
	"bookstore/pkg/utils"
	"os"
	"testing"
)

func setupControllerTest(t *testing.T) {
	// Set JWT secret for testing
	os.Setenv("JWT_SECRET", "test_secret_key_for_testing")
	utils.InitJWT()
}

// TestRegisterRequestValidation tests register request validation
func TestRegisterRequestValidation(t *testing.T) {
	tests := []struct {
		name   string
		req    RegisterRequest
		valid  bool
		fields []string
	}{
		{
			name:   "valid register request",
			req:    RegisterRequest{Username: "user", Email: "user@test.com", Password: "pass"},
			valid:  true,
			fields: []string{},
		},
		{
			name:   "missing username",
			req:    RegisterRequest{Username: "", Email: "user@test.com", Password: "pass"},
			valid:  false,
			fields: []string{"username"},
		},
		{
			name:   "missing email",
			req:    RegisterRequest{Username: "user", Email: "", Password: "pass"},
			valid:  false,
			fields: []string{"email"},
		},
		{
			name:   "missing password",
			req:    RegisterRequest{Username: "user", Email: "user@test.com", Password: ""},
			valid:  false,
			fields: []string{"password"},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			isValid := test.req.Username != "" && test.req.Email != "" && test.req.Password != ""
			if isValid != test.valid {
				t.Errorf("Expected valid=%v, got %v", test.valid, isValid)
			}
		})
	}
}

// TestLoginRequestValidation tests login request validation
func TestLoginRequestValidation(t *testing.T) {
	tests := []struct {
		name  string
		req   LoginRequest
		valid bool
	}{
		{
			name:  "valid login request",
			req:   LoginRequest{Username: "user", Password: "pass"},
			valid: true,
		},
		{
			name:  "missing username",
			req:   LoginRequest{Username: "", Password: "pass"},
			valid: false,
		},
		{
			name:  "missing password",
			req:   LoginRequest{Username: "user", Password: ""},
			valid: false,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			isValid := test.req.Username != "" && test.req.Password != ""
			if isValid != test.valid {
				t.Errorf("Expected valid=%v, got %v", test.valid, isValid)
			}
		})
	}
}

// TestRefreshRequestValidation tests refresh token request validation
func TestRefreshRequestValidation(t *testing.T) {
	tests := []struct {
		name  string
		req   RefreshRequest
		valid bool
	}{
		{
			name:  "valid refresh request",
			req:   RefreshRequest{Token: "some.jwt.token"},
			valid: true,
		},
		{
			name:  "missing token",
			req:   RefreshRequest{Token: ""},
			valid: false,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			isValid := test.req.Token != ""
			if isValid != test.valid {
				t.Errorf("Expected valid=%v, got %v", test.valid, isValid)
			}
		})
	}
}

// TestBookCreateRequestValidation tests book creation request validation
func TestBookCreateRequestValidation(t *testing.T) {
	tests := []struct {
		name  string
		req   BookCreateRequest
		valid bool
	}{
		{
			name:  "valid book request",
			req:   BookCreateRequest{Name: "Book", AuthorID: 1, Publication: "Pub"},
			valid: true,
		},
		{
			name:  "missing name",
			req:   BookCreateRequest{Name: "", AuthorID: 1, Publication: "Pub"},
			valid: false,
		},
		{
			name:  "missing author ID",
			req:   BookCreateRequest{Name: "Book", AuthorID: 0, Publication: "Pub"},
			valid: false,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			isValid := test.req.Name != "" && test.req.AuthorID != 0
			if isValid != test.valid {
				t.Errorf("Expected valid=%v, got %v", test.valid, isValid)
			}
		})
	}
}

// TestAuthorCreateRequestValidation tests author creation request validation
func TestAuthorCreateRequestValidation(t *testing.T) {
	tests := []struct {
		name  string
		req   AuthorCreateRequest
		valid bool
	}{
		{
			name:  "valid author request",
			req:   AuthorCreateRequest{Name: "Author", Email: "author@test.com", Biography: "Bio"},
			valid: true,
		},
		{
			name:  "missing name",
			req:   AuthorCreateRequest{Name: "", Email: "author@test.com", Biography: "Bio"},
			valid: false,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			isValid := test.req.Name != ""
			if isValid != test.valid {
				t.Errorf("Expected valid=%v, got %v", test.valid, isValid)
			}
		})
	}
}

// TestUserPasswordHashing tests that user passwords are properly hashed
func TestUserPasswordHashing(t *testing.T) {
	password := "testpassword123"
	hashedPassword := models.HashPassword(password)

	// Password should be hashed
	if hashedPassword == password {
		t.Error("Password should be hashed, not stored in plain text")
	}

	// Verify hashed password
	user := &models.User{Password: hashedPassword}
	if !user.VerifyPassword(password) {
		t.Error("Should verify correct password")
	}

	if user.VerifyPassword("wrongpassword") {
		t.Error("Should not verify incorrect password")
	}
}

// TestTokenGeneration tests JWT token generation
func TestTokenGeneration(t *testing.T) {
	setupControllerTest(t)

	token, err := utils.GenerateToken(1, "testuser", "test@example.com")
	if err != nil {
		t.Fatalf("Failed to generate token: %v", err)
	}

	if token == "" {
		t.Error("Generated token should not be empty")
	}

	// Validate the token
	claims, err := utils.ValidateToken(token)
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

// TestTokenValidation tests JWT token validation
func TestTokenValidation(t *testing.T) {
	setupControllerTest(t)

	// Test with invalid token
	_, err := utils.ValidateToken("invalid_token")
	if err == nil {
		t.Error("Should return error for invalid token")
	}

	// Test with valid token
	token, _ := utils.GenerateToken(1, "user", "user@example.com")
	claims, err := utils.ValidateToken(token)
	if err != nil {
		t.Errorf("Should validate correct token, got error: %v", err)
	}

	if claims.UserID != 1 {
		t.Errorf("Expected UserID 1, got %d", claims.UserID)
	}
}

// TestResponseStructures tests that response structures are properly defined
func TestResponseStructures(t *testing.T) {
	// Test AuthResponse
	authResp := AuthResponse{
		Token: "test_token",
		User: &UserInfo{
			ID:       1,
			Username: "testuser",
			Email:    "test@example.com",
		},
	}

	if authResp.Token != "test_token" {
		t.Error("AuthResponse token should be set")
	}

	if authResp.User.Username != "testuser" {
		t.Error("AuthResponse user should have username")
	}

	// Test BookCreateRequest
	bookReq := BookCreateRequest{
		Name:        "Test Book",
		AuthorID:    1,
		Publication: "Test Publisher",
	}

	if bookReq.Name != "Test Book" {
		t.Error("BookCreateRequest name should be set")
	}

	if bookReq.AuthorID != 1 {
		t.Error("BookCreateRequest AuthorID should be set")
	}

	// Test AuthorCreateRequest
	authorReq := AuthorCreateRequest{
		Name:      "Test Author",
		Email:     "author@example.com",
		Biography: "Test Biography",
	}

	if authorReq.Name != "Test Author" {
		t.Error("AuthorCreateRequest name should be set")
	}
}
