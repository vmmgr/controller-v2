package v0

import (
	"fmt"
	"github.com/vmmgr/controller/pkg/api/core/group"
)

func check(input group.Group) error {
	// check
	if input.Org == "" {
		return fmt.Errorf("no data: org")
	}
	return nil
}
