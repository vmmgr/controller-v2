package auth

import (
	"github.com/vmmgr/controller/pkg/api/core/group"
	"github.com/vmmgr/controller/pkg/api/core/user"
)

type UserResult struct {
	User user.User
	Err  error
}

type GroupResult struct {
	Group group.Group
	User  user.User
	Err   error
}

type AdminStruct struct {
	User string
	Pass string
}

type AdminResult struct {
	AdminID uint
	Err     error
}
