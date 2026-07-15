package queries

import _ "embed"

//go:embed CreateUser.sql
var CreateUser string

//go:embed GetUserByName.sql
var GetUserByName string
