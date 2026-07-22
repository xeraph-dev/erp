package models

import (
	"context"
	"database/sql"
	"erp/internal/dtos"
	"erp/internal/middlewares"
	"erp/internal/vos"

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

func NewUserFromRegisterDTO(ctx context.Context, dto dtos.UserRegister) (model User, err error) {
	logger := middlewares.GetLogger(ctx)

	username, err := vos.NewUsername(dto.Username)
	if err != nil {
		logger.ErrorContext(ctx, "validating username", "error", err)
		return
	}

	passwordHash, err := vos.NewPasswordHash(dto.Password)
	if err != nil {
		logger.ErrorContext(ctx, "hashing password", "error", err)
		return
	}

	email, err := vos.NewEmail(dto.Email)
	if err != nil {
		logger.ErrorContext(ctx, "validating email", "error", err)
		return
	}

	model = User{
		Username:     username,
		PasswordHash: passwordHash,
		Email:        email,
	}

	return
}

func NewUserFromLoginDTO(ctx context.Context, dto dtos.UserLogin) (model User, err error) {
	logger := middlewares.GetLogger(ctx)

	username, err := vos.NewUsername(dto.Username)
	if err != nil {
		logger.ErrorContext(ctx, "validating username", "error", err)
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
