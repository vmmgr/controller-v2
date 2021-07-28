package v0

import (
	"github.com/vmmgr/controller/pkg/api/core"
)

func updateAdminUser(input, replace core.Group) (core.Group, error) {

	// uint boolean
	//Status
	if input.Status != replace.Status {
		replace.Status = input.Status
	}

	return replace, nil
}
