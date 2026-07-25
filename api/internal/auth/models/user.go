package models

import (
	"database/sql"
	"erp/internal/auth/dtos"
	"erp/internal/auth/vos"
	"fmt"

	"github.com/google/uuid"
)

type User struct {
	ID           uuid.UUID        `db:"id"`
	Username     vos.Username     `db:"username"`
	Email        vos.Email        `db:"email"`
	PasswordHash vos.PasswordHash `db:"password_hash"`
	FirstName    sql.NullString   `db:"first_name"`
	LastName     sql.NullString   `db:"last_name"`
}

func NewUserFromRegisterDTO(dto dtos.UserRegister) (model User, err error) {
	const op = "models/NewUserFromRegisterDTO"

	username, err := vos.NewUsername(dto.Username)
	if err != nil {
		return model, fmt.Errorf("%s: %w", op, err)
	}

	email, err := vos.NewEmail(dto.Email)
	if err != nil {
		return model, fmt.Errorf("%s: %w", op, err)
	}

	passwordHash, err := vos.NewPasswordHash(dto.Password)
	if err != nil {
		return model, fmt.Errorf("%s: %w", op, err)
	}

	return User{
		Username:     username,
		Email:        email,
		PasswordHash: passwordHash,
	}, nil
}

func NewUserFromLoginDTO(dto dtos.UserLogin) (model User, err error) {
	const op = "models/NewUserFromLoginDTO"

	username, err := vos.NewUsername(dto.Username)
	if err != nil {
		return model, fmt.Errorf("%s: %w", op, err)
	}

	passwordHash, err := vos.NewPasswordHash(dto.Password)
	if err != nil {
		return model, fmt.Errorf("%s: %w", op, err)
	}

	return User{
		Username:     username,
		PasswordHash: passwordHash,
	}, nil
}

func (user User) PasswordMatches(raw string) bool {
	return user.PasswordHash.Matches(raw)
}

func (user User) DTO() dtos.User {
	return dtos.User{
		ID:        user.ID,
		Username:  user.Username.String(),
		Email:     user.Email.String(),
		FirstName: user.FirstName.String,
		LastName:  user.LastName.String,
	}
}
