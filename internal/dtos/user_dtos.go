package dtos

import "github.com/google/uuid"

type User struct {
	ID        uuid.UUID `json:"id"`
	Username  string    `json:"username"`
	Email     string    `json:"email"`
	FirstName string    `json:"first_name"`
	LastName  string    `json:"last_name"`
}

type UserLogin struct {
	Username string `json:"username"`
	Password string `json:"password_text"`
}

type UserRegister struct {
	Username string `json:"username"`
	Password string `json:"password_text"`
	Email    string `json:"email"`
}
