package controllers

import (
	"bookstore/pkg/models"
	"bookstore/pkg/utils"
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

func GetBook(w http.ResponseWriter, r *http.Request) {
	// Get query parameters
	author := r.URL.Query().Get("author")
	publication := r.URL.Query().Get("publication")
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

	result := models.GetBooksWithPagination(author, publication, params)
	utils.RespondWithSuccess(w, http.StatusOK, result)
}

func GetBookById(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	bookId := vars["bookId"]
	ID, err := strconv.ParseInt(bookId, 10, 64)
	if err != nil {
		utils.RespondWithError(w, http.StatusBadRequest, "Invalid book ID", err.Error())
		return
	}

	book, dbResult := models.GetBookById(ID)
	if dbResult.Error != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, "Database error", dbResult.Error.Error())
		return
	}

	if book.ID == 0 {
		utils.RespondWithError(w, http.StatusNotFound, "Book not found", "")
		return
	}

	utils.RespondWithSuccess(w, http.StatusOK, book)
}

func CreateBook(w http.ResponseWriter, r *http.Request) {
	newBook := &models.Book{}
	err := json.NewDecoder(r.Body).Decode(newBook)
	if err != nil {
		utils.RespondWithError(w, http.StatusBadRequest, "Invalid request payload", err.Error())
		return
	}

	if newBook.Name == "" {
		utils.RespondWithError(w, http.StatusBadRequest, "Validation failed", "Book name is required")
		return
	}

	createdBook := newBook.CreateBook()
	utils.RespondWithSuccess(w, http.StatusCreated, createdBook)
}

func UpdateBook(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	bookId := vars["bookId"]
	ID, err := strconv.ParseInt(bookId, 10, 64)
	if err != nil {
		utils.RespondWithError(w, http.StatusBadRequest, "Invalid book ID", err.Error())
		return
	}

	updateData := &models.Book{}
	err = json.NewDecoder(r.Body).Decode(updateData)
	if err != nil {
		utils.RespondWithError(w, http.StatusBadRequest, "Invalid request payload", err.Error())
		return
	}

	book, dbResult := models.GetBookById(ID)
	if dbResult.Error != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, "Database error", dbResult.Error.Error())
		return
	}

	if book.ID == 0 {
		utils.RespondWithError(w, http.StatusNotFound, "Book not found", "")
		return
	}

	if updateData.Name != "" {
		book.Name = updateData.Name
	}
	if updateData.Author != "" {
		book.Author = updateData.Author
	}
	if updateData.Publication != "" {
		book.Publication = updateData.Publication
	}

	updatedBook := book.UpdateBook()
	utils.RespondWithSuccess(w, http.StatusOK, updatedBook)
}

func DeleteBook(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	bookId := vars["bookId"]
	ID, err := strconv.ParseInt(bookId, 10, 64)
	if err != nil {
		utils.RespondWithError(w, http.StatusBadRequest, "Invalid book ID", err.Error())
		return
	}

	book, dbResult := models.GetBookById(ID)
	if dbResult.Error != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, "Database error", dbResult.Error.Error())
		return
	}

	if book.ID == 0 {
		utils.RespondWithError(w, http.StatusNotFound, "Book not found", "")
		return
	}

	models.DeleteBook(ID)
	utils.RespondWithMessage(w, http.StatusOK, "Book deleted successfully")
}
