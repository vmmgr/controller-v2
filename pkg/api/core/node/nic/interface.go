package nic

import "github.com/jinzhu/gorm"

const (
	ID        = 0
	NodeID    = 1
	GroupID   = 2
	AdminOnly = 3
	Name      = 4
	Enable    = 5
	Virtual   = 6
	UpdateAll = 110
)

type NIC struct {
	gorm.Model
	NodeID    uint   `json:"node_id"`
	GroupID   uint   `json:"group_id"`
	AdminOnly *bool  `json:"admin"`
	Name      string `json:"name"`
	Enable    *bool  `json:"enable"`
	Virtual   *bool  `json:"virtual"`
	Type      uint   `json:"type"`
	Vlan      uint   `json:"vlan"`
	Speed     uint   `json:"speed"`
	MAC       string `json:"mac"`
	Comment   string `json:"comment"`
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
