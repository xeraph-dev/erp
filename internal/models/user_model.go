package models

import (
	"database/sql"
	"erp/internal/dtos"
	"erp/internal/vos"
	"fmt"

	"github.com/google/uuid"
)

type User struct {
	ID           uuid.UUID        `db:"id"`
	Username     vos.Username     `db:"username"`
	PasswordHash vos.PasswordHash `db:"password_hash"`
	Email        vos.Email        `db:"email"`
	FirstName    sql.NullString   `db:"first_name"`
	LastName     sql.NullString   `db:"last_name"`
}

func NewUserFromRegisterDTO(dto dtos.UserRegister) (model User, err error) {
	username, err := vos.NewUsername(dto.Username)
	if err != nil {
		err = fmt.Errorf("validating username: %w", err)
		return
	}

	passwordHash, err := vos.NewPasswordHash(dto.Password)
	if err != nil {
		err = fmt.Errorf("hashing password: %w", err)
		return
	}

	email, err := vos.NewEmail(dto.Email)
	if err != nil {
		err = fmt.Errorf("validating email: %w", err)
		return
	}

	model = User{
		Username:     username,
		PasswordHash: passwordHash,
		Email:        email,
	}

	return
}

func NewUserFromLoginDTO(dto dtos.UserLogin) (model User, err error) {
	username, err := vos.NewUsername(dto.Username)
	if err != nil {
		err = fmt.Errorf("validating username: %w", err)
		return
	}
	model = User{
		Username: username,
	}
	return
}

func (user User) PasswordMatches(raw string) bool {
	return user.PasswordHash.Matches(raw)
}

func (user User) DTO() dtos.User {
	return dtos.User{
		ID:        user.ID,
		Username:  string(user.Username),
		Email:     string(user.Email),
		FirstName: user.FirstName.String,
		LastName:  user.LastName.String,
	}
}
