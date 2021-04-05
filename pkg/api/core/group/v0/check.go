package v0

import (
	"fmt"
	"github.com/vmmgr/controller/pkg/api/core"
)

func check(input core.Group) error {
	// check
	if input.Org == "" {
		return fmt.Errorf("no data: org")
	}
	return nil
}
