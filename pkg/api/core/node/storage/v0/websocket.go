package v0

import (
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"github.com/vmmgr/controller/pkg/api/core"
	auth "github.com/vmmgr/controller/pkg/api/core/auth/v0"
	"github.com/vmmgr/controller/pkg/api/core/node/storage"
	"github.com/vmmgr/controller/pkg/api/core/tool/gen"
	dbNode "github.com/vmmgr/controller/pkg/api/store/node/v0"
	"log"
)

func GetWebSocketProgressAdmin(c *gin.Context) {
	conn, err := storage.WsUpgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		log.Printf("Failed to set websocket upgrade: %+v\n", err)
		return
	}

	defer conn.Close()

	resultNode := dbNode.GetAll()
	if resultNode.Err != nil {
		log.Println(resultNode.Err)
		delete(storage.Clients, &storage.WebSocket{Admin: true, GroupID: 0, Socket: conn, Error: resultNode.Err})
		return
	}

	uuid := gen.GenerateUUID()

	// WebSocket送信
	storage.Clients[&storage.WebSocket{Admin: true, UUID: uuid, GroupID: 0, Socket: conn}] = true

	//WebSocket受信
	for {
		var msg storage.WebSocketResult
		err = conn.ReadJSON(&msg)
		if err != nil {
			log.Printf("error: %v", err)
			delete(storage.Clients, &storage.WebSocket{Admin: true, GroupID: 0, Socket: conn})
			break
		}
	}
}

func GetWebSocketProgress(c *gin.Context) {

	// /vm?user_token=accessID?access_token=token
	//  user_token = UserToken, access_token = AccessToken

	userToken := c.Query("user_token")
	accessToken := c.Query("access_token")

	conn, err := storage.WsUpgrader.Upgrade(c.Writer, c.Request, nil)
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
	storage.Clients[&storage.WebSocket{Admin: false, UserID: result.User.ID, GroupID: result.Group.ID, Socket: conn}] = true
}

func HandleMessagesProgress(admin bool) {
	for {
		msg := <-storage.ClientBroadcast

		//登録されているクライアント宛にデータ送信する
		//コントローラが管理者モードの場合
		for client := range storage.Clients {
			if admin {
				err := client.Socket.WriteJSON(msg)
				if err != nil {
					log.Printf("error: %v", err)
					client.Socket.Close()
					delete(storage.Clients, client)
				}
			} else {
				//コントローラがユーザモードの場合

				//Pathを空で上書きする
				msg.FilePath = ""

				//if client.GroupID == msg.GroupID {
				//	err := client.Socket.WriteJSON(msg)
				//	if err != nil {
				//		log.Printf("error: %v", err)
				//		client.Socket.Close()
				//		delete(storage.Clients, client)
				//	}
				//}
			}
		}
	}
}
