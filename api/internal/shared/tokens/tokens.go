package tokens

import (
	"context"
	"crypto/rand"
	"errors"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

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

func (tokens Tokens) ParseAccessToken(accessToken string) (userID uuid.UUID, err error) {
	token, err := jwt.ParseWithClaims(accessToken, jwt.RegisteredClaims{},
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
		return userID, ErrInvalidAccessToken
	}

	subject, _ := token.Claims.GetSubject()

	userID, err = uuid.Parse(subject)
	if err != nil {
		return userID, ErrInvalidAccessToken
	}

	return
}

func (tokens Tokens) IssuePair(ctx context.Context, userID uuid.UUID) (pair Pair, err error) {
	accessToken, accessTokenExpiresAt, err := tokens.issueAccessToken(userID)
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

func (tokens Tokens) issueAccessToken(userID uuid.UUID) (token string, expiresAt time.Time, err error) {
	now := time.Now()
	expiresAt = now.Add(accessTokenTTL)

	claims := jwt.RegisteredClaims{
		Subject:   userID.String(),
		ExpiresAt: jwt.NewNumericDate(expiresAt),
		IssuedAt:  jwt.NewNumericDate(now),
		NotBefore: jwt.NewNumericDate(now),
	}

	token, err = jwt.NewWithClaims(jwt.SigningMethodHS256, claims).SignedString([]byte(tokens.secret))
	if err != nil {
		return token, expiresAt, fmt.Errorf("services/JWT.IssueAccessToken: %w", err)
	}

	return
}

func (Tokens) issueRefreshToken() (token string, expiresAt time.Time) {
	return rand.Text(), time.Now().Add(refreshTokenTTL)
}
