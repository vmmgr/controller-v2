package nic

import (
	"github.com/jinzhu/gorm"
)

const (
	ID          = 0
	NodeNICID   = 1
	GroupID     = 2
	Name        = 3
	Type        = 4
	Vlan        = 5
	UpdateName  = 100
	UpdateNodeN = 101
	UpdateGroup = 102
	UpdateMac   = 103
	UpdateAll   = 110
)

type NIC struct {
	gorm.Model
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

type Result struct {
	Status bool   `json:"status"`
	Error  string `json:"error"`
	NIC    []NIC  `json:"nic"`
}

type ResultOne struct {
	Status bool   `json:"status"`
	Error  string `json:"error"`
	NIC    NIC    `json:"nic"`
}

type ResultDatabase struct {
	Err error
	NIC []NIC
}
