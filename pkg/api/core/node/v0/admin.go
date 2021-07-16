package v0

import (
	"github.com/gin-gonic/gin"
	"github.com/vmmgr/controller/pkg/api/core"
	auth "github.com/vmmgr/controller/pkg/api/core/auth/v0"
	"github.com/vmmgr/controller/pkg/api/core/common"
	"github.com/vmmgr/controller/pkg/api/core/group"
	"github.com/vmmgr/controller/pkg/api/core/node"
	"github.com/vmmgr/controller/pkg/api/core/region"
	"github.com/vmmgr/controller/pkg/api/core/region/zone"
	dbGroup "github.com/vmmgr/controller/pkg/api/store/group/v0"
	dbNode "github.com/vmmgr/controller/pkg/api/store/node/v0"
	dbZone "github.com/vmmgr/controller/pkg/api/store/region/zone/v0"
	"gorm.io/gorm"
	"log"
	"net/http"
	"strconv"
)

func AddAdmin(c *gin.Context) {
	var input core.Node

	resultAdmin := auth.AdminAuthentication(c.Request.Header.Get("ACCESS_TOKEN"))
	if resultAdmin.Err != nil {
		c.JSON(http.StatusUnauthorized, common.Error{Error: resultAdmin.Err.Error()})
		return
	}
	err := c.BindJSON(&input)
	log.Println(err)

	if err = check(input); err != nil {
		c.JSON(http.StatusBadRequest, common.Error{Error: err.Error()})
		return
	}

	if resultZone := dbZone.Get(zone.ID, &core.Zone{Model: gorm.Model{ID: input.ZoneID}}); resultZone.Err != nil {
		c.JSON(http.StatusBadRequest, common.Error{Error: "This zone is not found..."})
		return
	}
	if input.GroupID != nil {
		if resultGroup := dbGroup.Get(group.ID, &core.Group{Model: gorm.Model{ID: *input.GroupID}}); resultGroup.Err != nil {
			c.JSON(http.StatusBadRequest, common.Error{Error: "This group is not found..."})
			return
		}
	}

	if _, err = dbNode.Create(&input); err != nil {
		c.JSON(http.StatusInternalServerError, common.Error{Error: err.Error()})
		return
	}
	c.JSON(http.StatusOK, node.Result{})
}

func DeleteAdmin(c *gin.Context) {
	resultAdmin := auth.AdminAuthentication(c.Request.Header.Get("ACCESS_TOKEN"))
	if resultAdmin.Err != nil {
		c.JSON(http.StatusInternalServerError, common.Error{Error: resultAdmin.Err.Error()})
		return
	}

	var id, err = strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, common.Error{Error: err.Error()})
		return
	}

	if err = dbNode.Delete(&core.Node{Model: gorm.Model{ID: uint(id)}}); err != nil {
		c.JSON(http.StatusInternalServerError, common.Error{Error: err.Error()})
		return
	}
	c.JSON(http.StatusOK, node.Result{})
}

func UpdateAdmin(c *gin.Context) {
	var input core.Node

	resultAdmin := auth.AdminAuthentication(c.Request.Header.Get("ACCESS_TOKEN"))
	if resultAdmin.Err != nil {
		c.JSON(http.StatusUnauthorized, common.Error{Error: resultAdmin.Err.Error()})
		return
	}
	err := c.BindJSON(&input)
	log.Println(err)

	tmp := dbNode.Get(region.ID, &core.Node{Model: gorm.Model{ID: input.ID}})
	if tmp.Err != nil {
		c.JSON(http.StatusInternalServerError, common.Error{Error: tmp.Err.Error()})
		return
	}

	replace, err := updateAdminUser(input, tmp.Node[0])
	if err != nil {
		c.JSON(http.StatusInternalServerError, common.Error{Error: "error: this email is already registered"})
		return
	}

	if err = dbNode.Update(region.UpdateAll, replace); err != nil {
		c.JSON(http.StatusInternalServerError, common.Error{Error: err.Error()})
		return
	}
	c.JSON(http.StatusOK, node.Result{})
}

func GetAdmin(c *gin.Context) {
	resultAdmin := auth.AdminAuthentication(c.Request.Header.Get("ACCESS_TOKEN"))
	if resultAdmin.Err != nil {
		c.JSON(http.StatusUnauthorized, common.Error{Error: resultAdmin.Err.Error()})
		return
	}
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, common.Error{Error: err.Error()})
		return
	}

	result := dbNode.Get(region.ID, &core.Node{Model: gorm.Model{ID: uint(id)}})
	if result.Err != nil {
		c.JSON(http.StatusInternalServerError, common.Error{Error: result.Err.Error()})
		return
	}
	c.JSON(http.StatusOK, node.Result{Node: result.Node})
}

func GetAllAdmin(c *gin.Context) {
	resultAdmin := auth.AdminAuthentication(c.Request.Header.Get("ACCESS_TOKEN"))
	if resultAdmin.Err != nil {
		c.JSON(http.StatusUnauthorized, common.Error{Error: resultAdmin.Err.Error()})
		return
	}

	if result := dbNode.GetAll(); result.Err != nil {
		c.JSON(http.StatusInternalServerError, common.Error{Error: result.Err.Error()})
	} else {
		c.JSON(http.StatusOK, node.Result{Node: result.Node})
	}
}
