package group

import (
	"github.com/jinzhu/gorm"
)

const (
	ID           = 0
	OrgJa        = 1
	Org          = 2
	Email        = 3
	UpdateID     = 100
	UpdateOrg    = 101
	UpdateStatus = 102
	UpdateTechID = 103
	UpdateInfo   = 104
	UpdateAll    = 110
)

type VM struct {
	gorm.Model
	Status    uint   `json:"status"`
	Comment   string `json:"comment"`
	Vlan      uint   `json:"vlan"`
	Lock      bool   `json:"lock"`
	MaxVM     uint   `json:"max_VM"`
	MaxCPU    uint   `json:"max_cpu"`
	MaxMemory uint   `json:"max_memory"`
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
	VM  []VM
}
