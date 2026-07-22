package dtos

import "time"

type UserLogin struct {
	Username string `json:"username" xml:"username"`
	Password string `json:"password_text" xml:"password_text"`
}

type UserRegister struct {
	Username string `json:"username" xml:"username"`
	Password string `json:"password_text" xml:"password_text"`
	Email    string `json:"email" xml:"email"`
}

type TokenPair struct {
	AccessToken           string    `json:"access_token" xml:"access_token"`
	RefreshToken          string    `json:"refresh_token" xml:"refresh_token"`
	AccessTokenExpiresAt  time.Time `json:"-" xml:"-"`
	RefreshTokenExpiresAt time.Time `json:"-" xml:"-"`
}

var _ DTO = (*UserLogin)(nil)
var _ DTO = (*UserRegister)(nil)
var _ DTO = (*TokenPair)(nil)

func (UserLogin) __internal()     {}
func (UserRegister) __internal()  {}
func (TokenPair) __internal() {}
