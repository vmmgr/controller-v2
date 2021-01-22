package v0

import (
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	auth "github.com/vmmgr/controller/pkg/api/core/auth/v0"
	"github.com/vmmgr/controller/pkg/api/core/node"
	"github.com/vmmgr/controller/pkg/api/core/node/pci"
	dbNode "github.com/vmmgr/controller/pkg/api/store/node/v0"
	"log"
	"net/http"
	"strconv"
)

func Get(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))

	resultAdmin := auth.AdminAuthentication(c.Request.Header.Get("ACCESS_TOKEN"))
	if resultAdmin.Err != nil {
		c.JSON(http.StatusUnauthorized, pci.Result{Status: false, Error: resultAdmin.Err.Error()})
		return
	}

	resultNode := dbNode.Get(node.ID, &node.Node{Model: gorm.Model{ID: uint(id)}})
	if resultNode.Err != nil {
		log.Println(resultNode.Err)
		c.JSON(http.StatusInternalServerError, pci.Result{Status: false, Error: resultNode.Err.Error()})
		return
	}
	response, err := httpRequest(resultNode.Node[0].IP, resultNode.Node[0].Port)
	if err != nil {
		log.Println(err)
		c.JSON(http.StatusInternalServerError, pci.Result{Status: false, Error: err.Error()})
		return
	}

	c.JSON(http.StatusOK, pci.Result{Status: true, PCI: response})
}
