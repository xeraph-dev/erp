package db

import "slices"

type Permission string
type Permissions []Permission

func (permissions Permissions) Satisfies(other ...Permission) (ok bool) {
	for _, permission := range other {
		if !slices.Contains(permissions, permission) {
			return false
		} else {
			ok = true
		}
	}
	return ok
}

// const (
// 	PermissionUsersCreate Permission = "users:create"
// 	PermissionUsersRead   Permission = "users:read"
// 	PermissionUsersUpdate Permission = "users:update"
// 	PermissionUsersDelete Permission = "users:delete"

// 	PermissionRolesCreate Permission = "roles:create"
// 	PermissionRolesRead   Permission = "roles:read"
// 	PermissionRolesUpdate Permission = "roles:update"
// 	PermissionRolesDelete Permission = "roles:delete"

// 	PermissionRolesUsersAssign Permission = "roles:users:assign"
// )
