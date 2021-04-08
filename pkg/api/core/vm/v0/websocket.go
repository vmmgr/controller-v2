package v0

import (
	"encoding/xml"
	"github.com/gin-gonic/gin"
	websocket "github.com/gorilla/websocket"
	"github.com/libvirt/libvirt-go"
	libVirtXml "github.com/libvirt/libvirt-go-xml"
	"github.com/vmmgr/controller/pkg/api/core"
	auth "github.com/vmmgr/controller/pkg/api/core/auth/v0"
	"github.com/vmmgr/controller/pkg/api/core/tool/gen"
	"github.com/vmmgr/controller/pkg/api/core/vm"
	dbNode "github.com/vmmgr/controller/pkg/api/store/node/v0"
	"log"
	"time"
)

func GetWebSocketAdmin(c *gin.Context) {
	conn, err := vm.WsUpgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		log.Printf("Failed to set websocket upgrade: %+v\n", err)
		return
	}

	defer conn.Close()

	resultNode := dbNode.GetAll()
	if resultNode.Err != nil {
		log.Println(resultNode.Err)
		delete(vm.Clients, &vm.WebSocket{Admin: true, GroupID: 0, Socket: conn, Error: resultNode.Err})
		return
	}

	uuid := gen.GenerateUUID()

	// WebSocket送信
	vm.Clients[&vm.WebSocket{Admin: true, UUID: uuid, GroupID: 0, Socket: conn}] = true

	for _, tmpNode := range resultNode.Node {
		go GetWebSocketAdminVM(tmpNode, uuid)
	}

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

func GetWebSocketAdminVM(node core.Node, uuid string) {

	log.Println("qemu+ssh://" + node.UserName + "@" + node.IP + "/system")
	//libvirt.NewConnectWithAuth()
	conn, err := libvirt.NewConnect("qemu+ssh://" + node.UserName + "@" + node.IP + "/system")
	if err != nil {
		log.Println("failed to connect to qemu: " + err.Error())
		vm.ClientBroadcast <- vm.WebSocketResult{
			NodeID: node.ID,
			Err:    "failed to connect to qemu: " + err.Error(),
		}
		return
	}

	defer conn.Close()

	doms, err := conn.ListAllDomains(libvirt.CONNECT_LIST_DOMAINS_ACTIVE | libvirt.CONNECT_LIST_DOMAINS_INACTIVE)
	log.Println(doms)
	if err != nil {
		log.Println(err)
		vm.ClientBroadcast <- vm.WebSocketResult{
			NodeID: node.ID,
			Err:    err.Error(),
		}
		return
		//c.JSON(http.StatusInternalServerError, common.Error{Error: err.Error()})
	}

	for _, dom := range doms {
		t := libVirtXml.Domain{}
		//stat, _, _ := dom.GetState()
		xmlString, _ := dom.GetXMLDesc(libvirt.DOMAIN_XML_SECURE)
		xml.Unmarshal([]byte(xmlString), &t)
		vm.ClientBroadcast <- vm.WebSocketResult{
			NodeID:      node.ID,
			Name:        t.Name,
			Err:         "",
			CreatedAt:   time.Time{},
			UserToken:   "",
			AccessToken: "",
			UUID:        t.UUID,
			UserUUID:    uuid,
			VCPU:        t.VCPU.Value,
			Memory:      t.Memory.Value,
			Code:        0,
			GroupID:     0,
			FilePath:    "",
			Admin:       false,
			Message:     "",
			Progress:    0,
		}
	}
}
