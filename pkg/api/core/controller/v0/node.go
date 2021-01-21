package v0

import (
	"github.com/gin-gonic/gin"
	"github.com/vmmgr/controller/pkg/api/core/controller"
	request "github.com/vmmgr/controller/pkg/api/core/request/v0"
	"github.com/vmmgr/controller/pkg/api/core/vm"
	"log"
	"time"
)

func ReceiveNode(c *gin.Context) {
	var input controller.Node
	err := c.BindJSON(&input)
	if err != nil {
		log.Println(err)
	}
	log.Println(input)

	vm.ClientBroadcast <- vm.WebSocketResult{
		CreatedAt: time.Now(),
		Progress:  input.Progress,
		GroupID:   input.GroupID,
		Status:    input.Status,
		UUID:      input.UUID,
		FilePath:  input.FilePath,
		Message:   input.Comment,
	}

	if input.Code == 2 && input.Progress == 100 && input.Status {
		request.Delete(input.UUID)
	}

	if !input.Status {
		log.Println(input.Comment)
		request.Delete(input.UUID)
	}
}
