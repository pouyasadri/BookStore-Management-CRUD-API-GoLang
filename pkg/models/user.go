package models

import (
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type User struct {
	ID        uint   `gorm:"primaryKey" json:"id"`
	Username  string `json:"username" gorm:"type:varchar(255);uniqueIndex"`
	Email     string `json:"email" gorm:"type:varchar(191);uniqueIndex"`
	Password  string `json:"-" gorm:"type:varchar(72)"` // bcrypt hash is up to 60 chars, 72 bytes for margin
	CreatedAt int64  `json:"createdAt"`
	UpdatedAt int64  `json:"updatedAt"`
}

func (User) TableName() string {
	return "users"
}

// HashPassword creates a bcrypt hash of the password
func HashPassword(password string) string {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		// In case of error, return empty string. Caller should handle this.
		// This prevents password being stored in plaintext.
		return ""
	}
	return string(hash)
}

// VerifyPassword checks if the provided password matches the stored bcrypt hash
func (u *User) VerifyPassword(password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(password))
	return err == nil
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
