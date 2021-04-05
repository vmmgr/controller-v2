package notice

import "github.com/vmmgr/controller/pkg/api/core"

const (
	ID               = 0
	UserID           = 1
	GroupID          = 2
	UserIDAndGroupID = 3
	Everyone         = 4
	Data             = 5
	Important        = 10
	Fault            = 11
	Info             = 12
	UpdateAll        = 110
)

type Result struct {
	Status bool          `json:"status"`
	Error  string        `json:"error"`
	Notice []core.Notice `json:"notice"`
}

type ResultDatabase struct {
	Err    error
	Notice []core.Notice
}
