package models

import (
	"crypto/sha256"
	"fmt"
	"gorm.io/gorm"
)

type User struct {
	ID        uint   `gorm:"primaryKey" json:"id"`
	Username  string `json:"username" gorm:"uniqueIndex"`
	Email     string `json:"email" gorm:"uniqueIndex"`
	Password  string `json:"-"` // Never expose password in JSON
	CreatedAt int64  `json:"createdAt"`
	UpdatedAt int64  `json:"updatedAt"`
}

func (User) TableName() string {
	return "users"
}

// HashPassword creates a SHA256 hash of the password
func HashPassword(password string) string {
	hash := sha256.Sum256([]byte(password))
	return fmt.Sprintf("%x", hash)
}

// VerifyPassword checks if the provided password matches the stored hash
func (u *User) VerifyPassword(password string) bool {
	return u.Password == HashPassword(password)
}

func (u *User) CreateUser() *User {
	u.Password = HashPassword(u.Password)
	db.Create(&u)
	return u
}

func GetUserByUsername(username string) (*User, *gorm.DB) {
	var user User
	dbInstance := db.Where("username = ?", username).First(&user)
	return &user, dbInstance
}

func GetUserByEmail(email string) (*User, *gorm.DB) {
	var user User
	dbInstance := db.Where("email = ?", email).First(&user)
	return &user, dbInstance
}
