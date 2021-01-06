package v0

import "github.com/vmmgr/controller/pkg/api/core/vm"

func updateAdminUser(input, replace vm.VM) (vm.VM, error) {

	//Name
	if input.Name != "" {
		replace.Name = input.Name
	}
	if input.UUID != "" {
		replace.UUID = input.UUID
	}
	//if input.CPUModel != "" {
	//	replace.CPUModel = input.CPUModel
	//}
	if replace.GroupID != input.GroupID {
		replace.GroupID = input.GroupID
	}
	if replace.NodeID != input.NodeID {
		replace.NodeID = input.NodeID
	}
	//if replace.Boot != input.Boot {
	//	replace.Boot = input.Boot
	//}
	if replace.VNCPort != input.VNCPort {
		replace.VNCPort = input.VNCPort
	}
	//if replace.Memory != input.Memory {
	//	replace.Memory = input.Memory
	//}
	//if replace.CPU != input.CPU {
	//	replace.CPU = input.CPU
	//}
	//if replace.OS != input.OS {
	//	replace.OS = input.OS
	//}

	return replace, nil
}
