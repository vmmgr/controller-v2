package template

import "github.com/vmmgr/controller/pkg/api/core"

type TemplateByAdmin struct {
	Template []core.Template `json:"template"`
	Storage  []core.Storage  `json:"storage"`
	Node     []core.Node     `json:"node"`
}

type Template struct {
	Template []core.Template `json:"template"`
	Storage  []core.Storage  `json:"storage"`
	Node     []core.Node     `json:"node"`
}

type Image struct {
	ID        uint   `json:"id"`
	UUID      string `json:"uuid"`
	Name      string `json:"name"`
	Tag       string `json:"tag"`
	ImageName string `json:"image_name"`
	Plan      []Plan `json:"plan"`
}

//type Template struct {
//	ID        uint   `json:"id"`
//	UUID      string `json:"uuid"`
//	Name      string `json:"name"`
//	Tag       string `json:"tag"`
//	ImageName string `json:"image_name"`
//	Plan      []Plan `json:"plan"`
//}

type Plan struct {
	PlanID  uint `json:"plan_id"`
	CPU     uint `json:"cpu"`
	Mem     uint `json:"mem"`
	Storage uint `json:"storage"`
}
