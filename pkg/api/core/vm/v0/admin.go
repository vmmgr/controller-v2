package v0

import (
	"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	auth "github.com/vmmgr/controller/pkg/api/core/auth/v0"
	"github.com/vmmgr/controller/pkg/api/core/node"
	"github.com/vmmgr/controller/pkg/api/core/region"
	"github.com/vmmgr/controller/pkg/api/core/tool/client"
	"github.com/vmmgr/controller/pkg/api/core/vm"
	dbNode "github.com/vmmgr/controller/pkg/api/store/node/v0"
	dbVM "github.com/vmmgr/controller/pkg/api/store/vm/v0"
	"log"
	"net/http"
	"strconv"
)

type nodeOneVMResponse struct {
	Status int `json:"status"`
	Data   struct {
		VM vm.Detail `json:"vm"`
	} `json:"data"`
}

type nodeAllVMResponse struct {
	Status int `json:"status"`
	Data   struct {
		VM []vm.Detail `json:"vm"`
	} `json:"data"`
}

// #12 Issue

func AddAdmin(c *gin.Context) {
	var input vm.CreateAdmin

	resultAdmin := auth.AdminAuthentication(c.Request.Header.Get("ACCESS_TOKEN"))
	if resultAdmin.Err != nil {
		c.JSON(http.StatusUnauthorized, vm.Result{Status: false, Error: resultAdmin.Err.Error()})
		return
	}

	if err := c.BindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, vm.Result{Status: false, Error: err.Error()})
		return
	}

	// nodeIDが存在するか確認
	resultNode := dbNode.Get(node.ID, &node.Node{Model: gorm.Model{ID: input.NodeID}})
	if resultNode.Err != nil {
		c.JSON(http.StatusUnauthorized, vm.Result{Status: false, Error: resultNode.Err.Error()})
		return
	}

	if len(resultNode.Node) == 0 {
		c.JSON(http.StatusForbidden, vm.Result{Status: false, Error: "node id is not found..."})
		return
	}

	if !input.TemplateApply {
		//手動作成時
		//VM作成用のデータ
		body, _ := json.Marshal(input.VM)

		resultVMCreateProcess, err := client.Post(
			"http://"+resultNode.Node[0].IP+":"+strconv.Itoa(int(resultNode.Node[0].Port))+"/api/v1/vm",
			body)

		log.Println(resultVMCreateProcess)

		if err != nil {
			log.Println(err)
			c.JSON(http.StatusInternalServerError,
				vm.Result{Status: false, Error: err.Error() + "|" + resultVMCreateProcess})
			return
		}

	} else {
		//----ベースイメージコピー処理----
		h := NewVMAdminTemplateHandler(VMTemplateHandler{input: input.Template, node: resultNode.Node[0]})

		err := h.templateApply()
		if err != nil {
			c.JSON(http.StatusNotFound, vm.Result{Status: false, Error: "template is not found..."})
			return
		}
	}

	c.JSON(http.StatusOK, vm.Result{Status: true})
}

func DeleteAdmin(c *gin.Context) {
	var input vm.DeleteAdmin

	resultAdmin := auth.AdminAuthentication(c.Request.Header.Get("ACCESS_TOKEN"))
	if resultAdmin.Err != nil {
		c.JSON(http.StatusUnauthorized, vm.Result{Status: false, Error: resultAdmin.Err.Error()})
		return
	}

	if err := c.BindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, vm.Result{Status: false, Error: err.Error()})
		return
	}

	var id, err = strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, vm.Result{Status: false, Error: err.Error()})
		return
	}
	vmResult := dbVM.Get(vm.ID, &vm.VM{Model: gorm.Model{ID: uint(id)}})
	if vmResult.Err != nil {
		c.JSON(http.StatusInternalServerError, vm.Result{Status: false, Error: vmResult.Err.Error()})
		return
	}

	//nodeResult := dbNode.Get(node.ID, &node.Node{Model: gorm.Model{ID: vmResult.VMs[0].NodeID}})

	//client.Delete(
	//	"http://" + nodeResult.Node[0].IP + ":" + strconv.Itoa(int(nodeResult.Node[0].Port)) + "/api/v1/vm/" + vmResult.VMs[0].ID,
	//
	//)

	if err := dbVM.Delete(&vm.VM{Model: gorm.Model{ID: uint(id)}}); err != nil {
		c.JSON(http.StatusInternalServerError, vm.Result{Status: false, Error: err.Error()})
		return
	}
	c.JSON(http.StatusOK, vm.Result{Status: true})
}

