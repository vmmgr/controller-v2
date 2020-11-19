package v0

import (
	"encoding/json"
	"github.com/vmmgr/controller/pkg/api/core/tool/config"
	"github.com/vmmgr/controller/pkg/api/core/vm/template"
	"io/ioutil"
)

func GetTemplate(id, planID uint) (template.Template, error) {
	file, err := ioutil.ReadFile(config.Conf.Controller.TemplateConfPath)
	if err != nil {
		return template.Template{}, err
	}
	var data template.Template
	json.Unmarshal(file, &data)
	return data, nil
}
