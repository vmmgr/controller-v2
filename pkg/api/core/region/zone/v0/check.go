package v0

import (
	"fmt"
	"github.com/vmmgr/controller/pkg/api/core"
)

func check(input core.Zone) error {
	// check
	if input.RegionID == 0 {
		return fmt.Errorf("no data: regionID")
	}
	if input.Name == "" {
		return fmt.Errorf("no data: name")
	}
	if input.Postcode == "" {
		return fmt.Errorf("no data: postcode")
	}
	if input.Address == "" {
		return fmt.Errorf("no data: address")
	}
	if input.Mail == "" {
		return fmt.Errorf("no data: mail")
	}
	if input.Tel == "" {
		return fmt.Errorf("no data: tel")
	}

	return nil
}
