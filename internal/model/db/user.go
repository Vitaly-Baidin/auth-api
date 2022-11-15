package models

import "github.com/golang-jwt/jwt/v4"

type User struct {
	ID       uint   `json:"id"`
	Login    string `json:"login"`
	Email    string `json:"email"`
	Phone    uint64 `json:"phone"`
	Password []byte `json:"-"`
}

type UserClaims struct {
	jwt.RegisteredClaims
	Email string `json:"email"`
	Type  string `json:"type"`
}

func NewUser(login string, email string, phone uint64, password string) *User {
	return &User{
		Login:    login,
		Email:    email,
		Phone:    phone,
		Password: []byte(password),
	}
}
