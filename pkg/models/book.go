package models

import (
	"bookstore/pkg/config"
	"gorm.io/gorm"
)

var db *gorm.DB

type Book struct {
	ID          uint   `gorm:"primaryKey" json:"id"`
	Name        string `json:"name"`
	Author      string `json:"author"`
	AuthorID    uint   `json:"authorId,omitempty" gorm:"index"`
	Publication string `json:"publication"`
	CreatedAt   int64  `json:"createdAt"`
	UpdatedAt   int64  `json:"updatedAt"`
}

func init() {
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
