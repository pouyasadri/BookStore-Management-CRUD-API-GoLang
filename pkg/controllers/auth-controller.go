package controllers

import (
	"bookstore/pkg/models"
	"bookstore/pkg/utils"
	"encoding/json"
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
		utils.RespondWithError(w, http.StatusBadRequest, "Invalid request payload", err.Error())
		return
	}

	if req.Username == "" || req.Email == "" || req.Password == "" {
		utils.RespondWithError(w, http.StatusBadRequest, "Validation failed", "username, email, and password are required")
		return
	}

	// Check if user already exists
	existingUser, _ := models.GetUserByUsername(req.Username)
	if existingUser.ID != 0 {
		utils.RespondWithError(w, http.StatusConflict, "User already exists", "Username is already taken")
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
		utils.RespondWithError(w, http.StatusInternalServerError, "Failed to generate token", err.Error())
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
		utils.RespondWithError(w, http.StatusBadRequest, "Invalid request payload", err.Error())
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
		utils.RespondWithError(w, http.StatusInternalServerError, "Failed to generate token", err.Error())
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
