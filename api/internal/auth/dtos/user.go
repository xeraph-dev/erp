package dtos

import (
	"github.com/google/uuid"
)

type User struct {
	ID        uuid.UUID `json:"id" xml:"id"`
	Username  string    `json:"username" xml:"username"`
	Email     string    `json:"email" xml:"email"`
	FirstName string    `json:"first_name" xml:"first_name"`
	LastName  string    `json:"last_name" xml:"last_name"`
}

type UserLogin struct {
	Username string `json:"username" xml:"username"`
	Password string `json:"password" xml:"password"`
}

type UserRegister struct {
	Username string `json:"username" xml:"username"`
	Email    string `json:"email" xml:"email"`
	Password string `json:"password" xml:"password"`
}
