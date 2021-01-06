package v0

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	"github.com/vmmgr/controller/pkg/api/core/node"
	"github.com/vmmgr/controller/pkg/api/core/tool/client"
	"github.com/vmmgr/controller/pkg/api/core/vm"
	dbNode "github.com/vmmgr/controller/pkg/api/store/node/v0"
	"log"
	"net/http"
	"strconv"
)

func StartupAdmin(c *gin.Context) {
	vmID := c.Param("uuid")
	nodeID, err := strconv.Atoi(c.Param("node_id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, vm.ResultAdmin{Status: http.StatusBadRequest, Error: err.Error()})
		return
	}
	nodeResult := dbNode.Get(node.ID, &node.Node{Model: gorm.Model{ID: uint(nodeID)}})
	if nodeResult.Err != nil {
		c.JSON(http.StatusBadRequest, vm.ResultAdmin{Status: http.StatusBadRequest, Error: err.Error()})
		return
	}

	response, err := client.Put("http://"+nodeResult.Node[0].IP+":"+
		strconv.Itoa(int(nodeResult.Node[0].Port))+"/api/v1/vm/"+vmID+"/power", "")
	if err != nil {
		log.Println(err)
		c.JSON(http.StatusInternalServerError, vm.ResultAdmin{Status: http.StatusInternalServerError, Error: err.Error()})
		return
	}

	log.Println(response)

	c.JSON(http.StatusOK, vm.ResultAdmin{Status: http.StatusOK})

}

func ShutdownAdmin(c *gin.Context) {
	vmID := c.Param("uuid")
	nodeID, err := strconv.Atoi(c.Param("node_id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, vm.ResultAdmin{Status: http.StatusBadRequest, Error: err.Error()})
		return
	}

	body, err := c.GetRawData()
	if err != nil {
		fmt.Println("JSON marshal error: ", err)
		return
	}

	nodeResult := dbNode.Get(node.ID, &node.Node{Model: gorm.Model{ID: uint(nodeID)}})
	if nodeResult.Err != nil {
		c.JSON(http.StatusBadRequest, vm.ResultAdmin{Status: http.StatusBadRequest, Error: err.Error()})
		return
	}

	response, err := client.Delete("http://"+nodeResult.Node[0].IP+":"+
		strconv.Itoa(int(nodeResult.Node[0].Port))+"/api/v1/vm/"+vmID+"/power", string(body))
	if err != nil {
		log.Println(err)
		c.JSON(http.StatusInternalServerError, vm.ResultAdmin{Status: http.StatusInternalServerError, Error: err.Error()})
		return
	}

	log.Println(response)
	c.JSON(http.StatusOK, vm.ResultAdmin{Status: http.StatusOK})

}

func ResetAdmin(c *gin.Context) {
	vmID := c.Param("uuid")
	nodeID, err := strconv.Atoi(c.Param("node_id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, vm.ResultAdmin{Status: http.StatusBadRequest, Error: err.Error()})
		return
	}

	nodeResult := dbNode.Get(node.ID, &node.Node{Model: gorm.Model{ID: uint(nodeID)}})
	if nodeResult.Err != nil {
		c.JSON(http.StatusBadRequest, vm.ResultAdmin{Status: http.StatusBadRequest, Error: err.Error()})
		return
	}

	response, err := client.Put("http://"+nodeResult.Node[0].IP+":"+
		strconv.Itoa(int(nodeResult.Node[0].Port))+"/api/v1/vm/"+vmID+"/reset", "")
	if err != nil {
		log.Println(err)
		c.JSON(http.StatusInternalServerError, vm.ResultAdmin{Status: http.StatusInternalServerError, Error: err.Error()})
		return
	}

	log.Println(response)

	c.JSON(http.StatusOK, vm.ResultAdmin{Status: http.StatusOK})
}
