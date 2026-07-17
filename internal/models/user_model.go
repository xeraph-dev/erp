package models

import (
	"context"
	"database/sql"
	"erp/internal/dtos"
	"erp/internal/middlewares"
	"erp/internal/vos"
	"errors"

	"github.com/google/uuid"
)

var (
	ErrGeneratingPassword = errors.New("failed generating password")
)

type User struct {
	ID           uuid.UUID        `db:"id"`
	Username     vos.Username     `db:"username"`
	PasswordHash vos.PasswordHash `db:"password_hash"`
	Email        vos.Email        `db:"email"`
	FirstName    sql.NullString   `db:"first_name"`
	LastName     sql.NullString   `db:"last_name"`
}

func NewUserFromRegisterDTO(ctx context.Context, dto dtos.UserRegister) (model User, err error) {
	logger := middlewares.GetLogger(ctx)

	username, err := vos.NewUsername(dto.Username)
	if err != nil {
		return
	}

	passwordHash, err := vos.NewPasswordHash(ctx, dto.Password)
	if err != nil {
		logger.ErrorContext(ctx, "hashing password", "error", err)
		err = ErrGeneratingPassword
		return
	}

	email, err := vos.NewEmail(dto.Email)
	if err != nil {
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
		return
	}
	model = User{
		Username: username,
	}
	return
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
