package storage

import (
	"github.com/jinzhu/gorm"
)

const (
	ID            = 0
	NodeStorageID = 1
	GroupID       = 2
	Name          = 3
	NodeSAndVMID  = 4
	Lock          = 5
	UpdateName    = 100
	UpdateNodeS   = 101
	UpdateGroup   = 102
	UpdateAll     = 110
)

type Storage struct {
	gorm.Model
	VMID          uint   `json:"vm_id"`
	NodeStorageID uint   `json:"node_storage_id"`
	GroupID       uint   `json:"group_id"`
	Name          string `json:"name"`
	Type          uint   `json:"type"`
	FileType      uint   `json:"file_type"`
	MaxCapacity   uint   `json:"max_capacity"`
	UUID          string `json:"path"`
	ReadOnly      *bool  `json:"readonly"`
	Comment       string `json:"comment"`
	Lock          *bool  `json:"lock"`
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
