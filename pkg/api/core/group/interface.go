package group

import (
	"github.com/vmmgr/controller/pkg/api/core"
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

type Result struct {
	Group []core.Group `json:"group"`
}

type ResultOne struct {
	Group core.Group `json:"group"`
}

type ResultAll struct {
	Group core.Group `json:"group"`
}

type ResultDatabase struct {
	Err   error
	Group []core.Group
}
