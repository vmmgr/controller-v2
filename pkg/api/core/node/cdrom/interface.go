package nic

import (
	"github.com/vmmgr/controller/pkg/api/core"
)

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

type Post struct {
	Name string `json:"name"`
	URL  string `json:"url"`
}

type Delete struct {
	Name string `json:"name"`
}

type Result struct {
	Status bool       `json:"status"`
	Error  string     `json:"error"`
	NIC    []core.NIC `json:"nic"`
}

type ResultOne struct {
	Status bool     `json:"status"`
	Error  string   `json:"error"`
	NIC    core.NIC `json:"nic"`
}

type ResultDatabase struct {
	Err error
	NIC []core.NIC
}
