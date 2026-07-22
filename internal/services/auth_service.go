package services

import (
	"context"
	"crypto/rand"
	"erp/internal/dtos"
	"erp/internal/middlewares"
	"erp/internal/models"
	"erp/internal/repositories"
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
	ErrUsernameExists    = errors.New("username already exists")
	ErrUserEmailExists   = errors.New("user email already exists")
	ErrUsernameNotExists = errors.New("username does not exists")
	ErrPasswordNotMatch  = errors.New("user password does not match")
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
func (err ErrCreatingUserModel) Error() string { return fmt.Sprintf("creating model: %s", err.err.Error()) }
func (err ErrUserExists) Error() string        { return fmt.Sprintf("user exists: %s", err.err.Error()) }
func (err ErrUserNotExists) Error() string     { return fmt.Sprintf("user not exists: %s", err.err.Error()) }

type AuthService interface {
	Service
	Register(ctx context.Context, in dtos.UserRegister) (out dtos.TokenPair, err error)
	Login(ctx context.Context, in dtos.UserLogin) (out dtos.TokenPair, err error)
	// Logout(ctx context.Context, id uuid.UUID) (err error)
	// Refresh(ctx context.Context, refreshToken string) (out dtos.TokenPair, err error)
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
	logger := middlewares.GetLogger(ctx)

	model, err := models.NewUserFromRegisterDTO(ctx, in)
	if err != nil {
		logger.ErrorContext(ctx, "creating model", "error", err)
		err = NewErrCreatingUserModel(err)
		return
	}

	exists, err := service.user.UsernameExists(ctx, service.db, model.Username)
	if err != nil {
		logger.ErrorContext(ctx, "checking username exists", "error", err)
		return
	} else if exists {
		err = NewErrUserExists(ErrUsernameExists)
		return
	}

	exists, err = service.user.EmailExists(ctx, service.db, model.Email)
	if err != nil {
		logger.ErrorContext(ctx, "checking user email exists", "error", err)
		return
	} else if exists {
		err = NewErrUserExists(ErrUserEmailExists)
		return
	}

	err = withTx(ctx, service.db, func(tx pgx.Tx) (err error) {
		user, err := service.user.Create(ctx, tx, model)
		if err != nil {
			logger.ErrorContext(ctx, "creating entry", "error", err)
			return
		}

		role, err := service.role.GetUser(ctx, tx)
		if err != nil {
			logger.ErrorContext(ctx, "getting role user", "error", err)
			return
		}

		if err = service.role.Assign(ctx, tx, role.ID, user.ID); err != nil {
			logger.ErrorContext(ctx, "assigning role to user", "role_id", role.ID, "user_id", user.ID)
			return
		}

		out, err = service.issueTokenPair(ctx, tx, user.ID)
		if err != nil {
			logger.ErrorContext(ctx, "issuing token pair", "error", err)
			return
		}

		return
	})

	return
}

func (service authServiceImpl) Login(ctx context.Context, in dtos.UserLogin) (out dtos.TokenPair, err error) {
	logger := middlewares.GetLogger(ctx)

	model, err := models.NewUserFromLoginDTO(ctx, in)
	if err != nil {
		logger.ErrorContext(ctx, "creating model", "error", err)
		err = NewErrCreatingUserModel(err)
		return
	}

	exists, err := service.user.UsernameExists(ctx, service.db, model.Username)
	if err != nil {
		logger.ErrorContext(ctx, "checking username exists", "error", err)
		return
	} else if !exists {
		err = NewErrUserNotExists(ErrUsernameNotExists)
		return
	}

	user, err := service.user.GetByUsername(ctx, service.db, model.Username)
	if err != nil {
		logger.ErrorContext(ctx, "getting user by username", "error", err)
		return
	}

	if !user.PasswordHash.Matches(in.Password) {
		err = NewErrUserNotExists(ErrPasswordNotMatch)
		return
	}

	out, err = service.issueTokenPair(ctx, service.db, user.ID)
	if err != nil {
		logger.ErrorContext(ctx, "issuing token pair", "error", err)
		return
	}

	return
}

func (service authServiceImpl) issueTokenPair(ctx context.Context, db repositories.Querier, userID uuid.UUID) (out dtos.TokenPair, err error) {
	logger := middlewares.GetLogger(ctx)

	accessToken, accessTokenExpiresAt, err := service.issueAccessToken(ctx, userID)
	if err != nil {
		logger.ErrorContext(ctx, "issuing access token", "error", err)
		return
	}

	refreshToken, refreshTokenExpiresAt := service.issueRefreshToken()

	model := models.NewRefreshToken(userID, refreshToken, refreshTokenExpiresAt)

	if err = service.refresh.Create(ctx, db, model); err != nil {
		logger.ErrorContext(ctx, "creating refresh token entry", "error", err)
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

func (service authServiceImpl) issueAccessToken(ctx context.Context, userID uuid.UUID) (token string, expiresAt time.Time, err error) {
	logger := middlewares.GetLogger(ctx)

	expiresAt = time.Now().Add(accessTokenTTL)

	claims := jwt.RegisteredClaims{
		ID:        uuid.NewString(),
		Subject:   userID.String(),
		IssuedAt:  jwt.NewNumericDate(time.Now()),
		ExpiresAt: jwt.NewNumericDate(expiresAt),
	}

	token, err = jwt.NewWithClaims(jwt.SigningMethodHS256, claims).SignedString([]byte(service.jwtSecret))
	if err != nil {
		logger.ErrorContext(ctx, "signing access token", "error", err)
		return
	}

	return
}

func (service authServiceImpl) issueRefreshToken() (token string, expiresAt time.Time) {
	return rand.Text(), time.Now().Add(refreshTokenTTL)
}
