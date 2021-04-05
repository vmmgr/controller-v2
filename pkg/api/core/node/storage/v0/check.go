package v0

import (
	"fmt"
	"github.com/vmmgr/controller/pkg/api/core"
)

func check(input core.Storage) error {
	// check
	if input.Name == "" {
		return fmt.Errorf("no data: name")
	}
	if input.NodeID == 0 {
		return fmt.Errorf("no data: nodeID")
	}
	if input.Path == "" {
		return fmt.Errorf("no data: path")
	}
	if input.MaxCapacity == 0 {
		return fmt.Errorf("no data: Max capacity")
	}
	if input.Type == 0 {
		return fmt.Errorf("no data: type")
	}

	return nil
}
