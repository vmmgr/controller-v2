package v0

import (
	"github.com/vmmgr/controller/pkg/api/core"
)

func updateAdminUser(input, replace core.Zone) (core.Zone, error) {

	//Name
	if input.Name != "" {
		replace.Name = input.Name
	}
	if input.Postcode != "" {
		replace.Postcode = input.Postcode
	}
	if input.Address != "" {
		replace.Address = input.Address
	}
	if input.Tel != "" {
		replace.Tel = input.Tel
	}
	if input.Mail != "" {
		replace.Mail = input.Mail
	}
	if replace.RegionID != input.RegionID {
		replace.RegionID = input.RegionID
	}

	return replace, nil
}
