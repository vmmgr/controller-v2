package v0

import (
	"fmt"
	"github.com/vmmgr/controller/pkg/api/core"
	"github.com/vmmgr/controller/pkg/api/core/user"
	"strings"
)

func checkAdmin(input user.CreateAdmin) error {
	// check
	if input.Name == "" {
		return fmt.Errorf("no data: name")
	}
	if input.NameEn == "" {
		return fmt.Errorf("no data: name_en")
	}
	if !strings.Contains(input.Mail, "@") {
		return fmt.Errorf("wrong email address")
	}
	if input.Pass == "" {
		return fmt.Errorf("no data: pass")
	}

	return nil
}

func check(input core.User) error {
	// check
	if input.Name == "" {
		return fmt.Errorf("no data: name")
	}
	return nil
}
