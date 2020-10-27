package node

import "github.com/jinzhu/gorm"

const (
	ID        = 0
	ZoneID    = 1
	GroupID   = 2
	AdminOnly = 3
	Name      = 4
	UpdateAll = 110
)

type Node struct {
	gorm.Model
	ZoneID    uint   `json:"zone_id"`
	GroupID   uint   `json:"group_ids"`
	AdminOnly *bool  `json:"admin_only"`
	Name      string `json:"name"`
	IP        string `json:"ip"`
	Port      uint   `json:"port"`
	WsPort    uint   `json:"ws_port"`
	ManageNet uint   `json:"manage_net"`
	Comment   string `json:"comment"`
}

type Result struct {
	Status bool   `json:"status"`
	Error  string `json:"error"`
	Node   []Node `json:"node"`
}

type ResultOne struct {
	Status bool   `json:"status"`
	Error  string `json:"error"`
	Node   Node   `json:"node"`
}

type ResultDatabase struct {
	Err  error
	Node []Node
}
