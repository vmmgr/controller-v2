package v0

import (
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	auth "github.com/vmmgr/controller/pkg/api/core/auth/v0"
	"github.com/vmmgr/controller/pkg/api/core/region"
	"github.com/vmmgr/controller/pkg/api/core/region/zone"
	dbRegion "github.com/vmmgr/controller/pkg/api/store/region/v0"
	dbZone "github.com/vmmgr/controller/pkg/api/store/region/zone/v0"
	"log"
	"net/http"
	"strconv"
)

func AddAdmin(c *gin.Context) {
	var input zone.Zone

	resultAdmin := auth.AdminAuthentication(c.Request.Header.Get("ACCESS_TOKEN"))
	if resultAdmin.Err != nil {
		c.JSON(http.StatusUnauthorized, zone.Result{Status: false, Error: resultAdmin.Err.Error()})
		return
	}
	err := c.BindJSON(&input)
	log.Println(err)

	resultRegion := dbRegion.Get(region.ID, &region.Region{Model: gorm.Model{ID: input.ID}})
	if resultRegion.Err != nil {
		c.JSON(http.StatusBadRequest, zone.Result{Status: false, Error: resultRegion.Err.Error()})
		return
	}
	if len(resultRegion.Region) != 0 {
		c.JSON(http.StatusBadRequest, zone.Result{Status: false, Error: "same id !!"})
		return
	}

	if err := check(input); err != nil {
		c.JSON(http.StatusBadRequest, zone.Result{Status: false, Error: err.Error()})
		return
	}

	if _, err := dbZone.Create(&input); err != nil {
		c.JSON(http.StatusInternalServerError, zone.Result{Status: false, Error: err.Error()})
		return
	}
	c.JSON(http.StatusOK, zone.Result{Status: true})
}

func DeleteAdmin(c *gin.Context) {
	resultAdmin := auth.AdminAuthentication(c.Request.Header.Get("ACCESS_TOKEN"))
	if resultAdmin.Err != nil {
		c.JSON(http.StatusInternalServerError, zone.Result{Status: false, Error: resultAdmin.Err.Error()})
		return
	}

	var id, err = strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, zone.Result{Status: false, Error: err.Error()})
		return
	}

	if err := dbZone.Delete(&zone.Zone{Model: gorm.Model{ID: uint(id)}}); err != nil {
		c.JSON(http.StatusInternalServerError, zone.Result{Status: false, Error: err.Error()})
		return
	}
	c.JSON(http.StatusOK, zone.Result{Status: true})
}

func UpdateAdmin(c *gin.Context) {
	var input zone.Zone

	resultAdmin := auth.AdminAuthentication(c.Request.Header.Get("ACCESS_TOKEN"))
	if resultAdmin.Err != nil {
		c.JSON(http.StatusUnauthorized, zone.Result{Status: false, Error: resultAdmin.Err.Error()})
		return
	}
	err := c.BindJSON(&input)
	log.Println(err)

	tmp := dbZone.Get(region.ID, &zone.Zone{Model: gorm.Model{ID: input.ID}})
	if tmp.Err != nil {
		c.JSON(http.StatusInternalServerError, zone.Result{Status: false, Error: tmp.Err.Error()})
		return
	}

	replace, err := updateAdminUser(input, tmp.Zone[0])
	if err != nil {
		c.JSON(http.StatusInternalServerError, zone.Result{Status: false, Error: "error: this email is already registered"})
		return
	}

	if err := dbZone.Update(region.UpdateAll, replace); err != nil {
		c.JSON(http.StatusInternalServerError, zone.Result{Status: false, Error: err.Error()})
		return
	}
	c.JSON(http.StatusOK, zone.Result{Status: true})
}

func GetAdmin(c *gin.Context) {
	resultAdmin := auth.AdminAuthentication(c.Request.Header.Get("ACCESS_TOKEN"))
	if resultAdmin.Err != nil {
		c.JSON(http.StatusUnauthorized, zone.Result{Status: false, Error: resultAdmin.Err.Error()})
		return
	}
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, zone.Result{Status: false, Error: err.Error()})
		return
	}

	result := dbZone.Get(region.ID, &zone.Zone{Model: gorm.Model{ID: uint(id)}})
	if result.Err != nil {
		c.JSON(http.StatusInternalServerError, zone.Result{Status: false, Error: result.Err.Error()})
		return
	}
	c.JSON(http.StatusOK, zone.Result{Status: true, Zone: result.Zone})
}

func GetAllAdmin(c *gin.Context) {
	resultAdmin := auth.AdminAuthentication(c.Request.Header.Get("ACCESS_TOKEN"))
	if resultAdmin.Err != nil {
		c.JSON(http.StatusUnauthorized, zone.Result{Status: false, Error: resultAdmin.Err.Error()})
		return
	}

	if result := dbZone.GetAll(); result.Err != nil {
		c.JSON(http.StatusInternalServerError, zone.Result{Status: false, Error: result.Err.Error()})
	} else {
		c.JSON(http.StatusOK, zone.Result{Status: true, Zone: result.Zone})
	}
}
