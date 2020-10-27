package v0

import (
	"fmt"
	"github.com/vmmgr/controller/pkg/api/core/node/nic"
)

func check(input nic.NIC) error {
	// check
	if input.Name == "" {
		return fmt.Errorf("no data: name")
	}
	if input.NodeID == 0 {
		return fmt.Errorf("no data: nodeID")
	}
	if input.MAC == "" {
		return fmt.Errorf("no data: mac")
	}
	if input.Vlan == 0 {
		return fmt.Errorf("no data: vlan")
	}
	if input.Type == 0 {
		return fmt.Errorf("no data: type")
	}
	if input.Speed == 0 {
		return fmt.Errorf("no data: speed")
	}
	if input.GroupID == 0 {
		return fmt.Errorf("no data: groupID")
	}

	return nil
}
