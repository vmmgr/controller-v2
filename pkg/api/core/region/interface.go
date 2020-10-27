package region

import "github.com/jinzhu/gorm"

const (
	ID        = 0
	Name      = 1
	UpdateAll = 110
)

type Region struct {
	gorm.Model
	Name    string `json:"name"`
	Comment string `json:"comment"`
	Lock    *bool  `json:"lock"`
}

type Result struct {
	Status bool     `json:"status"`
	Error  string   `json:"error"`
	Region []Region `json:"region"`
}

type ResultOne struct {
	Status bool   `json:"status"`
	Error  string `json:"error"`
	Region Region `json:"region"`
}

type ResultDatabase struct {
	Err    error
	Region []Region
}
