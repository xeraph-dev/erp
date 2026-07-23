package queries

import _ "embed"

//go:embed CreateUser.sql
var CreateUser string

//go:embed SetCurrentUserID.sql
var SetCurrentUserID string

//go:embed GetRoleUser.sql
var GetRoleUser string

//go:embed AssignRoleToUser.sql
var AssignRoleToUser string

//go:embed UsernameExists.sql
var UsernameExists string

//go:embed UserEmailExists.sql
var UserEmailExists string

//go:embed GetUserByUsername.sql
var GetUserByUsername string

//go:embed CreateRefreshToken.sql
var CreateRefreshToken string

//go:embed GetRefreshTokenByTokenHash.sql
var GetRefreshTokenByTokenHash string

//go:embed RevokeRefreshToken.sql
var RevokeRefreshToken string

//go:embed RevokeRefreshTokenFamily.sql
var RevokeRefreshTokenFamily string
