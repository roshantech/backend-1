package model

import (
	"golang.org/x/crypto/bcrypt"
)

type User struct {
	Username   string `json:"username"`
	Password   string `json:"password"`
	Email      string `json:"email"`
	ProfilePic string `json:"profilepic`
	Active     bool   `json:"active`
	CreatedAt  string `json:"createdat`
	UpdatedAt  string `json:"updatedat`
}

func (u *User) TableName() string {
	return "users"
}

func (u *User) VerifyPassword(password string) error {
	return bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(password))
}



func HashPassword(password string) (string, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hashedPassword), nil
}