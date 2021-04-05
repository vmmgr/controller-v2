package v0

import (
	"github.com/gin-gonic/gin"
	websocket "github.com/gorilla/websocket"
	"github.com/vmmgr/controller/pkg/api/core"
	auth "github.com/vmmgr/controller/pkg/api/core/auth/v0"
	"github.com/vmmgr/controller/pkg/api/core/vm"
	"log"
)

func GetWebSocketAdmin(c *gin.Context) {
	conn, err := vm.WsUpgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		log.Printf("Failed to set websocket upgrade: %+v\n", err)
		return
	}

	defer conn.Close()

	// WebSocket送信
	vm.Clients[&vm.WebSocket{Admin: true, GroupID: 0, Socket: conn}] = true

	//WebSocket受信
	for {
		var msg vm.WebSocketResult
		err = conn.ReadJSON(&msg)
		if err != nil {
			log.Printf("error: %v", err)
			delete(vm.Clients, &vm.WebSocket{Admin: true, GroupID: 0, Socket: conn})
			break
		}
	}
}

func GetWebSocket(c *gin.Context) {

	// /vm?user_token=accessID?access_token=token
	//  user_token = UserToken, access_token = AccessToken

	userToken := c.Query("user_token")
	accessToken := c.Query("access_token")

	conn, err := vm.WsUpgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		log.Printf("Failed to set websocket upgrade: %+v\n", err)
		return
	}

	defer conn.Close()

	result := auth.GroupAuthentication(0, core.Token{UserToken: userToken, AccessToken: accessToken})
	if result.Err != nil {
		log.Println("ws:// support error:Auth error")
		conn.WriteMessage(websocket.TextMessage, []byte("error: auth error"))
		return
	}

	// WebSocket送信
	vm.Clients[&vm.WebSocket{Admin: false, UserID: result.User.ID, GroupID: result.Group.ID, Socket: conn}] = true
}

func HandleMessages(admin bool) {
	for {
		msg := <-vm.ClientBroadcast

		//登録されているクライアント宛にデータ送信する
		//コントローラが管理者モードの場合
		for client := range vm.Clients {
			if admin {
				err := client.Socket.WriteJSON(msg)
				if err != nil {
					log.Printf("error: %v", err)
					client.Socket.Close()
					delete(vm.Clients, client)
				}
			} else {
				//コントローラがユーザモードの場合

				//Pathを空で上書きする
				msg.FilePath = ""

				if client.GroupID == msg.GroupID {
					err := client.Socket.WriteJSON(msg)
					if err != nil {
						log.Printf("error: %v", err)
						client.Socket.Close()
						delete(vm.Clients, client)
					}
				}
			}
		}
	}
}
