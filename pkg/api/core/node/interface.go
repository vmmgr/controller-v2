package node

import (
	libvirtxml "github.com/libvirt/libvirt-go-xml"
	"github.com/vmmgr/controller/pkg/api/core"
)

const (
	ID        = 0
	ZoneID    = 1
	GroupID   = 2
	AdminOnly = 3
	Name      = 4
	UpdateAll = 110
)

type Result struct {
	Node []core.Node `json:"node"`
}

type ResultDevice struct {
	PCI []libvirtxml.NodeDevice `json:"pci"`
	USB []libvirtxml.NodeDevice `json:"usb"`
}

type ResultOne struct {
	Node core.Node `json:"node"`
}

type ResultDatabase struct {
	Err  error
	Node []core.Node
}