func UpdateAdmin(c *gin.Context) {
	var input vm.VM

	resultAdmin := auth.AdminAuthentication(c.Request.Header.Get("ACCESS_TOKEN"))
	if resultAdmin.Err != nil {
		c.JSON(http.StatusUnauthorized, vm.Result{Status: false, Error: resultAdmin.Err.Error()})
		return
	}
	err := c.BindJSON(&input)
	log.Println(err)

	tmp := dbVM.Get(region.ID, &vm.VM{Model: gorm.Model{ID: input.ID}})
	if tmp.Err != nil {
		c.JSON(http.StatusInternalServerError, vm.Result{Status: false, Error: tmp.Err.Error()})
		return
	}

	replace, err := updateAdminUser(input, tmp.VMs[0])
	if err != nil {
		c.JSON(http.StatusInternalServerError, vm.Result{Status: false, Error: "error: this email is already registered"})
		return
	}

	if err := dbVM.Update(region.UpdateAll, replace); err != nil {
		c.JSON(http.StatusInternalServerError, vm.Result{Status: false, Error: err.Error()})
		return
	}
	c.JSON(http.StatusOK, vm.Result{Status: true})
}

func GetAdmin(c *gin.Context) {
	resultAdmin := auth.AdminAuthentication(c.Request.Header.Get("ACCESS_TOKEN"))
	if resultAdmin.Err != nil {
		c.JSON(http.StatusUnauthorized, vm.Result{Status: false, Error: resultAdmin.Err.Error()})
		return
	}
	nodeID, err := strconv.Atoi(c.Param("node_id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, vm.Result{Status: false, Error: err.Error()})
		return
	}

	uuid := c.Param("uuid")

	resultNode := dbNode.Get(node.ID, &node.Node{Model: gorm.Model{ID: uint(nodeID)}})
	if resultNode.Err != nil {
		c.JSON(http.StatusInternalServerError, vm.Result{Status: false, Error: resultNode.Err.Error()})
		return
	}

	res, err := Get(resultNode.Node[0].IP, resultNode.Node[0].Port, uuid)
	if err != nil {
		log.Println(err)
		c.JSON(http.StatusInternalServerError, vm.Result{Status: false, Error: err.Error()})
		return
	}

	c.JSON(http.StatusOK, vm.ResultOneAdmin{Status: http.StatusOK, VM: vm.Detail{
		VM: res.Data.VM.VM, Stat: res.Data.VM.Stat, Node: uint(nodeID)}})
}

func GetAllAdmin(c *gin.Context) {
	var allVMs []vm.Detail

	resultAdmin := auth.AdminAuthentication(c.Request.Header.Get("ACCESS_TOKEN"))
	if resultAdmin.Err != nil {
		c.JSON(http.StatusUnauthorized, vm.Result{Status: false, Error: resultAdmin.Err.Error()})
		return
	}

	result := dbNode.GetAll()
	if result.Err != nil {
		c.JSON(http.StatusInternalServerError, vm.Result{Status: false, Error: result.Err.Error()})
		return
	}

	for _, node := range result.Node {
		var res nodeAllVMResponse

		response, err := client.Get("http://"+node.IP+":"+strconv.Itoa(int(node.Port))+"/api/v1/vm", "")
		if err == nil {
			if json.Unmarshal([]byte(response), &res) != nil {
				c.JSON(http.StatusInternalServerError, vm.Result{Status: false, Error: err.Error()})
				return
			}

			for _, virtualMachine := range res.Data.VM {
				allVMs = append(allVMs, vm.Detail{VM: virtualMachine.VM, Stat: virtualMachine.Stat, Node: node.ID})
			}
		}
	}

	c.JSON(http.StatusOK, vm.ResultAdmin{Status: 200, VM: allVMs})
}
