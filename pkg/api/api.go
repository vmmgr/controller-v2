package api

import (
	"github.com/gin-gonic/gin"
	controller "github.com/vmmgr/controller/pkg/api/core/controller/v0"
	group "github.com/vmmgr/controller/pkg/api/core/group/v0"
	nodeNIC "github.com/vmmgr/controller/pkg/api/core/node/nic/v0"
	cdrom "github.com/vmmgr/controller/pkg/api/core/node/storage/cdrom/v0"
	nodeStorage "github.com/vmmgr/controller/pkg/api/core/node/storage/v0"
	storage "github.com/vmmgr/controller/pkg/api/core/node/storage/v0"
	vmImage "github.com/vmmgr/controller/pkg/api/core/node/storage/vmImage/v0"
	node "github.com/vmmgr/controller/pkg/api/core/node/v0"
	notice "github.com/vmmgr/controller/pkg/api/core/notice/v0"
	region "github.com/vmmgr/controller/pkg/api/core/region/v0"
	zone "github.com/vmmgr/controller/pkg/api/core/region/zone/v0"
	ticket "github.com/vmmgr/controller/pkg/api/core/support/ticket/v0"
	token "github.com/vmmgr/controller/pkg/api/core/token/v0"
	"github.com/vmmgr/controller/pkg/api/core/tool/config"
	user "github.com/vmmgr/controller/pkg/api/core/user/v0"
	template "github.com/vmmgr/controller/pkg/api/core/vm/template/v0"
	vm "github.com/vmmgr/controller/pkg/api/core/vm/v0"
	wsVNC "github.com/vmmgr/controller/pkg/api/core/wsvnc/v0"
	"log"
	"net/http"
	"strconv"
)

