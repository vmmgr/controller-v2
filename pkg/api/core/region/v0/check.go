package v0

import (
	"fmt"
	"github.com/vmmgr/controller/pkg/api/core/region"
)

func check(input region.Region) error {
	// check
	if input.Name == "" {
		return fmt.Errorf("no data: name")
	}
	return nil
}
