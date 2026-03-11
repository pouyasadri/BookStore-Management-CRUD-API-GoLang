package models

import (
	"bookstore/pkg/config"
	"flag"
	"gorm.io/gorm"
	"os"
	"strings"
)

var db *gorm.DB

type Book struct {
	ID          uint   `gorm:"primaryKey" json:"id"`
	Name        string `json:"name"`
	Author      string `json:"author"`
	AuthorID    *uint  `json:"authorId,omitempty" gorm:"index"`
	Publication string `json:"publication"`
	CreatedAt   int64  `json:"createdAt"`
	UpdatedAt   int64  `json:"updatedAt"`
}

func init() {
	// When running `go test`, avoid connecting to a real database. Tests
	// should be isolated and not require a MySQL instance. Detect test
	// execution by checking for the testing flags (present when `go test`
	// is running) and skip DB initialization in that case. Checking
	// os.Args for any -test.* flags is more reliable during package init.
	for _, a := range os.Args {
		if strings.HasPrefix(a, "-test.") {
			return
		}
	}
	if flag.Lookup("test.v") != nil {
		return
	}

	config.Connect()
	db = config.GetDB()
	db.AutoMigrate(&User{}, &Author{}, &Book{})
}

func (b *Book) CreateBook() *Book {
	db.Create(&b)
	return b
}

func GetAllBooks() []Book {
	var books []Book
	db.Find(&books)
	return books
}

type PaginationParams struct {
	Page  int
	Limit int
}

type BookList struct {
	Data      []Book `json:"data"`
	Total     int64  `json:"total"`
	Page      int    `json:"page"`
	Limit     int    `json:"limit"`
	TotalPage int    `json:"totalPage"`
}

func GetBooksWithPagination(author, publication string, params PaginationParams) *BookList {
	var books []Book
	var total int64

	query := db
	if author != "" {
		query = query.Where("author LIKE ?", "%"+author+"%")
	}
	if publication != "" {
		query = query.Where("publication LIKE ?", "%"+publication+"%")
	}

	query.Model(&Book{}).Count(&total)

	if params.Page < 1 {
		params.Page = 1
	}
	if params.Limit < 1 {
		params.Limit = 10
	}
	if params.Limit > 100 {
		params.Limit = 100
	}

	offset := (params.Page - 1) * params.Limit
	query.Offset(offset).Limit(params.Limit).Find(&books)

	if books == nil {
		books = []Book{}
	}

	totalPage := int(total) / params.Limit
	if int(total)%params.Limit != 0 {
		totalPage++
	}

	return &BookList{
		Data:      books,
		Total:     total,
		Page:      params.Page,
		Limit:     params.Limit,
		TotalPage: totalPage,
	}
}

func GetBookById(id int64) (*Book, *gorm.DB) {
	var book Book
	dbInstance := db.Where("id = ?", id).First(&book)
	return &book, dbInstance
}

func (b *Book) UpdateBook() *Book {
	db.Save(b)
	return b
}

func DeleteBook(id int64) Book {
	var book Book
	db.Delete(&book, id)
	return book
}