func AdminRestAPI() error {

	router := gin.Default()
	router.Use(cors)

	api := router.Group("/api")
	{
		v1 := api.Group("/v1")
		{
			// Controller
			//
			v1.POST("/controller/chat", controller.ReceiveChatAdmin)
			//v1.POST("/controller/node", controller.ReceiveNode)

			// Notice
			//
			v1.POST("/notice", notice.AddAdmin)
			v1.DELETE("/notice/:id", notice.DeleteAdmin)
			v1.GET("/notice", notice.GetAllAdmin)
			v1.GET("/notice/:id", notice.GetAdmin)
			v1.PUT("/notice/:id", notice.UpdateAdmin)

			//
			// User
			//
			// User Create
			v1.POST("/user", user.AddAdmin)
			// User Delete
			v1.DELETE("/user", user.DeleteAdmin)
			// User Update
			v1.PUT("/user", user.UpdateAdmin)
			v1.GET("/user", user.GetAllAdmin)
			v1.GET("/user/:id", user.GetAdmin)

			//
			// Token
			//
			v1.POST("/token/generate", token.GenerateAdmin)

			v1.POST("/token", token.GenerateAdmin)
			// Token Delete
			v1.DELETE("/token", token.Delete)
			v1.DELETE("/token/:id", token.DeleteAdmin)
			// Token Update
			v1.PUT("/token", token.UpdateAdmin)
			v1.GET("/token", token.GetAllAdmin)
			v1.GET("/token/:id", token.GetAdmin)
			//
			// Group
			//
			v1.POST("/group", group.AddAdmin)
			// Group Delete
			v1.DELETE("/group", group.DeleteAdmin)
			// Group Update
			v1.PUT("/group", group.UpdateAdmin)
			v1.GET("/group", group.GetAllAdmin)
			v1.GET("/group/:id", group.GetAdmin)

			//
			// Support
			//
			v1.POST("/support", ticket.CreateAdmin)
			v1.GET("/support", ticket.GetAllAdmin)
			//v1.POST("/support/:id", chat.AddAdmin)
			v1.GET("/support/:id", ticket.GetAdmin)
			v1.PUT("/support/:id", ticket.UpdateAdmin)

			//
			// Region
			//
			v1.POST("/region", region.AddAdmin)
			v1.GET("/region", region.GetAllAdmin)
			v1.DELETE("/region/:id", region.DeleteAdmin)
			v1.GET("/region/:id", region.GetAdmin)
			v1.PUT("/region/:id", region.UpdateAdmin)

			//
			// Zone
			//
			v1.POST("/zone/:region_id", zone.AddAdmin)
			v1.GET("/zone", zone.GetAllAdmin)
			v1.DELETE("/zone/:region_id/:zone_id", zone.DeleteAdmin)
			v1.GET("/zone/:region_id/:zone_id", zone.GetAdmin)
			v1.PUT("/zone/:region_id/:zone_id", zone.UpdateAdmin)

			//
			// Node
			//
			v1.POST("/node", node.AddAdmin)
			v1.GET("/node", node.GetAllAdmin)
			v1.DELETE("/node/:id", node.DeleteAdmin)
			v1.GET("/node/:id", node.GetAdmin)
			v1.PUT("/node/:id", node.UpdateAdmin)
			v1.GET("/node/:id/device", node.GetAllDeviceAdmin)

			//
			// Storage
			//
			//v1.POST("/storage/:node_id", nodeStorage.AddAdmin)
			v1.GET("/storage", nodeStorage.GetAllAdmin)
			//v1.DELETE("/storage/:node_id/:storage_id", nodeStorage.DeleteAdmin)
			//v1.GET("/storage/:id/:storage_id", nodeStorage.GetAdmin)
			//v1.PUT("/storage/:id/:storage_id", nodeStorage.UpdateAdmin)

			//
			// Storage Floppy
			//
			v1.POST("/storage/:id/floppy", cdrom.AddAdmin)
			v1.DELETE("/storage/:id/floppy/:image_id", cdrom.DeleteAdmin)
			v1.GET("/storage/:id/floppy", cdrom.GetAllAdmin)

			//
			// Storage CDROM
			//
			v1.POST("/storage/:id/cdrom", cdrom.AddAdmin)
			v1.DELETE("/storage/:id/cdrom/:image_id", cdrom.DeleteAdmin)
			v1.GET("/storage/:id/cdrom", cdrom.GetAllAdmin)
			//v1.GET("/node/:id", image.GetAdmin)

			//
			// Storage VM DISK
			//
			//v1.POST("/storage/:id/disk", disk.AddAdmin)
			//v1.DELETE("/storage/:id/disk/:image_id", disk.DeleteAdmin)
			//v1.GET("/storage/:id/disk", disk.GetAllAdmin)
			//v1.GET("/node/:id", image.GetAdmin)

			//
			// Storage VM
			//
			v1.POST("/storage/:id/vm", vmImage.AddAdmin)
			v1.DELETE("/storage/:id/vm/:image_id", vmImage.DeleteAdmin)
			v1.GET("/storage/:id/vm", vmImage.GetAllAdmin)
			//v1.GET("/node/:id", image.GetAdmin)

			//
			// Node NIC
			//
			v1.POST("/nic/:node_id", nodeNIC.AddAdmin)
			v1.GET("/nic/:node_id", nodeNIC.GetAllAdmin)
			v1.DELETE("/nic/:node_id/:nic_id", nodeNIC.DeleteAdmin)
			v1.GET("/nic/:node_id/:nic_id", nodeNIC.GetAdmin)
			v1.PUT("/nic/:node_id/:nic_id", nodeNIC.UpdateAdmin)

			//
			//VM
			//
			v1.POST("/vm", vm.AddAdmin)
			v1.DELETE("/vm/:node_id/:vm_uuid", vm.DeleteAdmin)
			v1.PUT("/vm/:node_id/:vm_uuid", vm.UpdateAdmin)
			v1.GET("/vm/:node_id/:vm_uuid", vm.GetAdmin)

			//
			//VM
			//
			v1.PUT("/vm/:node_id/:vm_uuid/power", vm.StartupAdmin)
			v1.DELETE("/vm/:node_id/:vm_uuid/power", vm.ShutdownAdmin)
			v1.PUT("/vm/:node_id/:vm_uuid/reset", vm.ResetAdmin)
			v1.PUT("/vm/:node_id/:vm_uuid/suspend", vm.SuspendAdmin)
			v1.PUT("/vm/:node_id/:vm_uuid/resume", vm.ResumeAdmin)

			//
			//Template
			//
			v1.GET("/template", template.Get)
		}
	}
	ws := router.Group("/ws")
	{
		v1 := ws.Group("/v1")
		{
			v1.GET("/support", ticket.GetAdminWebSocket)
			v1.GET("/vm/list", vm.GetListWebSocketAdmin)
			v1.GET("/storage/progress", storage.GetWebSocketProgressAdmin)
			v1.GET("/storage/vm/list", vmImage.GetVMListWebSocketAdmin)
			v1.GET("/storage/no_vm/list", cdrom.GetCDROMListWebSocketAdmin)
			// noVNC
			v1.GET("/vnc/:user_token/:access_token/:node", wsVNC.Get)
		}
	}

	go ticket.HandleMessagesAdmin()
	go vm.ListHandleMessages(true)
	go storage.HandleMessagesProgress(true)
	go cdrom.CDROMListHandleMessages(true)
	go vmImage.VMListHandleMessages(true)
	log.Fatal(http.ListenAndServe(":"+strconv.Itoa(config.Conf.Controller.Admin.Port), router))
	return nil
}

