package controllers

import (
	"bookstore/pkg/models"
	"bookstore/pkg/utils"
	"encoding/json"
	"log"
	"net/http"
)

type RegisterRequest struct {
	Username string `json:"username" example:"john_doe"`
	Email    string `json:"email" example:"john@example.com"`
	Password string `json:"password" example:"securepassword123"`
}

type LoginRequest struct {
	Username string `json:"username" example:"john_doe"`
	Password string `json:"password" example:"securepassword123"`
}

type RefreshRequest struct {
	Token string `json:"token" example:"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."`
}

type AuthResponse struct {
	Token string    `json:"token"`
	User  *UserInfo `json:"user"`
}

type UserInfo struct {
	ID       uint   `json:"id"`
	Username string `json:"username"`
	Email    string `json:"email"`
}

// Register godoc
// @Summary User registration
// @Description Create a new user account
// @Tags auth
// @Accept json
// @Produce json
// @Param input body RegisterRequest true "User registration data"
// @Success 201 {object} AuthResponse
// @Failure 400 {object} utils.ErrorResponse
// @Failure 409 {object} utils.ErrorResponse
// @Router /auth/register [post]
func Register(w http.ResponseWriter, r *http.Request) {
	var req RegisterRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		utils.RespondWithError(w, http.StatusBadRequest, "Invalid request payload", "")
		return
	}

	if req.Username == "" || req.Email == "" || req.Password == "" {
		utils.RespondWithError(w, http.StatusBadRequest, "Validation failed", "username, email, and password are required")
		return
	}

	// Check if username already exists
	existingUser, _ := models.GetUserByUsername(req.Username)
	if existingUser.ID != 0 {
		utils.RespondWithError(w, http.StatusConflict, "User already exists", "Username is already taken")
		return
	}

	// Check if email already exists
	existingEmail, _ := models.GetUserByEmail(req.Email)
	if existingEmail.ID != 0 {
		utils.RespondWithError(w, http.StatusConflict, "User already exists", "Email is already registered")
		return
	}

	user := &models.User{
		Username: req.Username,
		Email:    req.Email,
		Password: req.Password,
	}

	user.CreateUser()

	token, err := utils.GenerateToken(user.ID, user.Username, user.Email)
	if err != nil {
		log.Printf("Failed to generate token for user %d: %v", user.ID, err)
		utils.RespondWithError(w, http.StatusInternalServerError, "Internal server error", "")
		return
	}

	response := AuthResponse{
		Token: token,
		User: &UserInfo{
			ID:       user.ID,
			Username: user.Username,
			Email:    user.Email,
		},
	}

	utils.RespondWithSuccess(w, http.StatusCreated, response)
}

// Login godoc
// @Summary User login
// @Description Login with username and password to get JWT token
// @Tags auth
// @Accept json
// @Produce json
// @Param input body LoginRequest true "User login credentials"
// @Success 200 {object} AuthResponse
// @Failure 400 {object} utils.ErrorResponse
// @Failure 401 {object} utils.ErrorResponse
// @Router /auth/login [post]
func Login(w http.ResponseWriter, r *http.Request) {
	var req LoginRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		utils.RespondWithError(w, http.StatusBadRequest, "Invalid request payload", "")
		return
	}

	if req.Username == "" || req.Password == "" {
		utils.RespondWithError(w, http.StatusBadRequest, "Validation failed", "username and password are required")
		return
	}

	user, dbResult := models.GetUserByUsername(req.Username)
	if dbResult.Error != nil || user.ID == 0 {
		utils.RespondWithError(w, http.StatusUnauthorized, "Invalid credentials", "")
		return
	}

	if !user.VerifyPassword(req.Password) {
		utils.RespondWithError(w, http.StatusUnauthorized, "Invalid credentials", "")
		return
	}

	token, err := utils.GenerateToken(user.ID, user.Username, user.Email)
	if err != nil {
		log.Printf("Failed to generate token for user %d: %v", user.ID, err)
		utils.RespondWithError(w, http.StatusInternalServerError, "Internal server error", "")
		return
	}

	response := AuthResponse{
		Token: token,
		User: &UserInfo{
			ID:       user.ID,
			Username: user.Username,
			Email:    user.Email,
		},
	}

	utils.RespondWithSuccess(w, http.StatusOK, response)
}

// Refresh godoc
// @Summary Refresh JWT token
// @Description Generate a new JWT token from an existing valid token
// @Tags auth
// @Accept json
// @Produce json
// @Param input body RefreshRequest true "Current JWT token"
// @Success 200 {object} AuthResponse
// @Failure 400 {object} utils.ErrorResponse
// @Failure 401 {object} utils.ErrorResponse
// @Failure 500 {object} utils.ErrorResponse
// @Router /auth/refresh [post]
func Refresh(w http.ResponseWriter, r *http.Request) {
	var req RefreshRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		utils.RespondWithError(w, http.StatusBadRequest, "Invalid request payload", "")
		return
	}

	if req.Token == "" {
		utils.RespondWithError(w, http.StatusBadRequest, "Validation failed", "token is required")
		return
	}

	// Validate the token
	claims, err := utils.ValidateToken(req.Token)
	if err != nil {
		utils.RespondWithError(w, http.StatusUnauthorized, "Invalid or expired token", "")
		return
	}

	// Generate a new token with the same claims
	newToken, err := utils.GenerateToken(claims.UserID, claims.Username, claims.Email)
	if err != nil {
		log.Printf("Failed to generate refreshed token for user %d: %v", claims.UserID, err)
		utils.RespondWithError(w, http.StatusInternalServerError, "Internal server error", "")
		return
	}

	response := AuthResponse{
		Token: newToken,
		User: &UserInfo{
			ID:       claims.UserID,
			Username: claims.Username,
			Email:    claims.Email,
		},
	}

	utils.RespondWithSuccess(w, http.StatusOK, response)
}
