package v0

import (
	"github.com/vmmgr/controller/pkg/api/core"
)

func updateAdminUser(input, replace core.Region) (core.Region, error) {

	//Title
	if input.Name != "" {
		replace.Name = input.Name
	}

	return replace, nil
}
