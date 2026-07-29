package tokens

import (
	"context"
	"crypto/rand"
	"erp/internal/shared/db"
	"errors"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

type userClaims struct {
	Permissions db.Permissions `json:"permissions"`
	jwt.RegisteredClaims
}

var (
	ErrUnexpectedSigningMethod = errors.New("unexpected signing method")
	ErrInvalidAccessToken      = errors.New("invalid access token")
)

const (
	accessTokenTTL  = 15 * time.Minute
	refreshTokenTTL = 7 * 24 * time.Hour
)

type Tokens struct {
	secret string
}

func New(secret string) Tokens {
	return Tokens{secret: secret}
}

func (tokens Tokens) ParseAccessToken(accessToken string) (userID uuid.UUID, permissions db.Permissions, err error) {
	var claims userClaims
	token, err := jwt.ParseWithClaims(accessToken, &claims,
		func(t *jwt.Token) (any, error) {
			if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("%w: %v", ErrUnexpectedSigningMethod, t.Header["alg"])
			}
			return []byte(tokens.secret), nil
		},
		jwt.WithValidMethods([]string{"HS256"}),
		jwt.WithExpirationRequired(),
		jwt.WithNotBeforeRequired(),
		jwt.WithLeeway(30*time.Second),
	)
	if err != nil {
		return
	}

	if !token.Valid {
		return userID, permissions, ErrInvalidAccessToken
	}

	userID, err = uuid.Parse(claims.Subject)
	if err != nil {
		return userID, permissions, ErrInvalidAccessToken
	}

	permissions = claims.Permissions
	return
}

func (tokens Tokens) IssuePair(ctx context.Context, userID uuid.UUID, permissions db.Permissions) (pair Pair, err error) {
	accessToken, accessTokenExpiresAt, err := tokens.issueAccessToken(userID, permissions)
	if err != nil {
		return
	}

	refreshToken, refreshTokenExpiresAt := tokens.issueRefreshToken()

	return Pair{
		AccessToken:           accessToken,
		AccessTokenExpiresAt:  accessTokenExpiresAt,
		RefreshToken:          refreshToken,
		RefreshTokenExpiresAt: refreshTokenExpiresAt,
	}, nil
}

func (tokens Tokens) issueAccessToken(userID uuid.UUID, permissions db.Permissions) (token string, expiresAt time.Time, err error) {
	now := time.Now()
	expiresAt = now.Add(accessTokenTTL)

	claims := userClaims{
		Permissions: permissions,
		RegisteredClaims: jwt.RegisteredClaims{
			Subject:   userID.String(),
			ExpiresAt: jwt.NewNumericDate(expiresAt),
			IssuedAt:  jwt.NewNumericDate(now),
			NotBefore: jwt.NewNumericDate(now),
		},
	}

	token, err = jwt.NewWithClaims(jwt.SigningMethodHS256, claims).SignedString([]byte(tokens.secret))
	if err != nil {
		return
	}

	return
}

func (Tokens) issueRefreshToken() (token string, expiresAt time.Time) {
	return rand.Text(), time.Now().Add(refreshTokenTTL)
}
