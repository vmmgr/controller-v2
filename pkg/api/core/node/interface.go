package node

import (
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
	Status bool        `json:"status"`
	Error  string      `json:"error"`
	Node   []core.Node `json:"node"`
}

type ResultOne struct {
	Status bool      `json:"status"`
	Error  string    `json:"error"`
	Node   core.Node `json:"node"`
}

type ResultDatabase struct {
	Err  error
	Node []core.Node
}
