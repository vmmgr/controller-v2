package v0

import (
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	"github.com/vmmgr/controller/pkg/api/core"
	auth "github.com/vmmgr/controller/pkg/api/core/auth/v0"
	"github.com/vmmgr/controller/pkg/api/core/common"
	"github.com/vmmgr/controller/pkg/api/core/node"
	"github.com/vmmgr/controller/pkg/api/core/vm"
	dbNode "github.com/vmmgr/controller/pkg/api/store/node/v0"
	"net/http"
)

// #13 Issue
func UserCreate(c *gin.Context) {
	var input vm.Template
	userToken := c.Request.Header.Get("USER_TOKEN")
	accessToken := c.Request.Header.Get("ACCESS_TOKEN")

	if err := c.BindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, common.Error{Error: err.Error()})
		return
	}

	result := auth.GroupAuthentication(0, core.Token{UserToken: userToken, AccessToken: accessToken})
	if result.Err != nil {
		c.JSON(http.StatusUnauthorized, common.Error{Error: result.Err.Error()})
		return
	}

	// nodeIDが存在するか確認
	resultNode := dbNode.Get(node.ID, &core.Node{Model: gorm.Model{ID: input.NodeID}})
	if resultNode.Err != nil {
		c.JSON(http.StatusBadRequest, common.Error{Error: resultNode.Err.Error()})
		return
	}

	//nodeIDの数が0である場合は
	if len(resultNode.Node) == 0 {
		c.JSON(http.StatusNotFound, common.Error{Error: "node id is not found..."})
		return
	}

	//　NodeのGroupIDが0かつAdminOnlyがfalseの時の以外である場合、
	if !(resultNode.Node[0].GroupID == 0 && resultNode.Node[0].AdminOnly == &[]bool{false}[0]) {
		c.JSON(http.StatusForbidden, common.Error{Error: "You can't use this node..."})
		return
	}

	// NodeIDとStoragePathTypeがGroupで使用可能か確認

	//----ベースイメージコピー処理----
	h := NewVMUserTemplateHandler(VMTemplateHandler{
		template: input, node: resultNode.Node[0], authUser: result, admin: false})

	err := h.templateApply()
	if err != nil {
		c.JSON(http.StatusNotFound, common.Error{Error: "template is not found..."})
		return
	}

	c.JSON(http.StatusOK, vm.Result{})
}

func UserDelete(c *gin.Context) {
	//var input vm.Template
	//userToken := c.Request.Header.Get("USER_TOKEN")
	//accessToken := c.Request.Header.Get("ACCESS_TOKEN")
	//
	//if err := c.BindJSON(&input); err != nil {
	//	c.JSON(http.StatusBadRequest, common.Error{ Error: err.Error()})
	//	return
	//}
	//
	//result := auth.GroupAuthentication(core.Token{UserToken: userToken, AccessToken: accessToken})
	//if result.Err != nil {
	//	c.JSON(http.StatusUnauthorized, common.Error{ Error: result.Err.Error()})
	//	return
	//}
	//
	//// nodeIDが存在するか確認
	//resultNode := dbNode.Get(node.ID, &core.Node{Model: gorm.Model{ID: input.NodeID}})
	//if resultNode.Err != nil {
	//	c.JSON(http.StatusBadRequest, common.Error{ Error: resultNode.Err.Error()})
	//	return
	//}
	//
	////nodeIDの数が0である場合は
	//if len(resultNode.Node) == 0 {
	//	c.JSON(http.StatusNotFound, common.Error{ Error: "node id is not found..."})
	//	return
	//}
	//
	//////NodeのGroupIDが0かつAdminOnlyがfalseの時の以外である場合、
	//if !(resultNode.Node[0].GroupID == 0 && resultNode.Node[0].AdminOnly == &[]bool{false}[0]) {
	//	c.JSON(http.StatusForbidden, common.Error{ Error: "You can't use this node..."})
	//	return
	//}
	//c.JSON(http.StatusOK, vm.Result{})
}
