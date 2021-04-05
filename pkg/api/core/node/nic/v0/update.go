package v0

import (
	"github.com/vmmgr/controller/pkg/api/core"
)

func updateAdminUser(input, replace core.NIC) (core.NIC, error) {

	//Name
	if input.Name != "" {
		replace.Name = input.Name
	}
	if input.Comment != "" {
		replace.Comment = input.Comment
	}
	if input.MAC != "" {
		replace.MAC = input.MAC
	}
	if replace.AdminOnly != input.AdminOnly {
		replace.AdminOnly = input.AdminOnly
	}
	if replace.NodeID != input.NodeID {
		replace.NodeID = input.NodeID
	}
	if replace.GroupID != input.GroupID {
		replace.GroupID = input.GroupID
	}
	if replace.Speed != input.Speed {
		replace.Speed = input.Speed
	}
	if replace.Type != input.Type {
		replace.Type = input.Type
	}
	if replace.Vlan != input.Vlan {
		replace.Vlan = input.Vlan
	}

	return replace, nil
}
