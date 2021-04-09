package token

import "github.com/vmmgr/controller/pkg/api/core"

const (
	ID                      = 0
	UserToken               = 10
	AccessToken             = 11
	UserTokenAndAccessToken = 12
	ExpiredTime             = 13
	AdminToken              = 20
	AddToken                = 100
	UpdateToken             = 101
	UpdateAll               = 110
)

type Result struct {
	Token []core.Token `json:"token"`
}

type ResultTmpToken struct {
	Token string `json:"token"`
}

type ResultDatabase struct {
	Err   error
	Token []core.Token
}
