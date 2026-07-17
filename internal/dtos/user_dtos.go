package dtos

import "github.com/google/uuid"

type User struct {
	ID        uuid.UUID `json:"id" xml:"id"`
	Username  string    `json:"username" xml:"username"`
	Email     string    `json:"email" xml:"email"`
	FirstName string    `json:"first_name" xml:"first_name"`
	LastName  string    `json:"last_name" xml:"last_name"`
}

type UserLogin struct {
	Username string `json:"username" xml:"username"`
	Password string `json:"password_text" xml:"password_text"`
}

type UserRegister struct {
	Username string `json:"username" xml:"username"`
	Password string `json:"password_text" xml:"password_text"`
	Email    string `json:"email" xml:"email"`
}

var _ DTO = User{}
var _ DTO = UserLogin{}
var _ DTO = UserRegister{}

func (User) __internal()         {}
func (UserLogin) __internal()    {}
func (UserRegister) __internal() {}
