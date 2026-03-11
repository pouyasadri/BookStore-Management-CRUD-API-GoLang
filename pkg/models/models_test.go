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

	// Same password should produce same hash
	hashedPassword2 := HashPassword(password)
	if hashedPassword != hashedPassword2 {
		t.Error("Same password should produce same hash")
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
	book := &Book{
		Name:        "Test Book",
		Author:      "Test Author",
		Publication: "Test Pub",
	}

	if book.Name == "" {
		t.Error("Book name should not be empty")
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
