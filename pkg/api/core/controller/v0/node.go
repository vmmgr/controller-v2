package v0

import (
	"github.com/gin-gonic/gin"
	"github.com/vmmgr/controller/pkg/api/core/controller"
	"github.com/vmmgr/controller/pkg/api/core/vm"
	"log"
	"time"
)

func ReceiveNode(c *gin.Context) {
	var input controller.Node
	c.BindJSON(&input)

	log.Println(input)

	vm.Broadcast <- vm.WebSocketResult{
		CreatedAt: time.Now(),
		Progress:  input.Progress,
		GroupID:   input.GroupID,
		UUID:      input.UUID,
		FilePath:  input.FilePath,
	}
}
