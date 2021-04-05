package v0

import (
	"fmt"
	"github.com/vmmgr/controller/pkg/api/core"
	"github.com/vmmgr/controller/pkg/api/core/auth"
	"github.com/vmmgr/controller/pkg/api/core/token"
	"github.com/vmmgr/controller/pkg/api/core/tool/config"
	dbToken "github.com/vmmgr/controller/pkg/api/store/token/v0"
)

func AdminRadiusAuthentication(data auth.AdminStruct) auth.AdminResult {

	if config.Conf.Controller.Admin.AdminAuth.User == data.User && config.Conf.Controller.Admin.AdminAuth.Pass == data.Pass {
		return auth.AdminResult{AdminID: 0, Err: nil}
	}
	// Todo Radius認証追加予定
	return auth.AdminResult{Err: fmt.Errorf("failed")}
}

func AdminAuthentication(accessToken string) auth.AdminResult {
	tokenResult := dbToken.Get(token.AdminToken, &core.Token{AccessToken: accessToken})
	if tokenResult.Err != nil {
		return auth.AdminResult{Err: tokenResult.Err}
	}
	return auth.AdminResult{Err: nil}
}
