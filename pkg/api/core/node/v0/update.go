package v0

import (
	"github.com/vmmgr/controller/pkg/api/core"
)

func updateAdminUser(input, replace core.Node) (core.Node, error) {

	//Name
	if input.Name != "" {
		replace.Name = input.Name
	}
	if input.Comment != "" {
		replace.Comment = input.Comment
	}
	if input.ManageNet != "" {
		replace.ManageNet = input.ManageNet
	}
	if input.IP != "" {
		replace.IP = input.IP
	}
	if replace.ZoneID != input.ZoneID {
		replace.ZoneID = input.ZoneID
	}
	if replace.GroupID != input.GroupID {
		replace.GroupID = input.GroupID
	}
	if replace.Port != input.Port {
		replace.Port = input.Port
	}
	if replace.WsPort != input.WsPort {
		replace.WsPort = input.WsPort
	}

	return replace, nil
}
