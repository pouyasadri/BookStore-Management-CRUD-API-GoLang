package controllers

import (
	"bookstore/pkg/models"
	"bookstore/pkg/utils"
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

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

func GetAuthorById(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	authorId := vars["authorId"]
	ID, err := strconv.ParseInt(authorId, 10, 64)
	if err != nil {
		utils.RespondWithError(w, http.StatusBadRequest, "Invalid author ID", err.Error())
		return
	}

	author, dbResult := models.GetAuthorById(ID)
	if dbResult.Error != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, "Database error", dbResult.Error.Error())
		return
	}

	if author.ID == 0 {
		utils.RespondWithError(w, http.StatusNotFound, "Author not found", "")
		return
	}

	utils.RespondWithSuccess(w, http.StatusOK, author)
}

func CreateAuthor(w http.ResponseWriter, r *http.Request) {
	newAuthor := &models.Author{}
	err := json.NewDecoder(r.Body).Decode(newAuthor)
	if err != nil {
		utils.RespondWithError(w, http.StatusBadRequest, "Invalid request payload", err.Error())
		return
	}

	if newAuthor.Name == "" {
		utils.RespondWithError(w, http.StatusBadRequest, "Validation failed", "Author name is required")
		return
	}

	createdAuthor := newAuthor.CreateAuthor()
	utils.RespondWithSuccess(w, http.StatusCreated, createdAuthor)
}

func UpdateAuthor(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	authorId := vars["authorId"]
	ID, err := strconv.ParseInt(authorId, 10, 64)
	if err != nil {
		utils.RespondWithError(w, http.StatusBadRequest, "Invalid author ID", err.Error())
		return
	}

	updateData := &models.Author{}
	err = json.NewDecoder(r.Body).Decode(updateData)
	if err != nil {
		utils.RespondWithError(w, http.StatusBadRequest, "Invalid request payload", err.Error())
		return
	}

	author, dbResult := models.GetAuthorById(ID)
	if dbResult.Error != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, "Database error", dbResult.Error.Error())
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

func DeleteAuthor(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	authorId := vars["authorId"]
	ID, err := strconv.ParseInt(authorId, 10, 64)
	if err != nil {
		utils.RespondWithError(w, http.StatusBadRequest, "Invalid author ID", err.Error())
		return
	}

	author, dbResult := models.GetAuthorById(ID)
	if dbResult.Error != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, "Database error", dbResult.Error.Error())
		return
	}

	if author.ID == 0 {
		utils.RespondWithError(w, http.StatusNotFound, "Author not found", "")
		return
	}

	models.DeleteAuthor(ID)
	utils.RespondWithMessage(w, http.StatusOK, "Author deleted successfully")
}
