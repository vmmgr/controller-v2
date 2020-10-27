package v0

import (
	"github.com/vmmgr/controller/pkg/api/core/region"
)

func updateAdminUser(input, replace region.Region) (region.Region, error) {

	//Title
	if input.Name != "" {
		replace.Name = input.Name
	}

	return replace, nil
}
