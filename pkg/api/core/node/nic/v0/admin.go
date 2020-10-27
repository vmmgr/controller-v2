package v0

import (
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	auth "github.com/vmmgr/controller/pkg/api/core/auth/v0"
	"github.com/vmmgr/controller/pkg/api/core/node/nic"
	"github.com/vmmgr/controller/pkg/api/core/region"
	dbNIC "github.com/vmmgr/controller/pkg/api/store/node/nic/v0"
	"log"
	"net/http"
	"strconv"
)

func AddAdmin(c *gin.Context) {
	var input nic.NIC

	resultAdmin := auth.AdminAuthentication(c.Request.Header.Get("ACCESS_TOKEN"))
	if resultAdmin.Err != nil {
		c.JSON(http.StatusUnauthorized, nic.Result{Status: false, Error: resultAdmin.Err.Error()})
		return
	}
	err := c.BindJSON(&input)
	log.Println(err)

	if err := check(input); err != nil {
		c.JSON(http.StatusBadRequest, nic.Result{Status: false, Error: err.Error()})
		return
	}

	if _, err := dbNIC.Create(&input); err != nil {
		c.JSON(http.StatusInternalServerError, nic.Result{Status: false, Error: err.Error()})
		return
	}
	c.JSON(http.StatusOK, nic.Result{Status: true})
}

func DeleteAdmin(c *gin.Context) {
	resultAdmin := auth.AdminAuthentication(c.Request.Header.Get("ACCESS_TOKEN"))
	if resultAdmin.Err != nil {
		c.JSON(http.StatusInternalServerError, nic.Result{Status: false, Error: resultAdmin.Err.Error()})
		return
	}

	var id, err = strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, nic.Result{Status: false, Error: err.Error()})
		return
	}

	if err := dbNIC.Delete(&nic.NIC{Model: gorm.Model{ID: uint(id)}}); err != nil {
		c.JSON(http.StatusInternalServerError, nic.Result{Status: false, Error: err.Error()})
		return
	}
	c.JSON(http.StatusOK, nic.Result{Status: true})
}

func UpdateAdmin(c *gin.Context) {
	var input nic.NIC

	resultAdmin := auth.AdminAuthentication(c.Request.Header.Get("ACCESS_TOKEN"))
	if resultAdmin.Err != nil {
		c.JSON(http.StatusUnauthorized, nic.Result{Status: false, Error: resultAdmin.Err.Error()})
		return
	}
	err := c.BindJSON(&input)
	log.Println(err)

	tmp := dbNIC.Get(region.ID, &nic.NIC{Model: gorm.Model{ID: input.ID}})
	if tmp.Err != nil {
		c.JSON(http.StatusInternalServerError, nic.Result{Status: false, Error: tmp.Err.Error()})
		return
	}

	replace, err := updateAdminUser(input, tmp.NIC[0])
	if err != nil {
		c.JSON(http.StatusInternalServerError, nic.Result{Status: false, Error: "error: this email is already registered"})
		return
	}

	if err := dbNIC.Update(region.UpdateAll, replace); err != nil {
		c.JSON(http.StatusInternalServerError, nic.Result{Status: false, Error: err.Error()})
		return
	}
	c.JSON(http.StatusOK, nic.Result{Status: true})
}

func GetAdmin(c *gin.Context) {
	resultAdmin := auth.AdminAuthentication(c.Request.Header.Get("ACCESS_TOKEN"))
	if resultAdmin.Err != nil {
		c.JSON(http.StatusUnauthorized, nic.Result{Status: false, Error: resultAdmin.Err.Error()})
		return
	}
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, nic.Result{Status: false, Error: err.Error()})
		return
	}

	result := dbNIC.Get(region.ID, &nic.NIC{Model: gorm.Model{ID: uint(id)}})
	if result.Err != nil {
		c.JSON(http.StatusInternalServerError, nic.Result{Status: false, Error: result.Err.Error()})
		return
	}
	c.JSON(http.StatusOK, nic.Result{Status: true, NIC: result.NIC})
}

func GetAllAdmin(c *gin.Context) {
	resultAdmin := auth.AdminAuthentication(c.Request.Header.Get("ACCESS_TOKEN"))
	if resultAdmin.Err != nil {
		c.JSON(http.StatusUnauthorized, nic.Result{Status: false, Error: resultAdmin.Err.Error()})
		return
	}

	if result := dbNIC.GetAll(); result.Err != nil {
		c.JSON(http.StatusInternalServerError, nic.Result{Status: false, Error: result.Err.Error()})
	} else {
		c.JSON(http.StatusOK, nic.Result{Status: true, NIC: result.NIC})
	}
}
