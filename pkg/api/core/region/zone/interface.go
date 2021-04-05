package zone

import (
	"github.com/vmmgr/controller/pkg/api/core"
)

const (
	ID        = 0
	RegionID  = 1
	Name      = 2
	UpdateAll = 110
)

type Result struct {
	Status bool        `json:"status"`
	Error  string      `json:"error"`
	Zone   []core.Zone `json:"zone"`
}

type ResultOne struct {
	Status bool      `json:"status"`
	Error  string    `json:"error"`
	Zone   core.Zone `json:"zone"`
}

type ResultDatabase struct {
	Err  error
	Zone []core.Zone
}
