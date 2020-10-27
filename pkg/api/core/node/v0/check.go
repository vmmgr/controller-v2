package v0

import (
	"fmt"
	"github.com/vmmgr/controller/pkg/api/core/node"
)

func check(input node.Node) error {
	// check
	if input.ZoneID == 0 {
		return fmt.Errorf("no data: zoneID")
	}
	if input.Name == "" {
		return fmt.Errorf("no data: name")
	}
	if input.GroupID == 0 {
		return fmt.Errorf("no data: groupID")
	}
	if input.IP == "" {
		return fmt.Errorf("no data: ip")
	}
	if input.Port == 0 {
		return fmt.Errorf("no data: port")
	}
	if input.WsPort == 0 {
		return fmt.Errorf("no data: webSocket port")
	}
	if input.ManageNet == "" {
		return fmt.Errorf("no data: management network")
	}

	return nil
}
