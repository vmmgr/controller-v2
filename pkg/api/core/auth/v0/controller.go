package v0

import (
	"fmt"
	"github.com/vmmgr/controller/pkg/api/core/controller"
	"github.com/vmmgr/controller/pkg/api/core/tool/config"
	"github.com/vmmgr/controller/pkg/api/core/tool/hash"
)

func ControllerAuthentication(con controller.Controller) error {
	if con.Token1 != config.Conf.Controller.Auth.Token1 {
		return fmt.Errorf("auth error! ")
	}
	if con.Token2 != hash.Generate(config.Conf.Controller.Auth.Token2+config.Conf.Controller.Auth.Token3) {
		return fmt.Errorf("auth error! ")
	}
	return nil
}
