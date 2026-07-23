package middlewares

import (
	"context"
	"errors"
	"net/http"
	"strings"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

var (
	ErrUnexpectedSigningMethod = errors.New("unexpected signing method")
)

var userIDKey = "user-id"

func AuthMaybeExpired(secret string, acceptExpired bool) Middleware {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			authorization := r.Header.Get("Authorization")
			switch {
			case authorization == "":
				http.Error(w, "missing Authorization header", http.StatusUnauthorized)
				return
			case !strings.HasPrefix(authorization, "Bearer "):
				http.Error(w, "malformed Authorization header", http.StatusBadRequest)
				return
			}
			authorization = strings.TrimPrefix(authorization, "Bearer ")

			var claims jwt.RegisteredClaims
			token, err := jwt.ParseWithClaims(authorization, &claims, func(t *jwt.Token) (any, error) {
				if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
					return nil, ErrUnexpectedSigningMethod
				}
				return []byte(secret), nil
			}, jwt.WithValidMethods([]string{"HS256"}))
			if err != nil {
				if !(acceptExpired && errors.Is(err, jwt.ErrTokenExpired)) {
					http.Error(w, "unauthorized", http.StatusUnauthorized)
					return
				}
				err = nil
			}

			subject, _ := token.Claims.GetSubject()

			userID, err := uuid.Parse(subject)
			if err != nil {
				http.Error(w, "invalid user ID", http.StatusUnprocessableEntity)
				return
			}

			ctx := context.WithValue(r.Context(), userIDKey, userID)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

func Auth(secret string) Middleware {
	return func(next http.Handler) http.Handler {
		return AuthMaybeExpired(secret, false)(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			next.ServeHTTP(w, r)
		}))
	}
}

func GetUserID(ctx context.Context) uuid.UUID {
	id, ok := ctx.Value(userIDKey).(uuid.UUID)
	if !ok {
		id = uuid.Nil
	}
	return id
}
