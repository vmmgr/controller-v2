package region

import (
	"github.com/vmmgr/controller/pkg/api/core"
)

const (
	ID        = 0
	Name      = 1
	UpdateAll = 110
)

type Result struct {
	Status bool          `json:"status"`
	Error  string        `json:"error"`
	Region []core.Region `json:"region"`
}

type ResultOne struct {
	Status bool        `json:"status"`
	Error  string      `json:"error"`
	Region core.Region `json:"region"`
}

type ResultDatabase struct {
	Err    error
	Region []core.Region
}
