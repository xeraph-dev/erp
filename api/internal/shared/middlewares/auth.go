package middlewares

import (
	"context"
	"erp/internal/shared/db"
	"erp/internal/shared/tokens"
	"net/http"
	"strings"

	"github.com/google/uuid"
)

type authTransport int

const (
	noAuthTransport authTransport = iota
	authTransportHeader
	authTransportCookie
)

type authCtxKey int

const (
	userIDKey authCtxKey = iota
	userPermissionsKey
)

func Auth(jwtSecret string) func(next http.Handler) http.Handler {
	t := tokens.New(jwtSecret)
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ctx := r.Context()

			var accessToken string

			if authorization := r.Header.Get("Authorization"); authorization != "" {
				var ok bool
				accessToken, ok = strings.CutPrefix(authorization, "Bearer ")
				if !ok {
					http.Error(w, "malformed Authorization header", http.StatusUnauthorized)
					return
				}
				if accessToken == "" {
					http.Error(w, "unauthorized", http.StatusUnauthorized)
					return
				}
			} else {
				cookie, err := r.Cookie("access_token")
				if err != nil || cookie.Value == "" {
					http.Error(w, "unauthorized", http.StatusUnauthorized)
					return
				}
				accessToken = cookie.Value
			}

			userID, permissions, err := t.ParseAccessToken(accessToken)
			if err != nil {
				http.Error(w, "unauthorized", http.StatusUnauthorized)
				return
			}

			ctx = context.WithValue(ctx, userIDKey, userID)
			ctx = context.WithValue(ctx, userPermissionsKey, permissions)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

func GetUserID(ctx context.Context) uuid.UUID {
	if userID, ok := ctx.Value(userIDKey).(uuid.UUID); ok {
		return userID
	}
	return uuid.Nil
}

func GetUserPermissions(ctx context.Context) db.Permissions {
	if permissions, ok := ctx.Value(userPermissionsKey).(db.Permissions); ok {
		return permissions
	}
	return nil
}
