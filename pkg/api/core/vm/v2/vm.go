package v2

import (
	"github.com/libvirt/libvirt-go"
	"github.com/vmmgr/controller/pkg/api/core"
	"github.com/vmmgr/controller/pkg/api/core/vm"
)

type VMHandler struct {
	Conn    *libvirt.Connect
	VM      vm.VirtualMachine
	Node    core.Node
	GroupID uint
	IPID    uint
}

func NewVMHandler(input VMHandler) *VMHandler {
	return &VMHandler{
		Conn:    input.Conn,
		VM:      input.VM,
		Node:    input.Node,
		GroupID: input.GroupID,
		IPID:    input.IPID,
	}
}
