package storage

import (
	"github.com/vmmgr/controller/pkg/api/core"
)

const (
	ID        = 0
	NodeID    = 1
	AdminOnly = 2
	Name      = 3
	UpdateAll = 110
)

type Result struct {
	Storage []core.Storage `json:"storage"`
}

type ResultOne struct {
	Status  bool         `json:"status"`
	Error   string       `json:"error"`
	Storage core.Storage `json:"storage"`
}

type ResultDatabase struct {
	Err     error
	Storage []core.Storage
}
