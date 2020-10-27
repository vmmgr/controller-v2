package storage

import "github.com/jinzhu/gorm"

const (
	ID        = 0
	NodeID    = 1
	AdminOnly = 2
	Name      = 3
	UpdateAll = 110
)

type Storage struct {
	gorm.Model
	NodeID      uint   `json:"node_id"`
	AdminOnly   *bool  `json:"admin"`
	Name        uint   `json:"name"`
	Type        uint   `json:"type"`
	Path        string `json:"path"`
	MaxCapacity uint   `json:"max_capacity"`
	Comment     string `json:"comment"`
}

type Result struct {
	Status  bool      `json:"status"`
	Error   string    `json:"error"`
	Storage []Storage `json:"storage"`
}

type ResultOne struct {
	Status  bool    `json:"status"`
	Error   string  `json:"error"`
	Storage Storage `json:"storage"`
}

type ResultDatabase struct {
	Err     error
	Storage []Storage
}
