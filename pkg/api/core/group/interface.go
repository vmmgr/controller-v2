package group

import (
	"github.com/jinzhu/gorm"
)

const (
	ID           = 0
	OrgJa        = 1
	Org          = 2
	UpdateID     = 100
	UpdateOrg    = 101
	UpdateStatus = 102
	UpdateTechID = 103
	UpdateInfo   = 104
	UpdateAll    = 110
)

type Group struct {
	gorm.Model
	Org       string `json:"org"`
	Status    uint   `json:"status"`
	Comment   string `json:"comment"`
	Vlan      uint   `json:"vlan"`
	Lock      bool   `json:"lock"`
	MaxVM     uint   `json:"max_VM"`
	MaxCPU    uint   `json:"max_cpu"`
	MaxMemory uint   `json:"max_memory"`
}

type Result struct {
	Status bool    `json:"status"`
	Error  string  `json:"error"`
	Group  []Group `json:"group"`
}

type ResultOne struct {
	Status bool   `json:"status"`
	Error  string `json:"error"`
	Group  Group  `json:"group"`
}

type ResultAll struct {
	Status bool   `json:"status"`
	Error  string `json:"error"`
	Group  Group  `json:"group"`
}

type ResultDatabase struct {
	Err   error
	Group []Group
}
