package services

import (
	"context"
	"erp/internal/auth/dtos"
	"erp/internal/auth/models"
	"erp/internal/auth/stores"
	"erp/internal/auth/vos"
	"erp/internal/shared/db"
	"erp/internal/shared/tokens"
	"errors"
	"fmt"
	"time"
)

var dummyPasswordHash, _ = vos.NewPasswordHash("Dummy-Password-For-Timing-1!")

var (
	ErrInvalidCredentials   = errors.New("invalid credentials")   // -> 401
	ErrRecordConflict       = errors.New("record already exists") // -> 409
	ErrValidationFailed     = errors.New("validation failed")     // -> 422
	ErrUserPasswordNotMatch = errors.New("user password does not match")
	ErrRefreshTokenReused   = errors.New("refresh token reuse detected")
	ErrRefreshTokenExpired  = errors.New("refresh token expired")
)

type Auth interface {
	Register(ctx context.Context, in dtos.UserRegister) (out tokens.Pair, err error)
	Login(ctx context.Context, in dtos.UserLogin) (out tokens.Pair, err error)
	Logout(ctx context.Context, refreshToken string) (err error)
	Refresh(ctx context.Context, refreshToken string) (out tokens.Pair, err error)
}

type authImpl struct {
	db      db.TxBeginner
	token   tokens.Tokens
	user    stores.User
	refresh stores.RefreshToken
}

var _ Auth = (*authImpl)(nil)

func NewAuth(db db.TxBeginner, token tokens.Tokens, user stores.User, refresh stores.RefreshToken) Auth {
	return authImpl{db: db, token: token, user: user, refresh: refresh}
}

func (service authImpl) Register(ctx context.Context, in dtos.UserRegister) (out tokens.Pair, err error) {
	model, err := models.NewUserFromRegisterDTO(in)
	if err != nil {
		return out, fmt.Errorf("%w: %w", ErrValidationFailed, err)
	}

	if err = db.WithTx(ctx, service.db, func(tx db.Tx) (err error) {
		user, err := service.user.Create(ctx, tx, model)
		if err != nil {
			switch {
			case errors.Is(err, stores.ErrUsernameTaken), errors.Is(err, stores.ErrUserEmailTaken):
				return fmt.Errorf("%w: %w", ErrRecordConflict, err)
			default:
				return
			}
		}

		out, err = service.token.IssuePair(ctx, user.ID)
		if err != nil {
			return
		}

		refresh := models.NewRefreshToken(out.RefreshToken, user.ID, out.RefreshTokenExpiresAt)
		if _, err = service.refresh.Create(ctx, tx, refresh); err != nil {
			return
		}

		return
	}); err != nil {
		return
	}

	return
}

func (service authImpl) Login(ctx context.Context, in dtos.UserLogin) (out tokens.Pair, err error) {
	model, err := models.NewUserFromLoginDTO(in)
	if err != nil {
		return out, fmt.Errorf("%w: %w", ErrValidationFailed, err)
	}

	user, err := service.user.GetByUsername(ctx, service.db, model.Username)
	if err != nil {
		dummyPasswordHash.Matches(in.Password)
		return out, fmt.Errorf("%w: %w", ErrInvalidCredentials, err)
	}

	if !user.PasswordMatches(in.Password) {
		return out, fmt.Errorf("%w: %w", ErrInvalidCredentials, ErrUserPasswordNotMatch)
	}

	out, err = service.token.IssuePair(ctx, user.ID)
	if err != nil {
		return
	}

	refresh := models.NewRefreshToken(out.RefreshToken, user.ID, out.RefreshTokenExpiresAt)
	if _, err = service.refresh.Create(ctx, service.db, refresh); err != nil {
		return
	}

	return
}

func (service authImpl) Logout(ctx context.Context, refreshToken string) (err error) {
	refresh, err := service.refresh.GetByHash(ctx, service.db, vos.NewTokenHash(refreshToken))
	if err != nil {
		if errors.Is(err, stores.ErrRefreshTokenNotFound) {
			return fmt.Errorf("%w: %w", ErrInvalidCredentials, ErrRefreshTokenExpired)
		}
		return
	}

	if _, err = service.refresh.RevokeByFamilyID(ctx, service.db, refresh.FamilyID); err != nil {
		return
	}

	return
}

func (service authImpl) Refresh(ctx context.Context, refreshToken string) (out tokens.Pair, err error) {
	refresh, err := service.refresh.GetByHash(ctx, service.db, vos.NewTokenHash(refreshToken))
	if err != nil {
		if errors.Is(err, stores.ErrRefreshTokenNotFound) {
			return out, fmt.Errorf("%w: %w", ErrInvalidCredentials, ErrRefreshTokenExpired)
		}
		return
	}

	if refresh.RevokedAt.Valid {
		if _, err = service.refresh.RevokeByFamilyID(ctx, service.db, refresh.FamilyID); err != nil {
			return
		}
		return out, fmt.Errorf("%w: %w", ErrInvalidCredentials, ErrRefreshTokenReused)
	}

	if time.Now().After(refresh.ExpiresAt) {
		return out, fmt.Errorf("%w: %w", ErrInvalidCredentials, ErrRefreshTokenExpired)
	}

	if err = db.WithTx(ctx, service.db, func(tx db.Tx) (err error) {
		refresh, err := service.refresh.RevokeByHash(ctx, tx, refresh.TokenHash)
		if err != nil {
			return
		}

		out, err = service.token.IssuePair(ctx, refresh.UserID)
		if err != nil {
			return
		}

		refresh = models.NewRefreshTokenInFamily(out.RefreshToken, refresh.UserID, refresh.FamilyID, out.RefreshTokenExpiresAt)
		if _, err = service.refresh.Create(ctx, tx, refresh); err != nil {
			return
		}

		return
	}); err != nil {
		return
	}

	return
}
