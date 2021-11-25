package v2

import (
	"github.com/libvirt/libvirt-go"
	"github.com/vmmgr/controller/pkg/api/core"
	"github.com/vmmgr/controller/pkg/api/core/node"
	dbNode "github.com/vmmgr/controller/pkg/api/store/node/v0"
	"gorm.io/gorm"
	"log"
)

func connectLibvirt(nodeID uint) (*core.Node, *libvirt.Connect, error) {
	resultNode := dbNode.Get(node.ID, &core.Node{Model: gorm.Model{ID: nodeID}})
	if resultNode.Err != nil {
		return nil, nil, resultNode.Err
	}

	conn, err := libvirt.NewConnect("qemu+ssh://" + resultNode.Node[0].User + "@" + resultNode.Node[0].HostName + "/system")
	if err != nil {
		log.Fatalf("failed to connect to qemu")
		return nil, nil, err
	}

	return &resultNode.Node[0], conn, nil
}
