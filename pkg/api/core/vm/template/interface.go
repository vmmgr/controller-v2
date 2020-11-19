package template

type Template struct {
	ID        uint   `json:"id"`
	Name      string `json:"name"`
	Tag       string `json:"tag"`
	ImageName string `json:"image_name"`
	Plan      []Plan `json:"plan"`
}

type Plan struct {
	PlanID  uint `json:"plan_id"`
	CPU     uint `json:"cpu"`
	Mem     uint `json:"mem"`
	Storage uint `json:"storage"`
}
