package v0

import (
	"encoding/json"
	"fmt"
	"github.com/vmmgr/controller/pkg/api/core/tool/config"
	"github.com/vmmgr/controller/pkg/api/core/vm/template"
	"io/ioutil"
)

func GetTemplate(id, planID uint) (template.Template, template.Plan, error) {
	file, err := ioutil.ReadFile(config.Conf.Controller.TemplateConfPath)
	if err != nil {
		return template.Template{}, template.Plan{}, err
	}
	var data template.Root
	json.Unmarshal(file, &data)
	for _, tmp := range data.Template {
		if tmp.ID == id {
			for _, tmpPlan := range tmp.Plan {
				if tmpPlan.PlanID == planID {
					return template.Template{
							UUID: tmp.UUID, Name: tmp.Name, Tag: tmp.Tag, ImageName: tmp.ImageName,
						}, template.Plan{
							PlanID: tmpPlan.PlanID, CPU: tmpPlan.CPU, Mem: tmpPlan.Mem, Storage: tmpPlan.Storage,
						}, nil
				}
			}
		}
	}
	return template.Template{}, template.Plan{}, fmt.Errorf("not found: template ")
}
