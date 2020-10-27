package v0

import (
	"github.com/vmmgr/controller/pkg/api/core/region/zone"
)

func updateAdminUser(input, replace zone.Zone) (zone.Zone, error) {

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
