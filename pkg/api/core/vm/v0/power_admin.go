package v0

import (
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	"github.com/libvirt/libvirt-go"
	"github.com/vmmgr/controller/pkg/api/core"
	"github.com/vmmgr/controller/pkg/api/core/common"
	"github.com/vmmgr/controller/pkg/api/core/node"
	"github.com/vmmgr/controller/pkg/api/core/vm"
	dbNode "github.com/vmmgr/controller/pkg/api/store/node/v0"
	"log"
	"net/http"
	"strconv"
)

func StartupAdmin(c *gin.Context) {
	vmUUID := c.Param("vm_uuid")
	nodeID, err := strconv.Atoi(c.Param("node_id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, common.Error{Error: err.Error()})
		return
	}
	resultNode := dbNode.Get(node.ID, &core.Node{Model: gorm.Model{ID: uint(nodeID)}})
	if resultNode.Err != nil {
		c.JSON(http.StatusBadRequest, common.Error{Error: resultNode.Err.Error()})
		return
	}

	log.Println("qemu+ssh://" + resultNode.Node[0].UserName + "@" + resultNode.Node[0].IP + "/system")
	conn, err := libvirt.NewConnect("qemu+ssh://" + resultNode.Node[0].UserName + "@" + resultNode.Node[0].IP + "/system")
	if err != nil {
		log.Println("failed to connect to qemu: " + err.Error())
		c.JSON(http.StatusInternalServerError, common.Error{Error: err.Error()})
		return
	}

	defer conn.Close()

	dom, err := conn.LookupDomainByUUIDString(vmUUID)
	log.Println(dom)
	if err != nil {
		log.Println(err)
		c.JSON(http.StatusInternalServerError, common.Error{Error: err.Error()})
		return
	}

	stat, _, err := dom.GetState()
	if err != nil {
		log.Println(err)
		c.JSON(http.StatusInternalServerError, common.Error{Error: err.Error()})
		return
	}

	if stat != libvirt.DOMAIN_RUNNING {
		if err = dom.Create(); err != nil {
			c.JSON(http.StatusInternalServerError, common.Error{Error: err.Error()})
		}
	}

	err = dom.Free()
	if err != nil {
		log.Println(err)
	}

	c.JSON(http.StatusOK, common.Result{})

}

func ShutdownAdmin(c *gin.Context) {
	var input vm.VirtualMachineStop

	err := c.BindJSON(&input)
	if err != nil {
		c.JSON(http.StatusBadRequest, common.Error{Error: err.Error()})
		return
	}

	vmUUID := c.Param("vm_uuid")
	nodeID, err := strconv.Atoi(c.Param("node_id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, common.Error{Error: err.Error()})
		return
	}

	resultNode := dbNode.Get(node.ID, &core.Node{Model: gorm.Model{ID: uint(nodeID)}})
	if resultNode.Err != nil {
		c.JSON(http.StatusBadRequest, common.Error{Error: resultNode.Err.Error()})
		return
	}

	log.Println("qemu+ssh://" + resultNode.Node[0].UserName + "@" + resultNode.Node[0].IP + "/system")
	conn, err := libvirt.NewConnect("qemu+ssh://" + resultNode.Node[0].UserName + "@" + resultNode.Node[0].IP + "/system")
	if err != nil {
		log.Println("failed to connect to qemu: " + err.Error())
		c.JSON(http.StatusInternalServerError, common.Error{Error: err.Error()})
		return
	}

	defer conn.Close()

	dom, err := conn.LookupDomainByUUIDString(vmUUID)
	log.Println(dom)
	if err != nil {
		log.Println(err)
		c.JSON(http.StatusInternalServerError, common.Error{Error: err.Error()})
		return
	}

	stat, _, err := dom.GetState()
	if err != nil {
		log.Println(err)
		c.JSON(http.StatusInternalServerError, common.Error{Error: err.Error()})
		return
	}

	if stat != libvirt.DOMAIN_SHUTOFF {
		// Forceがtrueである場合、強制終了
		if input.Force {
			if err = dom.Destroy(); err != nil {
				c.JSON(http.StatusInternalServerError, common.Error{Error: err.Error()})
			}
		} else {
			if err = dom.Shutdown(); err != nil {
				c.JSON(http.StatusInternalServerError, common.Error{Error: err.Error()})
			}
		}
	}

	err = dom.Free()
	if err != nil {
		log.Println(err)
	}

	c.JSON(http.StatusOK, common.Result{})
}

func ResetAdmin(c *gin.Context) {
	vmUUID := c.Param("vm_uuid")
	nodeID, err := strconv.Atoi(c.Param("node_id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, common.Error{Error: err.Error()})
		return
	}

	resultNode := dbNode.Get(node.ID, &core.Node{Model: gorm.Model{ID: uint(nodeID)}})
	if resultNode.Err != nil {
		c.JSON(http.StatusBadRequest, common.Error{Error: resultNode.Err.Error()})
		return
	}

	log.Println("qemu+ssh://" + resultNode.Node[0].UserName + "@" + resultNode.Node[0].IP + "/system")
	conn, err := libvirt.NewConnect("qemu+ssh://" + resultNode.Node[0].UserName + "@" + resultNode.Node[0].IP + "/system")
	if err != nil {
		log.Println("failed to connect to qemu: " + err.Error())
		c.JSON(http.StatusInternalServerError, common.Error{Error: err.Error()})
		return
	}

	defer conn.Close()

	dom, err := conn.LookupDomainByUUIDString(vmUUID)
	log.Println(dom)
	if err != nil {
		log.Println(err)
		c.JSON(http.StatusInternalServerError, common.Error{Error: err.Error()})
		return
	}

	if err = dom.Reset(0); err != nil {
		c.JSON(http.StatusInternalServerError, common.Error{Error: err.Error()})
	}

	err = dom.Free()
	if err != nil {
		log.Println(err)
	}

	c.JSON(http.StatusOK, common.Result{})
}

func SuspendAdmin(c *gin.Context) {
	vmUUID := c.Param("vm_uuid")
	nodeID, err := strconv.Atoi(c.Param("node_id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, common.Error{Error: err.Error()})
		return
	}

	resultNode := dbNode.Get(node.ID, &core.Node{Model: gorm.Model{ID: uint(nodeID)}})
	if resultNode.Err != nil {
		c.JSON(http.StatusBadRequest, common.Error{Error: resultNode.Err.Error()})
		return
	}

	log.Println("qemu+ssh://" + resultNode.Node[0].UserName + "@" + resultNode.Node[0].IP + "/system")
	conn, err := libvirt.NewConnect("qemu+ssh://" + resultNode.Node[0].UserName + "@" + resultNode.Node[0].IP + "/system")
	if err != nil {
		log.Println("failed to connect to qemu: " + err.Error())
		c.JSON(http.StatusInternalServerError, common.Error{Error: err.Error()})
		return
	}

	defer conn.Close()

	dom, err := conn.LookupDomainByUUIDString(vmUUID)
	log.Println(dom)
	if err != nil {
		log.Println(err)
		c.JSON(http.StatusInternalServerError, common.Error{Error: err.Error()})
		return
	}

	if err = dom.Suspend(); err != nil {
		c.JSON(http.StatusInternalServerError, common.Error{Error: err.Error()})
	}

	err = dom.Free()
	if err != nil {
		log.Println(err)
	}

	c.JSON(http.StatusOK, common.Result{})
}

func ResumeAdmin(c *gin.Context) {
	vmUUID := c.Param("vm_uuid")
	nodeID, err := strconv.Atoi(c.Param("node_id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, common.Error{Error: err.Error()})
		return
	}

	resultNode := dbNode.Get(node.ID, &core.Node{Model: gorm.Model{ID: uint(nodeID)}})
	if resultNode.Err != nil {
		c.JSON(http.StatusBadRequest, common.Error{Error: resultNode.Err.Error()})
		return
	}

	log.Println("qemu+ssh://" + resultNode.Node[0].UserName + "@" + resultNode.Node[0].IP + "/system")
	conn, err := libvirt.NewConnect("qemu+ssh://" + resultNode.Node[0].UserName + "@" + resultNode.Node[0].IP + "/system")
	if err != nil {
		log.Println("failed to connect to qemu: " + err.Error())
		c.JSON(http.StatusInternalServerError, common.Error{Error: err.Error()})
		return
	}

	defer conn.Close()

	dom, err := conn.LookupDomainByUUIDString(vmUUID)
	log.Println(dom)
	if err != nil {
		log.Println(err)
		c.JSON(http.StatusInternalServerError, common.Error{Error: err.Error()})
		return
	}

	if err = dom.Resume(); err != nil {
		c.JSON(http.StatusInternalServerError, common.Error{Error: err.Error()})
	}

	err = dom.Free()
	if err != nil {
		log.Println(err)
	}

	c.JSON(http.StatusOK, common.Result{})
}
