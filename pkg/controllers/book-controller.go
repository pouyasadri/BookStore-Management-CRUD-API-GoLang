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

type BookCreateRequest struct {
	Name        string `json:"name" example:"The Go Programming Language"`
	AuthorID    uint   `json:"authorId" example:"1"`
	Publication string `json:"publication" example:"Addison-Wesley Professional"`
}

type BookUpdateRequest struct {
	Name        string `json:"name" example:"The Go Programming Language"`
	AuthorID    uint   `json:"authorId" example:"1"`
	Publication string `json:"publication" example:"Addison-Wesley Professional"`
}

// GetBook godoc
// @Summary Get books with pagination
// @Description Get all books with optional filtering by author or publication
// @Tags books
// @Accept json
// @Produce json
// @Param author query string false "Filter by author name"
// @Param publication query string false "Filter by publication"
// @Param page query integer false "Page number (default: 1)"
// @Param limit query integer false "Number of items per page (default: 10, max: 100)"
// @Success 200 {object} models.BookList
// @Failure 500 {object} utils.ErrorResponse
// @Router /books [get]
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

// GetBookById godoc
// @Summary Get a single book by ID
// @Description Retrieve a specific book by its ID
// @Tags books
// @Accept json
// @Produce json
// @Param bookId path integer true "Book ID"
// @Success 200 {object} models.Book
// @Failure 400 {object} utils.ErrorResponse
// @Failure 404 {object} utils.ErrorResponse
// @Failure 500 {object} utils.ErrorResponse
// @Router /books/{bookId} [get]
func GetBookById(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	bookId := vars["bookId"]
	ID, err := strconv.ParseInt(bookId, 10, 64)
	if err != nil {
		utils.RespondWithError(w, http.StatusBadRequest, "Invalid book ID", "")
		return
	}

	book, dbResult := models.GetBookById(ID)
	if dbResult.Error != nil {
		log.Printf("Database error fetching book (ID=%d): %v", ID, dbResult.Error)
		utils.RespondWithError(w, http.StatusInternalServerError, "Internal server error", "")
		return
	}

	if book.ID == 0 {
		utils.RespondWithError(w, http.StatusNotFound, "Book not found", "")
		return
	}

	utils.RespondWithSuccess(w, http.StatusOK, book)
}

// CreateBook godoc
// @Summary Create a new book
// @Description Create a new book with required fields
// @Tags books
// @Accept json
// @Produce json
// @Param input body BookCreateRequest true "Book data"
// @Success 201 {object} models.Book
// @Failure 400 {object} utils.ErrorResponse
// @Failure 500 {object} utils.ErrorResponse
// @Router /books [post]
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

	if newBook.AuthorID == 0 {
		utils.RespondWithError(w, http.StatusBadRequest, "Validation failed", "Book author ID is required")
		return
	}

	createdBook := newBook.CreateBook()
	utils.RespondWithSuccess(w, http.StatusCreated, createdBook)
}

// UpdateBook godoc
// @Summary Update an existing book
// @Description Update book fields (name, authorId, publication)
// @Tags books
// @Accept json
// @Produce json
// @Param bookId path integer true "Book ID"
// @Param input body BookUpdateRequest true "Updated book data"
// @Success 200 {object} models.Book
// @Failure 400 {object} utils.ErrorResponse
// @Failure 404 {object} utils.ErrorResponse
// @Failure 500 {object} utils.ErrorResponse
// @Router /books/{bookId} [put]
func UpdateBook(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	bookId := vars["bookId"]
	ID, err := strconv.ParseInt(bookId, 10, 64)
	if err != nil {
		utils.RespondWithError(w, http.StatusBadRequest, "Invalid book ID", "")
		return
	}

	updateData := &models.Book{}
	err = json.NewDecoder(r.Body).Decode(updateData)
	if err != nil {
		utils.RespondWithError(w, http.StatusBadRequest, "Invalid request payload", "")
		return
	}

	book, dbResult := models.GetBookById(ID)
	if dbResult.Error != nil {
		log.Printf("Database error fetching book (ID=%d): %v", ID, dbResult.Error)
		utils.RespondWithError(w, http.StatusInternalServerError, "Internal server error", "")
		return
	}

	if book.ID == 0 {
		utils.RespondWithError(w, http.StatusNotFound, "Book not found", "")
		return
	}

	if updateData.Name != "" {
		book.Name = updateData.Name
	}
	if updateData.AuthorID != 0 {
		book.AuthorID = updateData.AuthorID
	}
	if updateData.Publication != "" {
		book.Publication = updateData.Publication
	}

	updatedBook := book.UpdateBook()
	utils.RespondWithSuccess(w, http.StatusOK, updatedBook)
}

// DeleteBook godoc
// @Summary Delete a book
// @Description Delete a book by its ID
// @Tags books
// @Accept json
// @Produce json
// @Param bookId path integer true "Book ID"
// @Success 200 {object} utils.SuccessResponse
// @Failure 400 {object} utils.ErrorResponse
// @Failure 404 {object} utils.ErrorResponse
// @Failure 500 {object} utils.ErrorResponse
// @Router /books/{bookId} [delete]
func DeleteBook(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	bookId := vars["bookId"]
	ID, err := strconv.ParseInt(bookId, 10, 64)
	if err != nil {
		utils.RespondWithError(w, http.StatusBadRequest, "Invalid book ID", "")
		return
	}

	book, dbResult := models.GetBookById(ID)
	if dbResult.Error != nil {
		log.Printf("Database error fetching book (ID=%d): %v", ID, dbResult.Error)
		utils.RespondWithError(w, http.StatusInternalServerError, "Internal server error", "")
		return
	}

	if book.ID == 0 {
		utils.RespondWithError(w, http.StatusNotFound, "Book not found", "")
		return
	}

	models.DeleteBook(ID)
	utils.RespondWithMessage(w, http.StatusOK, "Book deleted successfully")
}
