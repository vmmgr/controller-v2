package v0

import (
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	"github.com/vmmgr/controller/pkg/api/core"
	authInterface "github.com/vmmgr/controller/pkg/api/core/auth"
	auth "github.com/vmmgr/controller/pkg/api/core/auth/v0"
	"github.com/vmmgr/controller/pkg/api/core/common"
	"github.com/vmmgr/controller/pkg/api/core/token"
	toolToken "github.com/vmmgr/controller/pkg/api/core/tool/token"
	dbToken "github.com/vmmgr/controller/pkg/api/store/token/v0"
	"net/http"
	"strconv"
	"time"
)

func GenerateAdmin(c *gin.Context) {
	resultAuth := auth.AdminRadiusAuthentication(authInterface.AdminStruct{
		User: c.Request.Header.Get("USER"), Pass: c.Request.Header.Get("PASS")})
	if resultAuth.Err != nil {
		c.JSON(http.StatusInternalServerError, common.Error{Error: resultAuth.Err.Error()})
		return
	}
	accessToken, _ := toolToken.Generate(2)

	if err := dbToken.Create(&core.Token{
		UserID: 0, ExpiredAt: time.Now().Add(60 * time.Minute),
		Admin:       &[]bool{true}[0],
		AccessToken: accessToken,
		Debug:       "User: " + c.Request.Header.Get("USER"),
	}); err != nil {
		c.JSON(http.StatusInternalServerError, common.Error{Error: err.Error()})
		return
	}
	c.JSON(http.StatusOK, token.Result{Token: []core.Token{{AccessToken: accessToken}}})
}

func AddAdmin(c *gin.Context) {
	var input core.Token

	resultAdmin := auth.AdminAuthentication(c.Request.Header.Get("ACCESS_TOKEN"))
	if resultAdmin.Err != nil {
		c.JSON(http.StatusInternalServerError, common.Error{Error: resultAdmin.Err.Error()})
		return
	}
	c.BindJSON(&input)

	accessToken, _ := toolToken.Generate(2)

	if err := dbToken.Create(&core.Token{
		Admin:       &[]bool{true}[0],
		AccessToken: accessToken,
		Debug:       "User: " + c.Request.Header.Get("USER"),
	}); err != nil {
		c.JSON(http.StatusInternalServerError, common.Error{Error: err.Error()})
		return
	}
	c.JSON(http.StatusOK, token.ResultTmpToken{Token: accessToken})
}

func DeleteAdmin(c *gin.Context) {
	resultAdmin := auth.AdminAuthentication(c.Request.Header.Get("ACCESS_TOKEN"))
	if resultAdmin.Err != nil {
		c.JSON(http.StatusInternalServerError, common.Error{Error: resultAdmin.Err.Error()})
		return
	}

	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusInternalServerError, common.Error{Error: err.Error()})
		return
	}

	if err = dbToken.Delete(&core.Token{Model: gorm.Model{ID: uint(id)}}); err != nil {
		c.JSON(http.StatusInternalServerError, common.Error{Error: err.Error()})
		return
	}
	c.JSON(http.StatusOK, token.Result{})
}

func DeleteAllAdmin(c *gin.Context) {
	resultAdmin := auth.AdminAuthentication(c.Request.Header.Get("ACCESS_TOKEN"))
	if resultAdmin.Err != nil {
		c.JSON(http.StatusInternalServerError, common.Error{Error: resultAdmin.Err.Error()})
		return
	}

	if err := dbToken.DeleteAll(); err != nil {
		c.JSON(http.StatusInternalServerError, common.Error{Error: err.Error()})
		return
	}
	c.JSON(http.StatusOK, token.Result{})
}

func UpdateAdmin(c *gin.Context) {
	var input core.Token

	resultAdmin := auth.AdminAuthentication(c.Request.Header.Get("ACCESS_TOKEN"))
	if resultAdmin.Err != nil {
		c.JSON(http.StatusInternalServerError, common.Error{Error: resultAdmin.Err.Error()})
		return
	}
	c.BindJSON(&input)

	if err := dbToken.Update(token.UpdateAll, &input); err != nil {
		c.JSON(http.StatusInternalServerError, common.Error{Error: err.Error()})
		return
	}
	c.JSON(http.StatusOK, token.Result{})
}

func GetAdmin(c *gin.Context) {
	resultAdmin := auth.AdminAuthentication(c.Request.Header.Get("ACCESS_TOKEN"))
	if resultAdmin.Err != nil {
		c.JSON(http.StatusInternalServerError, common.Error{Error: resultAdmin.Err.Error()})
		return
	}
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusInternalServerError, common.Error{Error: err.Error()})
		return
	}

	result := dbToken.Get(token.ID, &core.Token{Model: gorm.Model{ID: uint(id)}})
	if result.Err != nil {
		c.JSON(http.StatusInternalServerError, common.Error{Error: err.Error()})
		return
	}
	c.JSON(http.StatusOK, token.Result{Token: result.Token})
}

func GetAllAdmin(c *gin.Context) {
	resultAdmin := auth.AdminAuthentication(c.Request.Header.Get("ACCESS_TOKEN"))
	if resultAdmin.Err != nil {
		c.JSON(http.StatusInternalServerError, common.Error{Error: resultAdmin.Err.Error()})
		return
	}

	if result := dbToken.GetAll(); result.Err != nil {
		c.JSON(http.StatusInternalServerError, common.Error{Error: result.Err.Error()})
	} else {
		c.JSON(http.StatusOK, token.Result{Token: result.Token})
	}
}