func UserRestAPI() {
	router := gin.Default()
	router.Use(cors)

	api := router.Group("/api")
	{
		v1 := api.Group("/v1")
		{
			// Node
			//

			// Controller
			//
			v1.POST("/controller/chat", controller.ReceiveChatUser)
			//v1.POST("/controller/node", controller.ReceiveNode)

			//
			// User
			//
			// User Delete
			//router.DELETE("/user", user.Delete)
			// User Get
			v1.GET("/user", user.Get)
			v1.GET("/user/all", user.GetGroup)
			// User ID Get
			// v1.GET("/user/:id",user.GetId)
			// User Update
			v1.PUT("/user/:id", user.Update)
			// User Mail MailVerify
			v1.GET("/user/verify/:token", user.MailVerify)
			//
			// Token
			//
			// get token for CHAP authentication
			v1.GET("/token/init", token.GenerateInit)
			// get token for user
			v1.GET("/token", token.Generate)
			// delete
			v1.DELETE("/token", token.Delete)

			//
			// Group
			//
			// Group Create
			v1.POST("/group", group.Add)
			v1.GET("/group", group.Get)
			v1.PUT("/group", group.Update)
			v1.GET("/group/all", group.GetAll)

			//
			// Support
			//
			v1.POST("/support", ticket.Create)
			//v1.GET("/support", ticket.GetTitle)
			v1.GET("/support/:id", ticket.Get)

			//
			// Notice
			//
			v1.GET("/notice", notice.Get)

			//
			// VM
			//
			v1.POST("/vm", vm.UserCreate)
			v1.DELETE("/vm/:id", vm.UserDelete)
		}
	}

	ws := router.Group("/ws")
	{
		v1 := ws.Group("/v1")
		{
			v1.GET("/support", ticket.GetWebSocket)
			v1.GET("/vm", vm.GetListWebSocket)
		}
	}

	go ticket.HandleMessages()
	go vm.ListHandleMessages(false)

	log.Fatal(http.ListenAndServe(":"+strconv.Itoa(config.Conf.Controller.User.Port), router))
}

func cors(c *gin.Context) {

	//c.Header("Access-Control-Allow-Headers", "Accept, Content-ID, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, Access-Control-Request-Headers, Access-Control-Request-Method, Connection, Host, Origin, User-Agent, Referer, Cache-Control, X-header")
	c.Header("Access-Control-Allow-Origin", "*")
	c.Header("Access-Control-Allow-Methods", "*")
	c.Header("Access-Control-Allow-Headers", "*")
	c.Header("Content-ID", "application/json")
	c.Header("Access-Control-Allow-Credentials", "true")
	//c.Header("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")

	if c.Request.Method != "OPTIONS" {
		c.Next()
	} else {
		c.AbortWithStatus(http.StatusOK)
	}
}
