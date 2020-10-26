package v0

import "github.com/vmmgr/controller/pkg/api/core/group"

func updateAdminUser(input, replace group.Group) (group.Group, error) {

	// uint boolean
	//Lock
	if input.Lock != replace.Lock {
		replace.Lock = input.Lock
	}
	//Status
	if input.Status != replace.Status {
		replace.Status = input.Status
	}

	return replace, nil
}
