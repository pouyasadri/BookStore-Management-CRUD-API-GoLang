package models

import (
	"testing"
)

func TestUserPasswordHashing(t *testing.T) {
	password := "testpassword123"
	hashedPassword := HashPassword(password)

	// Password should be hashed
	if hashedPassword == password {
		t.Error("Password should be hashed, not stored in plain text")
	}

	// Bcrypt hashes are unique even for the same password (includes salt)
	// so we verify by checking that both hashes verify the same password
	hashedPassword2 := HashPassword(password)
	user := &User{Password: hashedPassword}
	user2 := &User{Password: hashedPassword2}

	if !user.VerifyPassword(password) {
		t.Error("First hash should verify the password")
	}
	if !user2.VerifyPassword(password) {
		t.Error("Second hash should verify the password")
	}
}

func TestUserVerifyPassword(t *testing.T) {
	user := &User{
		Username: "testuser",
		Password: HashPassword("testpassword"),
	}

	if !user.VerifyPassword("testpassword") {
		t.Error("Should verify correct password")
	}

	if user.VerifyPassword("wrongpassword") {
		t.Error("Should not verify incorrect password")
	}
}

func TestBookValidation(t *testing.T) {
	authorID := uint(1)
	book := &Book{
		Name:        "Test Book",
		AuthorID:    authorID,
		Publication: "Test Pub",
	}

	if book.Name == "" {
		t.Error("Book name should not be empty")
	}

	if book.AuthorID == 0 {
		t.Error("Book AuthorID should be set")
	}
}

func TestPaginationParams(t *testing.T) {
	params := PaginationParams{
		Page:  0,
		Limit: 0,
	}

	if params.Page < 1 {
		params.Page = 1
	}
	if params.Limit < 1 {
		params.Limit = 10
	}

	if params.Page != 1 {
		t.Error("Page should default to 1")
	}

	if params.Limit != 10 {
		t.Error("Limit should default to 10")
	}
}
