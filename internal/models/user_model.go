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
	Username     string           `db:"username"`
	PasswordHash vos.PasswordHash `db:"password_hash"`
	Email        string           `db:"email"`
	FirstName    sql.NullString   `db:"first_name"`
	LastName     sql.NullString   `db:"last_name"`
}

func NewUserFromRegisterDTO(ctx context.Context, dto dtos.UserRegister) (model User, err error) {
	logger := middlewares.GetLogger(ctx)

	passwordHash, err := vos.NewPasswordHash(dto.Password)
	if err != nil {
		logger.ErrorContext(ctx, ErrGeneratingPassword.Error(), "error", err)
		err = ErrGeneratingPassword
		return
	}

	model = User{
		Username:     dto.Username,
		PasswordHash: passwordHash,
		Email:        dto.Email,
	}

	return
}

func NewUserFromLoginDTO(dto dtos.UserLogin) (model User) {
	return User{
		Username: dto.Username,
	}
}

func (user User) DTO() dtos.User {
	return dtos.User{
		ID:        user.ID,
		Username:  user.Username,
		Email:     user.Email,
		FirstName: user.FirstName.String,
		LastName:  user.LastName.String,
	}
}
