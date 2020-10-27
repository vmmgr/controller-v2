package vm

import (
	"github.com/jinzhu/gorm"
)

const (
	ID             = 0
	NodeID         = 1
	GroupID        = 2
	UUID           = 3
	NodeAndVNCPort = 4
	Lock           = 5
	UpdateName     = 100
	UpdateNode     = 101
	UpdateGroup    = 102
	UpdateBoot     = 103
	UpdateInfo     = 104
	UpdateAll      = 110
)

type VM struct {
	gorm.Model
	NodeID   uint   `json:"node_id"`
	GroupID  uint   `json:"group_id"`
	Name     string `json:"name"`
	UUID     string `json:"uuid"`
	OS       uint   `json:"os"` //32bit=> 32 64bit=> 64
	CPU      uint   `json:"cpu"`
	CPUModel string `json:"cpu_mode"`
	Memory   uint   `json:"memory"`
	VNCPort  uint   `json:"vnc_port"`
	Boot     uint   `json:"boot"` //0: hd 1:cdrom 2:floppy
	Lock     *bool  `json:"lock"`
}

type Result struct {
	Status bool   `json:"status"`
	Error  string `json:"error"`
	VM     []VM   `json:"vm"`
}

type ResultOne struct {
	Status bool   `json:"status"`
	Error  string `json:"error"`
	VM     VM     `json:"vm"`
}

type ResultDatabase struct {
	Err error
	VMs []VM
}
