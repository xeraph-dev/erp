package auth

import (
	"erp/internal/auth/handlers"
	"erp/internal/auth/services"
	"erp/internal/auth/stores"
	"erp/internal/shared/api"
	"erp/internal/shared/db"
	"erp/internal/shared/tokens"
)

func Handlers(db db.TxBeginner, jwtSecret string) func(group *api.Group) {
	return func(group *api.Group) {
		user := stores.NewUser()
		refresh := stores.NewRefreshToken()
		role := stores.NewRole()
		permission := stores.NewPermission()

		token := tokens.New(jwtSecret)
		auth := services.NewAuth(db, token, user, refresh, role, permission)

		group.HandleFunc("POST /api/auth/register", handlers.Register(auth))
		group.HandleFunc("POST /api/auth/login", handlers.Login(auth))
		group.HandleFunc("POST /api/auth/logout", handlers.Logout(auth))
		group.HandleFunc("POST /api/auth/refresh", handlers.Refresh(auth))
	}
}
