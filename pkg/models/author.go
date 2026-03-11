package models

import "gorm.io/gorm"

type Author struct {
	ID        uint   `gorm:"primaryKey" json:"id"`
	Name      string `json:"name" gorm:"type:varchar(255);index"`
	Email     string `json:"email" gorm:"type:varchar(255)"`
	Biography string `json:"biography"`
	Books     []Book `json:"books,omitempty" gorm:"foreignKey:AuthorID"`
	CreatedAt int64  `json:"createdAt"`
	UpdatedAt int64  `json:"updatedAt"`
}

func (Author) TableName() string {
	return "authors"
}

func (a *Author) CreateAuthor() *Author {
	db.Create(&a)
	return a
}

func GetAllAuthors() []Author {
	var authors []Author
	db.Preload("Books").Find(&authors)
	return authors
}

func GetAuthorById(id int64) (*Author, *gorm.DB) {
	var author Author
	dbInstance := db.Preload("Books").Where("id = ?", id).First(&author)
	return &author, dbInstance
}

func (a *Author) UpdateAuthor() *Author {
	db.Save(a)
	return a
}

func DeleteAuthor(id int64) {
	db.Delete(&Author{}, id)
}

type AuthorList struct {
	Data      []Author `json:"data"`
	Total     int64    `json:"total"`
	Page      int      `json:"page"`
	Limit     int      `json:"limit"`
	TotalPage int      `json:"totalPage"`
}

func GetAuthorsWithPagination(name string, params PaginationParams) *AuthorList {
	var authors []Author
	var total int64

	query := db
	if name != "" {
		query = query.Where("name LIKE ?", "%"+name+"%")
	}

	// Count the filtered results (before pagination)
	query.Model(&Author{}).Count(&total)

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
	query.Preload("Books").Offset(offset).Limit(params.Limit).Find(&authors)

	if authors == nil {
		authors = []Author{}
	}

	totalPage := int(total) / params.Limit
	if int(total)%params.Limit != 0 {
		totalPage++
	}

	return &AuthorList{
		Data:      authors,
		Total:     total,
		Page:      params.Page,
		Limit:     params.Limit,
		TotalPage: totalPage,
	}
}
