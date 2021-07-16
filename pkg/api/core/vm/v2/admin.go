package v2

import (
	"encoding/xml"
	"github.com/gin-gonic/gin"
	"github.com/libvirt/libvirt-go"
	libVirtXml "github.com/libvirt/libvirt-go-xml"
	"github.com/vmmgr/controller/pkg/api/core"
	auth "github.com/vmmgr/controller/pkg/api/core/auth/v0"
	"github.com/vmmgr/controller/pkg/api/core/common"
	"github.com/vmmgr/controller/pkg/api/core/region"
	"github.com/vmmgr/controller/pkg/api/core/tool/config"
	"github.com/vmmgr/controller/pkg/api/core/vm"
	dbNode "github.com/vmmgr/controller/pkg/api/store/node/v0"
	dbVM "github.com/vmmgr/controller/pkg/api/store/vm/v0"
	"gorm.io/gorm"
	"log"
	"net/http"
	"strconv"
)

type VMHandler struct {
	Conn    *libvirt.Connect
	VM      vm.VirtualMachine
	Node    core.Node
	GroupID uint
	IPID    uint
}

func NewVMHandler(input VMHandler) *VMHandler {
	return &VMHandler{
		Conn:    input.Conn,
		VM:      input.VM,
		Node:    input.Node,
		GroupID: input.GroupID,
		IPID:    input.IPID,
	}
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
	node, conn, err := connectLibvirt(input.NodeID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, common.Error{Error: err.Error()})
		return
	}

	if !input.TemplateApply {
		//手動作成時
		//VM作成用のデータ
		h := NewVMHandler(VMHandler{
			Conn: conn,
			VM:   input.VM,
			Node: *node,
		})

		err = h.CreateVM()
		if err != nil {
			c.JSON(http.StatusBadRequest, common.Error{Error: err.Error()})
			return
		}
	} else {
		// storage
		var vmBasePath core.Storage = core.Storage{}
		for _, tmp := range node.Storage {
			if tmp.ID == input.Template.StorageID {
				vmBasePath = tmp
				break
			}
		}
		if vmBasePath.ID == 0 {
			c.JSON(http.StatusBadRequest, common.Error{Error: err.Error()})
			return
		}

		node.Storage = []core.Storage{vmBasePath}

		//----ベースイメージコピー処理----
		h := NewVMAdminTemplateHandler(VMTemplateHandler{
			input:    input.VM,
			template: input.Template,
			node:     *node,
			conn:     conn,
		})

		err = h.templateApply()
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

	uuid := c.Param("uuid")

	nodeID, err := strconv.Atoi(c.Param("node_id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, common.Error{Error: err.Error()})
		return
	}

	_, conn, err := connectLibvirt(uint(nodeID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, common.Error{Error: err.Error()})
		return
	}

	dom, err := conn.LookupDomainByUUIDString(uuid)
	if err != nil {
		c.JSON(http.StatusInternalServerError, common.Error{Error: err.Error()})
		return
	}

	stat, _, err := dom.GetState()
	if err != nil {
		c.JSON(http.StatusInternalServerError, common.Error{Error: err.Error()})
		return
	}
	if stat == libvirt.DOMAIN_SHUTOFF {
		log.Println("power off")
	} else {
		if err = dom.Destroy(); err != nil {
			c.JSON(http.StatusInternalServerError, common.Error{Error: err.Error()})
			return
		}
	}

	if err = dom.Undefine(); err != nil {
		c.JSON(http.StatusInternalServerError, common.Error{Error: err.Error()})
		return
	}

	if err = dom.Free(); err != nil {
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

	//replace, err := updateAdminUser(input, tmp.VMs[0])
	//if err != nil {
	//	c.JSON(http.StatusInternalServerError, common.Error{Error: "error: this email is already registered"})
	//	return
	//}
	//
	//if err = dbVM.Update(region.UpdateAll, replace); err != nil {
	//	c.JSON(http.StatusInternalServerError, common.Error{Error: err.Error()})
	//	return
	//}
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

	_, conn, err := connectLibvirt(uint(nodeID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, common.Error{Error: err.Error()})
		return
	}

	dom, err := conn.LookupDomainByUUIDString(uuid)
	if err != nil {
		c.JSON(http.StatusInternalServerError, common.Error{Error: err.Error()})
		return
	}

	stat, _, err := dom.GetState()
	if err != nil {
		log.Println(err)
		c.JSON(http.StatusInternalServerError, common.Error{Error: err.Error()})
		return
	}

	// 初期定義
	t := libVirtXml.Domain{}

	// XMLをStructに代入
	tmpXml, _ := dom.GetXMLDesc(libvirt.DOMAIN_XML_SECURE)
	xml.Unmarshal([]byte(tmpXml), &t)

	if err = dom.Undefine(); err != nil {
		c.JSON(http.StatusInternalServerError, common.Error{Error: err.Error()})
		return
	}

	if err = dom.Free(); err != nil {
		c.JSON(http.StatusInternalServerError, common.Error{Error: err.Error()})
		return
	}

	c.JSON(http.StatusOK, vm.ResultOneAdmin{Status: http.StatusOK, VM: vm.Detail{
		VM:   t,
		Stat: uint(stat),
		Node: uint(nodeID),
	}})
}

func GetAllAdmin(c *gin.Context) {
	var vms []vm.Detail

	resultNode := dbNode.GetAll()
	if resultNode.Err != nil {
		log.Println(resultNode.Err)
		c.JSON(http.StatusInternalServerError, common.Error{Error: resultNode.Err.Error()})
		return
	}

	for _, tmpNode := range resultNode.Node {
		conn, err := libvirt.NewConnect("qemu+ssh://" + config.Conf.Node.User + "@" + tmpNode.IP + "/system")
		if err != nil {
			log.Println("failed to connect to qemu: " + err.Error())
			return
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
