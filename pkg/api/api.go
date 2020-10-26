package api

import (
	"github.com/gin-gonic/gin"
	controller "github.com/vmmgr/controller/pkg/api/core/controller/v0"
	group "github.com/vmmgr/controller/pkg/api/core/group/v0"
	notice "github.com/vmmgr/controller/pkg/api/core/notice/v0"
	ticket "github.com/vmmgr/controller/pkg/api/core/support/ticket/v0"
	token "github.com/vmmgr/controller/pkg/api/core/token/v0"
	"github.com/vmmgr/controller/pkg/api/core/tool/config"
	user "github.com/vmmgr/controller/pkg/api/core/user/v0"
	"log"
	"net/http"
	"strconv"
)

func ControllerRestAPI() {
	router := gin.Default()
	router.Use(cors)

	go token.TokenRemove()

	api := router.Group("/api")
	{
		v1 := api.Group("/v1")
		{
			// Controller
			//
			v1.POST("/controller/chat", controller.ReceiveChatAdmin)

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

			v1.POST("/token", token.AddAdmin)
			// Token Delete
			v1.DELETE("/token", token.DeleteAllAdmin)
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
		}
	}
	ws := router.Group("/ws")
	{
		v1 := ws.Group("/v1")
		{
			v1.GET("/support", ticket.GetAdminWebSocket)
		}
	}

	go ticket.HandleMessagesAdmin()
	log.Fatal(http.ListenAndServe(":"+strconv.Itoa(config.Conf.Controller.Admin.Port), router))
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
