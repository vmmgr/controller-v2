package v0

import (
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	"github.com/vmmgr/controller/pkg/api/core"
	auth "github.com/vmmgr/controller/pkg/api/core/auth/v0"
	"github.com/vmmgr/controller/pkg/api/core/common"
	"github.com/vmmgr/controller/pkg/api/core/node"
	"github.com/vmmgr/controller/pkg/api/core/node/usb"
	dbNode "github.com/vmmgr/controller/pkg/api/store/node/v0"
	"net/http"
	"strconv"
)

func Get(c *gin.Context) {

	id, _ := strconv.Atoi(c.Param("id"))

	resultAdmin := auth.AdminAuthentication(c.Request.Header.Get("ACCESS_TOKEN"))
	if resultAdmin.Err != nil {
		c.JSON(http.StatusUnauthorized, common.Error{Error: resultAdmin.Err.Error()})
		return
	}

	resultNode := dbNode.Get(node.ID, &core.Node{Model: gorm.Model{ID: uint(id)}})
	if resultNode.Err != nil {
		c.JSON(http.StatusInternalServerError, common.Error{Error: resultNode.Err.Error()})
		return
	}
	response, err := httpRequest(resultNode.Node[0].IP, resultNode.Node[0].Port)
	if err != nil {
		c.JSON(http.StatusInternalServerError, common.Error{Error: err.Error()})
		return
	}
	c.JSON(http.StatusOK, usb.Result{USB: response})

}
