package controllers

import (
	"bookstore/pkg/models"
	"bookstore/pkg/utils"
	"encoding/json"
	"log"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

type AuthorCreateRequest struct {
	Name      string `json:"name" example:"Robert C. Martin"`
	Email     string `json:"email" example:"robert@example.com"`
	Biography string `json:"biography" example:"Author of Clean Code"`
}

type AuthorUpdateRequest struct {
	Name      string `json:"name" example:"Robert C. Martin"`
	Email     string `json:"email" example:"robert@example.com"`
	Biography string `json:"biography" example:"Author of Clean Code"`
}

// GetAuthors godoc
// @Summary Get authors with pagination
// @Description Get all authors with optional filtering by name
// @Tags authors
// @Accept json
// @Produce json
// @Param name query string false "Filter by author name"
// @Param page query integer false "Page number (default: 1)"
// @Param limit query integer false "Number of items per page (default: 10, max: 100)"
// @Success 200 {object} models.AuthorList
// @Failure 500 {object} utils.ErrorResponse
// @Router /authors [get]
func GetAuthors(w http.ResponseWriter, r *http.Request) {
	name := r.URL.Query().Get("name")
	page := r.URL.Query().Get("page")
	limit := r.URL.Query().Get("limit")

	pageNum := 1
	if page != "" {
		p, err := strconv.Atoi(page)
		if err == nil && p > 0 {
			pageNum = p
		}
	}

	limitNum := 10
	if limit != "" {
		l, err := strconv.Atoi(limit)
		if err == nil && l > 0 {
			limitNum = l
		}
	}

	params := models.PaginationParams{
		Page:  pageNum,
		Limit: limitNum,
	}

	result := models.GetAuthorsWithPagination(name, params)
	utils.RespondWithSuccess(w, http.StatusOK, result)
}

// GetAuthorById godoc
// @Summary Get a single author by ID
// @Description Retrieve a specific author by their ID
// @Tags authors
// @Accept json
// @Produce json
// @Param authorId path integer true "Author ID"
// @Success 200 {object} models.Author
// @Failure 400 {object} utils.ErrorResponse
// @Failure 404 {object} utils.ErrorResponse
// @Failure 500 {object} utils.ErrorResponse
// @Router /authors/{authorId} [get]
func GetAuthorById(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	authorId := vars["authorId"]
	ID, err := strconv.ParseInt(authorId, 10, 64)
	if err != nil {
		utils.RespondWithError(w, http.StatusBadRequest, "Invalid author ID", "")
		return
	}

	author, dbResult := models.GetAuthorById(ID)
	if dbResult.Error != nil {
		log.Printf("Database error fetching author (ID=%d): %v", ID, dbResult.Error)
		utils.RespondWithError(w, http.StatusInternalServerError, "Internal server error", "")
		return
	}

	if author.ID == 0 {
		utils.RespondWithError(w, http.StatusNotFound, "Author not found", "")
		return
	}

	utils.RespondWithSuccess(w, http.StatusOK, author)
}

// CreateAuthor godoc
// @Summary Create a new author
// @Description Create a new author with required fields
// @Tags authors
// @Accept json
// @Produce json
// @Param input body AuthorCreateRequest true "Author data"
// @Success 201 {object} models.Author
// @Failure 400 {object} utils.ErrorResponse
// @Failure 500 {object} utils.ErrorResponse
// @Router /authors [post]
func CreateAuthor(w http.ResponseWriter, r *http.Request) {
	newAuthor := &models.Author{}
	err := json.NewDecoder(r.Body).Decode(newAuthor)
	if err != nil {
		utils.RespondWithError(w, http.StatusBadRequest, "Invalid request payload", "")
		return
	}

	if newAuthor.Name == "" {
		utils.RespondWithError(w, http.StatusBadRequest, "Validation failed", "Author name is required")
		return
	}

	createdAuthor := newAuthor.CreateAuthor()
	utils.RespondWithSuccess(w, http.StatusCreated, createdAuthor)
}

// UpdateAuthor godoc
// @Summary Update an existing author
// @Description Update author fields (name, email, biography)
// @Tags authors
// @Accept json
// @Produce json
// @Param authorId path integer true "Author ID"
// @Param input body AuthorUpdateRequest true "Updated author data"
// @Success 200 {object} models.Author
// @Failure 400 {object} utils.ErrorResponse
// @Failure 404 {object} utils.ErrorResponse
// @Failure 500 {object} utils.ErrorResponse
// @Router /authors/{authorId} [put]
func UpdateAuthor(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	authorId := vars["authorId"]
	ID, err := strconv.ParseInt(authorId, 10, 64)
	if err != nil {
		utils.RespondWithError(w, http.StatusBadRequest, "Invalid author ID", "")
		return
	}

	updateData := &models.Author{}
	err = json.NewDecoder(r.Body).Decode(updateData)
	if err != nil {
		utils.RespondWithError(w, http.StatusBadRequest, "Invalid request payload", "")
		return
	}

	author, dbResult := models.GetAuthorById(ID)
	if dbResult.Error != nil {
		log.Printf("Database error fetching author (ID=%d): %v", ID, dbResult.Error)
		utils.RespondWithError(w, http.StatusInternalServerError, "Internal server error", "")
		return
	}

	if author.ID == 0 {
		utils.RespondWithError(w, http.StatusNotFound, "Author not found", "")
		return
	}

	if updateData.Name != "" {
		author.Name = updateData.Name
	}
	if updateData.Email != "" {
		author.Email = updateData.Email
	}
	if updateData.Biography != "" {
		author.Biography = updateData.Biography
	}

	updatedAuthor := author.UpdateAuthor()
	utils.RespondWithSuccess(w, http.StatusOK, updatedAuthor)
}

// DeleteAuthor godoc
// @Summary Delete an author
// @Description Delete an author by their ID
// @Tags authors
// @Accept json
// @Produce json
// @Param authorId path integer true "Author ID"
// @Success 200 {object} utils.SuccessResponse
// @Failure 400 {object} utils.ErrorResponse
// @Failure 404 {object} utils.ErrorResponse
// @Failure 500 {object} utils.ErrorResponse
// @Router /authors/{authorId} [delete]
func DeleteAuthor(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	authorId := vars["authorId"]
	ID, err := strconv.ParseInt(authorId, 10, 64)
	if err != nil {
		utils.RespondWithError(w, http.StatusBadRequest, "Invalid author ID", "")
		return
	}

	author, dbResult := models.GetAuthorById(ID)
	if dbResult.Error != nil {
		log.Printf("Database error fetching author (ID=%d): %v", ID, dbResult.Error)
		utils.RespondWithError(w, http.StatusInternalServerError, "Internal server error", "")
		return
	}

	if author.ID == 0 {
		utils.RespondWithError(w, http.StatusNotFound, "Author not found", "")
		return
	}

	models.DeleteAuthor(ID)
	utils.RespondWithMessage(w, http.StatusOK, "Author deleted successfully")
}
