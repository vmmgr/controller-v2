package v0

import (
	"fmt"
	"github.com/vmmgr/controller/pkg/api/core/vm"
)

func check(input vm.VM) error {
	// check
	if input.GroupID == 0 {
		return fmt.Errorf("no data: groupID")
	}
	if input.NodeID == 0 {
		return fmt.Errorf("no data: nodeID")
	}
	if input.Name == "" {
		return fmt.Errorf("no data: name")
	}
	if input.UUID == "" {
		return fmt.Errorf("no data: uuid")
	}
	//if input.CPUModel == "" {
	//	return fmt.Errorf("no data: CPUModel")
	//}
	//if input.CPU == 0 {
	//	return fmt.Errorf("no data: CPU")
	//}
	if input.VNCPort == 0 {
		return fmt.Errorf("no data: vnc port")
	}
	//if input.Boot == 0 {
	//	return fmt.Errorf("no data: boot")
	//}
	//if input.Memory == 0 {
	//	return fmt.Errorf("no data: memory")
	//}
	//if input.OS == 0 {
	//	return fmt.Errorf("no data: os")
	//}

	return nil
}
