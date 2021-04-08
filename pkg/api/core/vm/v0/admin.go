package v0

import (
	"encoding/json"
	"encoding/xml"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	"github.com/libvirt/libvirt-go"
	libVirtXml "github.com/libvirt/libvirt-go-xml"
	"github.com/vmmgr/controller/pkg/api/core"
	auth "github.com/vmmgr/controller/pkg/api/core/auth/v0"
	"github.com/vmmgr/controller/pkg/api/core/common"
	"github.com/vmmgr/controller/pkg/api/core/node"
	"github.com/vmmgr/controller/pkg/api/core/region"
	"github.com/vmmgr/controller/pkg/api/core/tool/client"
	"github.com/vmmgr/controller/pkg/api/core/tool/config"
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

type VMHandler struct {
	Conn *libvirt.Connect
	//VM   vm.VirtualMachine
}

func NewVMHandler(input VMHandler) *VMHandler {
	return &VMHandler{Conn: input.Conn}
}

// #12 Issue

func AddAdmin(c *gin.Context) {
	var input vm.CreateAdmin

	resultAdmin := auth.AdminAuthentication(c.Request.Header.Get("ACCESS_TOKEN"))
	if resultAdmin.Err != nil {
		c.JSON(http.StatusUnauthorized, common.Error{Error: resultAdmin.Err.Error()})
		return
	}

	if err := c.BindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, common.Error{Error: err.Error()})
		return
	}

	// nodeIDが存在するか確認
	resultNode := dbNode.Get(node.ID, &core.Node{Model: gorm.Model{ID: input.NodeID}})
	if resultNode.Err != nil {
		c.JSON(http.StatusUnauthorized, common.Error{Error: resultNode.Err.Error()})
		return
	}

	if len(resultNode.Node) == 0 {
		c.JSON(http.StatusForbidden, common.Error{Error: "node id is not found..."})
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
				common.Error{Error: err.Error() + "|" + resultVMCreateProcess})
			return
		}

	} else {
		//----ベースイメージコピー処理----
		h := NewVMAdminTemplateHandler(VMTemplateHandler{input: input.VM, template: input.Template, node: resultNode.Node[0]})

		err := h.templateApply()
		if err != nil {
			c.JSON(http.StatusNotFound, common.Error{Error: "template is not found..."})
			return
		}
	}

	c.JSON(http.StatusOK, vm.Result{})
}

func DeleteAdmin(c *gin.Context) {
	var input vm.DeleteAdmin

	resultAdmin := auth.AdminAuthentication(c.Request.Header.Get("ACCESS_TOKEN"))
	if resultAdmin.Err != nil {
		c.JSON(http.StatusUnauthorized, common.Error{Error: resultAdmin.Err.Error()})
		return
	}

	if err := c.BindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, common.Error{Error: err.Error()})
		return
	}

	var id, err = strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, common.Error{Error: err.Error()})
		return
	}
	vmResult := dbVM.Get(vm.ID, &core.VM{Model: gorm.Model{ID: uint(id)}})
	if vmResult.Err != nil {
		c.JSON(http.StatusInternalServerError, common.Error{Error: vmResult.Err.Error()})
		return
	}

	//nodeResult := dbNode.Get(node.ID, &core.Node{Model: gorm.Model{ID: vmResult.VMs[0].NodeID}})

	//client.Delete(
	//	"http://" + nodeResult.Node[0].IP + ":" + strconv.Itoa(int(nodeResult.Node[0].Port)) + "/api/v1/vm/" + vmResult.VMs[0].ID,
	//
	//)

	if err = dbVM.Delete(&core.VM{Model: gorm.Model{ID: uint(id)}}); err != nil {
		c.JSON(http.StatusInternalServerError, common.Error{Error: err.Error()})
		return
	}
	c.JSON(http.StatusOK, vm.Result{})
}

func UpdateAdmin(c *gin.Context) {
	var input core.VM

	resultAdmin := auth.AdminAuthentication(c.Request.Header.Get("ACCESS_TOKEN"))
	if resultAdmin.Err != nil {
		c.JSON(http.StatusUnauthorized, common.Error{Error: resultAdmin.Err.Error()})
		return
	}
	err := c.BindJSON(&input)
	log.Println(err)

	tmp := dbVM.Get(region.ID, &core.VM{Model: gorm.Model{ID: input.ID}})
	if tmp.Err != nil {
		c.JSON(http.StatusInternalServerError, common.Error{Error: tmp.Err.Error()})
		return
	}

	replace, err := updateAdminUser(input, tmp.VMs[0])
	if err != nil {
		c.JSON(http.StatusInternalServerError, common.Error{Error: "error: this email is already registered"})
		return
	}

	if err = dbVM.Update(region.UpdateAll, replace); err != nil {
		c.JSON(http.StatusInternalServerError, common.Error{Error: err.Error()})
		return
	}
	c.JSON(http.StatusOK, vm.Result{})
}

