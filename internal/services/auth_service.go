package services

import (
	"context"
	"crypto/rand"
	"erp/internal/dtos"
	"erp/internal/models"
	"erp/internal/repositories"
	"erp/internal/vos"
	"errors"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

const (
	accessTokenTTL  = 15 * time.Minute
	refreshTokenTTL = 7 * 24 * time.Hour
)

var (
	ErrUsernameExists       = errors.New("username already exists")
	ErrUserEmailExists      = errors.New("user email already exists")
	ErrUsernameNotExists    = errors.New("username does not exists")
	ErrPasswordNotMatch     = errors.New("user password does not match")
	ErrRefreshTokenNotFound = errors.New("refresh token not found")
)

type ErrCreatingUserModel struct{ err error }
type ErrUserExists struct{ err error }
type ErrUserNotExists struct{ err error }

func NewErrCreatingUserModel(err error) ErrCreatingUserModel { return ErrCreatingUserModel{err} }
func NewErrUserExists(err error) ErrUserExists               { return ErrUserExists{err} }
func NewErrUserNotExists(err error) ErrUserNotExists         { return ErrUserNotExists{err} }

func (err ErrCreatingUserModel) Unwrap() error { return err.err }
func (err ErrUserExists) Unwrap() error        { return err.err }
func (err ErrUserNotExists) Unwrap() error     { return err.err }
func (err ErrCreatingUserModel) Error() string {
	return fmt.Sprintf("creating model: %s", err.err.Error())
}
func (err ErrUserExists) Error() string { return fmt.Sprintf("user exists: %s", err.err.Error()) }
func (err ErrUserNotExists) Error() string {
	return fmt.Sprintf("user not exists: %s", err.err.Error())
}

type AuthService interface {
	Service
	Register(ctx context.Context, in dtos.UserRegister) (out dtos.TokenPair, err error)
	Login(ctx context.Context, in dtos.UserLogin) (out dtos.TokenPair, err error)
	Logout(ctx context.Context, refreshToken string) (err error)
	Refresh(ctx context.Context, refreshToken string) (out dtos.TokenPair, err error)
}

type authServiceImpl struct {
	jwtSecret string
	db        *pgxpool.Pool
	user      repositories.UserRepository
	role      repositories.RoleRepository
	refresh   repositories.RefreshTokenRepository
}

var _ AuthService = (*authServiceImpl)(nil)

func NewAuthService(jwtSecret string, db *pgxpool.Pool, users repositories.UserRepository, roles repositories.RoleRepository, refresh repositories.RefreshTokenRepository) AuthService {
	return authServiceImpl{
		jwtSecret: jwtSecret,
		db:        db,
		user:      users,
		role:      roles,
		refresh:   refresh,
	}
}

func (authServiceImpl) __internal() {}

func (service authServiceImpl) Register(ctx context.Context, in dtos.UserRegister) (out dtos.TokenPair, err error) {
	model, err := models.NewUserFromRegisterDTO(in)
	if err != nil {
		err = NewErrCreatingUserModel(err)
		return
	}

	exists, err := service.user.UsernameExists(ctx, service.db, model.Username)
	if err != nil {
		err = fmt.Errorf("checking username exists: %w", err)
		return
	} else if exists {
		err = NewErrUserExists(ErrUsernameExists)
		return
	}

	exists, err = service.user.EmailExists(ctx, service.db, model.Email)
	if err != nil {
		err = fmt.Errorf("checking user email exists: %w", err)
		return
	} else if exists {
		err = NewErrUserExists(ErrUserEmailExists)
		return
	}

	err = withTx(ctx, service.db, func(tx pgx.Tx) (err error) {
		user, err := service.user.Create(ctx, tx, model)
		if err != nil {
			err = fmt.Errorf("creating entry: %w", err)
			return
		}

		role, err := service.role.GetUser(ctx, tx)
		if err != nil {
			err = fmt.Errorf("getting role user: %w", err)
			return
		}

		if err = service.role.Assign(ctx, tx, role.ID, user.ID); err != nil {
			err = fmt.Errorf("assigning role %s to user %s: %w", role.ID, user.ID, err)
			return
		}

		out, err = service.issueTokenPair(ctx, tx, user.ID)
		if err != nil {
			err = fmt.Errorf("issuing token pair: %w", err)
			return
		}

		return
	})

	return
}

func (service authServiceImpl) Login(ctx context.Context, in dtos.UserLogin) (out dtos.TokenPair, err error) {
	model, err := models.NewUserFromLoginDTO(in)
	if err != nil {
		err = NewErrCreatingUserModel(err)
		return
	}

	exists, err := service.user.UsernameExists(ctx, service.db, model.Username)
	if err != nil {
		err = fmt.Errorf("checking username exists: %w", err)
		return
	} else if !exists {
		err = NewErrUserNotExists(ErrUsernameNotExists)
		return
	}

	user, err := service.user.GetByUsername(ctx, service.db, model.Username)
	if err != nil {
		err = fmt.Errorf("getting user by username: %w", err)
		return
	}

	if !user.PasswordMatches(in.Password) {
		err = NewErrUserNotExists(ErrPasswordNotMatch)
		return
	}

	out, err = service.issueTokenPair(ctx, service.db, user.ID)
	if err != nil {
		err = fmt.Errorf("issuing token pair: %w", err)
		return
	}

	return
}

func (service authServiceImpl) Logout(ctx context.Context, refreshToken string) (err error) {
	refresh, err := service.refresh.GetByTokenHash(ctx, service.db, vos.NewTokenHash(refreshToken))
	if err != nil {
		if errors.Is(err, repositories.ErrRefreshTokenNotFound) {
			err = ErrRefreshTokenNotFound
			return
		}
		err = fmt.Errorf("getting refresh token by token hash: %w", err)
		return
	}

	_, err = service.refresh.RevokeFamily(ctx, service.db, refresh)
	if err != nil {
		err = fmt.Errorf("revoking refresh token family: %w", err)
		return
	}

	return
}

func (service authServiceImpl) Refresh(ctx context.Context, refreshToken string) (out dtos.TokenPair, err error) {
	refresh, err := service.refresh.GetByTokenHash(ctx, service.db, vos.NewTokenHash(refreshToken))
	if err != nil {
		if errors.Is(err, repositories.ErrRefreshTokenNotFound) {
			err = ErrRefreshTokenNotFound
			return
		}
		err = fmt.Errorf("getting refresh token by token hash: %w", err)
		return
	}

	if err = withTx(ctx, service.db, func(tx pgx.Tx) (err error) {
		refresh, err := service.refresh.Revoke(ctx, tx, refresh)
		if err != nil {
			err = fmt.Errorf("revoking refresh token: %w", err)
			return
		}

		out, err = service.issueTokenPairInFamily(ctx, service.db, refresh.UserID, refresh.FamilyID)
		if err != nil {
			err = fmt.Errorf("issuing token pair: %w", err)
			return
		}
		return
	}); err != nil {
		return
	}

	return
}

func (service authServiceImpl) issueTokenPairInFamily(ctx context.Context, db repositories.Querier, userID uuid.UUID, familyID uuid.UUID) (out dtos.TokenPair, err error) {
	accessToken, accessTokenExpiresAt, err := service.issueAccessToken(ctx, userID)
	if err != nil {
		err = fmt.Errorf("issuing access token: %w", err)
		return
	}

	refreshToken, refreshTokenExpiresAt := service.issueRefreshToken()

	model := models.NewRefreshToken(userID, familyID, refreshToken, refreshTokenExpiresAt)

	if _, err = service.refresh.Create(ctx, db, model); err != nil {
		err = fmt.Errorf("creating refresh token entry: %w", err)
		return
	}

	out = dtos.TokenPair{
		AccessToken:           accessToken,
		RefreshToken:          refreshToken,
		AccessTokenExpiresAt:  accessTokenExpiresAt,
		RefreshTokenExpiresAt: refreshTokenExpiresAt,
	}
	return
}

func (service authServiceImpl) issueTokenPair(ctx context.Context, db repositories.Querier, userID uuid.UUID) (out dtos.TokenPair, err error) {
	return service.issueTokenPairInFamily(ctx, db, userID, uuid.New())
}

func (service authServiceImpl) issueAccessToken(ctx context.Context, userID uuid.UUID) (token string, expiresAt time.Time, err error) {
	now := time.Now()
	expiresAt = now.Add(accessTokenTTL)

	claims := jwt.RegisteredClaims{
		ID:        uuid.NewString(),
		Subject:   userID.String(),
		ExpiresAt: jwt.NewNumericDate(expiresAt),
		IssuedAt:  jwt.NewNumericDate(now),
		NotBefore: jwt.NewNumericDate(now),
	}

	token, err = jwt.NewWithClaims(jwt.SigningMethodHS256, claims).SignedString([]byte(service.jwtSecret))
	if err != nil {
		err = fmt.Errorf("signing access token: %w", err)
		return
	}

	return
}

func (service authServiceImpl) issueRefreshToken() (token string, expiresAt time.Time) {
	return rand.Text(), time.Now().Add(refreshTokenTTL)
}
