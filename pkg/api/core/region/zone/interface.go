package zone

import "github.com/jinzhu/gorm"

const (
	ID        = 0
	RegionID  = 1
	Name      = 2
	UpdateAll = 110
)

type Zone struct {
	gorm.Model
	RegionID uint   `json:"region_id"`
	Name     string `json:"name"`
	Postcode string `json:"postcode"`
	Address  string `json:"address"`
	Tel      string `json:"tel"`
	Mail     string `json:"mail"`
	Comment  string `json:"comment"`
	Lock     *bool  `json:"lock"`
}

type Result struct {
	Status bool   `json:"status"`
	Error  string `json:"error"`
	Zone   []Zone `json:"zone"`
}

type ResultOne struct {
	Status bool   `json:"status"`
	Error  string `json:"error"`
	Zone   Zone   `json:"zone"`
}

type ResultDatabase struct {
	Err  error
	Zone []Zone
}
