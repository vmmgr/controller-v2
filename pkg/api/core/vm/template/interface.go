package template

type Root struct {
	Template []Template `json:"template"`
}

type Template struct {
	ID        uint   `json:"id"`
	UUID      string `json:"uuid"`
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

type Result struct {
	Error error `json:"error"`
	Root  Root  `json:"root"`
}
