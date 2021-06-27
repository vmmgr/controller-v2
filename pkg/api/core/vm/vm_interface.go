package vm

import (
	cloudinit "github.com/vmmgr/controller/pkg/api/core/vm/cloudinit/v0"
	"github.com/vmmgr/controller/pkg/api/core/vm/nic"
	"github.com/vmmgr/controller/pkg/api/core/vm/storage"
)

type VirtualMachine struct {
	Name           string              `json:"name"`
	UUID           string              `json:"uuid"`
	Memory         uint                `json:"memory"`
	CPUMode        uint                `json:"cpu_mode"` //0:custom 1:host-model 2:pass-through
	VCPU           uint                `json:"vcpu"`
	OS             OS                  `json:"os"`
	VNCPort        uint                `json:"vnc_port"`
	WebSocketPort  uint                `json:"websocket_port"`
	KeyMap         string              `json:"keymap"`
	NIC            []nic.NIC           `json:"nic"`
	Storage        []storage.VMStorage `json:"storage"`
	CloudInit      cloudinit.CloudInit `json:"cloudinit"`
	CloudInitApply bool                `json:"cloudinit_apply"`
	Template       TemplateVM          `json:"template"`
	Stat           uint                `json:"stat"`
}

type TemplateVM struct {
	Apply   bool            `json:"apply"`
	Storage storage.Storage `json:"storage"`
}

type OS struct {
	Boot   []string `json:"boot"`
	Kernel string   `json:"kernel"`
	Arch   uint     `json:"arch"`
	Type   string   `json:"type"`
}

type Address struct {
	PCICount  uint
	DiskCount uint
}

type NIC struct {
	VMID      uint   `json:"vm_id"`
	NodeNICID uint   `json:"node_nic_id"`
	GroupID   uint   `json:"group_id"`
	Name      string `json:"name"`
	Type      uint   `json:"type"`
	Driver    uint   `json:"driver"`
	Mode      uint   `json:"mode"`
	Mac       string `json:"mac"`
	Vlan      uint   `json:"vlan"`
	Comment   string `json:"comment"`
	Lock      *bool  `json:"lock"`
}
