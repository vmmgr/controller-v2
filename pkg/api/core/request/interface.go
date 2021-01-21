package notice

import (
	"github.com/jinzhu/gorm"
	"time"
)

const (
	ID             = 0
	UUID           = 1
	GroupID        = 2
	ExpirationDate = 3
)

type Request struct {
	gorm.Model
	ExpirationDate time.Time `json:"expiration_date"`
	GroupID        uint      `json:"group_id"`
	UUID           string    `json:"uuid"`
	Comment        string    `json:"comment"`
}

type Result struct {
	Status  bool      `json:"status"`
	Error   string    `json:"error"`
	Request []Request `json:"request"`
}

type ResultDatabase struct {
	Err     error
	Request []Request
}
