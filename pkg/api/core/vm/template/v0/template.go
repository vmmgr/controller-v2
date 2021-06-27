package v0

import (
	"github.com/gin-gonic/gin"
	auth "github.com/vmmgr/controller/pkg/api/core/auth/v0"
	"github.com/vmmgr/controller/pkg/api/core/common"
	"github.com/vmmgr/controller/pkg/api/core/vm/template"
	dbTemplate "github.com/vmmgr/controller/pkg/api/store/imacon/template/v0"
	dbStorage "github.com/vmmgr/controller/pkg/api/store/node/storage/v0"
	dbNode "github.com/vmmgr/controller/pkg/api/store/node/v0"
	"net/http"
)

func GetByAdmin(c *gin.Context) {
	resultAdmin := auth.AdminAuthentication(c.Request.Header.Get("ACCESS_TOKEN"))
	if resultAdmin.Err != nil {
		c.JSON(http.StatusUnauthorized, common.Error{Error: resultAdmin.Err.Error()})
		return
	}

	resultTemplate, err := dbTemplate.GetAll()
	if err != nil {
		c.JSON(http.StatusInternalServerError, common.Error{Error: err.Error()})
		return
	}

	resultStorage := dbStorage.GetAll()
	if resultStorage.Err != nil {
		c.JSON(http.StatusInternalServerError, common.Error{Error: resultStorage.Err.Error()})
		return
	}

	c.JSON(http.StatusOK, template.TemplateByAdmin{Template: resultTemplate, Storage: resultStorage.Storage})
}

func Get(c *gin.Context) {
	resultAdmin := auth.AdminAuthentication(c.Request.Header.Get("ACCESS_TOKEN"))
	if resultAdmin.Err != nil {
		c.JSON(http.StatusUnauthorized, common.Error{Error: resultAdmin.Err.Error()})
		return
	}

	resultTemplate, err := dbTemplate.GetAll()
	if err != nil {
		c.JSON(http.StatusInternalServerError, common.Error{Error: err.Error()})
		return
	}

	resultNode := dbNode.GetAll()
	if resultNode.Err != nil {
		c.JSON(http.StatusInternalServerError, common.Error{Error: resultNode.Err.Error()})
		return
	}

	c.JSON(http.StatusOK, template.Template{Template: resultTemplate, Node: resultNode.Node})
}

//func GetTemplate(id, planID uint) (template.Template, template.Plan, error) {
//	file, err := ioutil.ReadFile(config.Conf.Controller.TemplateConfPath)
//	if err != nil {
//		return template.Template{}, template.Plan{}, err
//	}
//	var data template.Root
//	json.Unmarshal(file, &data)
//	for _, tmp := range data.Template {
//		if tmp.ID == id {
//			for _, tmpPlan := range tmp.Plan {
//				if tmpPlan.PlanID == planID {
//					return template.Template{
//							UUID: tmp.UUID, Name: tmp.Name, Tag: tmp.Tag, ImageName: tmp.ImageName,
//						}, template.Plan{
//							PlanID: tmpPlan.PlanID, CPU: tmpPlan.CPU, Mem: tmpPlan.Mem, Storage: tmpPlan.Storage,
//						}, nil
//				}
//			}
//		}
//	}
//	return template.Template{}, template.Plan{}, fmt.Errorf("not found: template ")
//}