func GetAdmin(c *gin.Context) {
	resultAdmin := auth.AdminAuthentication(c.Request.Header.Get("ACCESS_TOKEN"))
	if resultAdmin.Err != nil {
		c.JSON(http.StatusUnauthorized, common.Error{Error: resultAdmin.Err.Error()})
		return
	}
	nodeID, err := strconv.Atoi(c.Param("node_id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, common.Error{Error: err.Error()})
		return
	}

	uuid := c.Param("uuid")

	resultNode := dbNode.Get(node.ID, &core.Node{Model: gorm.Model{ID: uint(nodeID)}})
	if resultNode.Err != nil {
		c.JSON(http.StatusInternalServerError, common.Error{Error: resultNode.Err.Error()})
		return
	}

	res, err := Get(resultNode.Node[0].IP, resultNode.Node[0].Port, uuid)
	if err != nil {
		log.Println(err)
		c.JSON(http.StatusInternalServerError, common.Error{Error: err.Error()})
		return
	}

	c.JSON(http.StatusOK, vm.ResultOneAdmin{Status: http.StatusOK, VM: vm.Detail{
		VM:   res.Data.VM.VM,
		Stat: res.Data.VM.Stat,
		Node: uint(nodeID),
	}})
}

func GetAllAdmin(c *gin.Context) {
	// Todo: websocketで処理をするべきかも
	var vms []vm.Detail

	resultNode := dbNode.GetAll()
	if resultNode.Err != nil {
		log.Println(resultNode.Err)
		c.JSON(http.StatusInternalServerError, common.Error{Error: resultNode.Err.Error()})
		return
	}

	for _, tmpNode := range resultNode.Node {
		log.Println("qemu+ssh://" + config.Conf.Node.User + "@" + tmpNode.IP + "/system")
		//libvirt.NewConnectWithAuth()
		conn, err := libvirt.NewConnect("qemu+ssh://" + config.Conf.Node.User + "@" + tmpNode.IP + "/system")
		if err != nil {
			log.Println("failed to connect to qemu: " + err.Error())
			//return
		}

		defer conn.Close()

		doms, err := conn.ListAllDomains(libvirt.CONNECT_LIST_DOMAINS_ACTIVE | libvirt.CONNECT_LIST_DOMAINS_INACTIVE)
		log.Println(doms)
		if err != nil {
			log.Println(err)
			c.JSON(http.StatusInternalServerError, common.Error{Error: err.Error()})
		}

		for _, dom := range doms {
			t := libVirtXml.Domain{}
			stat, _, _ := dom.GetState()
			xmlString, _ := dom.GetXMLDesc(libvirt.DOMAIN_XML_SECURE)
			xml.Unmarshal([]byte(xmlString), &t)

			vms = append(vms, vm.Detail{
				VM:   t,
				Stat: uint(stat),
			})
		}
	}

	c.JSON(http.StatusOK, vm.ResultAdmin{VM: vms})
}

//func GetAllAdmin(c *gin.Context) {
//	var allVMs []vm.Detail
//
//	resultAdmin := auth.AdminAuthentication(c.Request.Header.Get("ACCESS_TOKEN"))
//	if resultAdmin.Err != nil {
//		c.JSON(http.StatusUnauthorized, common.Error{ Error: resultAdmin.Err.Error()})
//		return
//	}
//
//	result := dbNode.GetAll()
//	if result.Err != nil {
//		c.JSON(http.StatusInternalServerError, common.Error{ Error: result.Err.Error()})
//		return
//	}
//
//	for _, node := range result.Node {
//		var res nodeAllVMResponse
//
//		response, err := client.Get("http://"+node.IP+":"+strconv.Itoa(int(node.Port))+"/api/v1/vm", "")
//		if err == nil {
//			if json.Unmarshal([]byte(response), &res) != nil {
//				c.JSON(http.StatusInternalServerError, common.Error{ Error: err.Error()})
//				return
//			}
//
//			for _, virtualMachine := range res.Data.VM {
//				allVMs = append(allVMs, vm.Detail{VM: virtualMachine.VM, Stat: virtualMachine.Stat, Node: node.ID})
//			}
//		}
//	}
//
//	c.JSON(http.StatusOK, vm.ResultAdmin{Status: 200, VM: allVMs})
//}

func (h *VMHandler) TestAdminGetAll() ([]vm.Detail, error) {
	var vms []vm.Detail

	doms, err := h.Conn.ListAllDomains(libvirt.CONNECT_LIST_DOMAINS_ACTIVE | libvirt.CONNECT_LIST_DOMAINS_INACTIVE)
	log.Println(doms)
	if err != nil {
		log.Println(err)
		//json.ResponseError(c, http.StatusInternalServerError, err)
		return vms, err
	}

	for _, dom := range doms {
		t := libVirtXml.Domain{}
		stat, _, _ := dom.GetState()
		xmlString, _ := dom.GetXMLDesc(libvirt.DOMAIN_XML_SECURE)
		xml.Unmarshal([]byte(xmlString), &t)

		//log.Println(len(t.Devices.Graphics))
		//log.Println(t.Devices.Graphics[0].VNC.Port)
		//log.Println(t.Devices.Graphics)
		vms = append(vms, vm.Detail{
			VM:   t,
			Stat: uint(stat),
		})
	}

	return vms, nil
}
