package v0

import (
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	auth "github.com/vmmgr/controller/pkg/api/core/auth/v0"
	"github.com/vmmgr/controller/pkg/api/core/group"
	"github.com/vmmgr/controller/pkg/api/core/node"
	"github.com/vmmgr/controller/pkg/api/core/region"
	"github.com/vmmgr/controller/pkg/api/core/region/zone"
	dbGroup "github.com/vmmgr/controller/pkg/api/store/group/v0"
	dbNode "github.com/vmmgr/controller/pkg/api/store/node/v0"
	dbZone "github.com/vmmgr/controller/pkg/api/store/region/zone/v0"
	"log"
	"net/http"
	"strconv"
)

func AddAdmin(c *gin.Context) {
	var input node.Node

	resultAdmin := auth.AdminAuthentication(c.Request.Header.Get("ACCESS_TOKEN"))
	if resultAdmin.Err != nil {
		c.JSON(http.StatusUnauthorized, node.Result{Status: false, Error: resultAdmin.Err.Error()})
		return
	}
	err := c.BindJSON(&input)
	log.Println(err)

	if err := check(input); err != nil {
		c.JSON(http.StatusBadRequest, node.Result{Status: false, Error: err.Error()})
		return
	}

	if resultZone := dbZone.Get(zone.ID, &zone.Zone{Model: gorm.Model{ID: input.ZoneID}}); resultZone.Err != nil {
		c.JSON(http.StatusBadRequest, node.Result{Status: false, Error: "This zone is not found..."})
		return
	}
	if input.GroupID != 0 {
		if resultGroup := dbGroup.Get(group.ID, &group.Group{Model: gorm.Model{ID: input.GroupID}}); resultGroup.Err != nil {
			c.JSON(http.StatusBadRequest, node.Result{Status: false, Error: "This group is not found..."})
			return
		}
	}

	if _, err := dbNode.Create(&input); err != nil {
		c.JSON(http.StatusInternalServerError, node.Result{Status: false, Error: err.Error()})
		return
	}
	c.JSON(http.StatusOK, node.Result{Status: true})
}

func DeleteAdmin(c *gin.Context) {
	resultAdmin := auth.AdminAuthentication(c.Request.Header.Get("ACCESS_TOKEN"))
	if resultAdmin.Err != nil {
		c.JSON(http.StatusInternalServerError, node.Result{Status: false, Error: resultAdmin.Err.Error()})
		return
	}

	var id, err = strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, node.Result{Status: false, Error: err.Error()})
		return
	}

	if err := dbNode.Delete(&node.Node{Model: gorm.Model{ID: uint(id)}}); err != nil {
		c.JSON(http.StatusInternalServerError, node.Result{Status: false, Error: err.Error()})
		return
	}
	c.JSON(http.StatusOK, node.Result{Status: true})
}

func UpdateAdmin(c *gin.Context) {
	var input node.Node

	resultAdmin := auth.AdminAuthentication(c.Request.Header.Get("ACCESS_TOKEN"))
	if resultAdmin.Err != nil {
		c.JSON(http.StatusUnauthorized, node.Result{Status: false, Error: resultAdmin.Err.Error()})
		return
	}
	err := c.BindJSON(&input)
	log.Println(err)

	tmp := dbNode.Get(region.ID, &node.Node{Model: gorm.Model{ID: input.ID}})
	if tmp.Err != nil {
		c.JSON(http.StatusInternalServerError, node.Result{Status: false, Error: tmp.Err.Error()})
		return
	}

	replace, err := updateAdminUser(input, tmp.Node[0])
	if err != nil {
		c.JSON(http.StatusInternalServerError, node.Result{Status: false, Error: "error: this email is already registered"})
		return
	}

	if err := dbNode.Update(region.UpdateAll, replace); err != nil {
		c.JSON(http.StatusInternalServerError, node.Result{Status: false, Error: err.Error()})
		return
	}
	c.JSON(http.StatusOK, node.Result{Status: true})
}

func GetAdmin(c *gin.Context) {
	resultAdmin := auth.AdminAuthentication(c.Request.Header.Get("ACCESS_TOKEN"))
	if resultAdmin.Err != nil {
		c.JSON(http.StatusUnauthorized, node.Result{Status: false, Error: resultAdmin.Err.Error()})
		return
	}
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, node.Result{Status: false, Error: err.Error()})
		return
	}

	result := dbNode.Get(region.ID, &node.Node{Model: gorm.Model{ID: uint(id)}})
	if result.Err != nil {
		c.JSON(http.StatusInternalServerError, node.Result{Status: false, Error: result.Err.Error()})
		return
	}
	c.JSON(http.StatusOK, node.Result{Status: true, Node: result.Node})
}

func GetAllAdmin(c *gin.Context) {
	resultAdmin := auth.AdminAuthentication(c.Request.Header.Get("ACCESS_TOKEN"))
	if resultAdmin.Err != nil {
		c.JSON(http.StatusUnauthorized, node.Result{Status: false, Error: resultAdmin.Err.Error()})
		return
	}

	if result := dbNode.GetAll(); result.Err != nil {
		c.JSON(http.StatusInternalServerError, node.Result{Status: false, Error: result.Err.Error()})
	} else {
		c.JSON(http.StatusOK, node.Result{Status: true, Node: result.Node})
	}
}
