package v2

import (
	"encoding/xml"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"github.com/jinzhu/gorm"
	"github.com/libvirt/libvirt-go"
	libVirtXml "github.com/libvirt/libvirt-go-xml"
	"github.com/vmmgr/controller/pkg/api/core"
	auth "github.com/vmmgr/controller/pkg/api/core/auth/v0"
	"github.com/vmmgr/controller/pkg/api/core/group"
	"github.com/vmmgr/controller/pkg/api/core/tool/gen"
	"github.com/vmmgr/controller/pkg/api/core/vm"
	"github.com/vmmgr/controller/pkg/api/core/vm/storage"
	dbGroup "github.com/vmmgr/controller/pkg/api/store/group/v0"
	"github.com/vmmgr/controller/pkg/api/store/ip"
	dbIP "github.com/vmmgr/controller/pkg/api/store/ip/v0"
	dbStorage "github.com/vmmgr/controller/pkg/api/store/node/storage/v0"
	dbNode "github.com/vmmgr/controller/pkg/api/store/node/v0"
	dbVM "github.com/vmmgr/controller/pkg/api/store/vm/v0"
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

	accessToken := c.Query("access_token")
	result := auth.AdminAuthentication(accessToken)
	if result.Err != nil {
		log.Println("ws:// support error:Auth error")
		conn.WriteMessage(websocket.TextMessage, []byte("error: auth error"))
		return
	}

	uuid := gen.GenerateUUID()

	// WebSocket送信
	vm.Clients[&vm.WebSocket{
		UUID:    uuid,
		Admin:   true,
		GroupID: 0,
		Socket:  conn,
	}] = true

	//WebSocket受信
	for {
		var msg vm.WebSocketInput
		err = conn.ReadJSON(&msg)
		if err != nil {
			log.Printf("error: %v", err)
			delete(vm.Clients, &vm.WebSocket{UUID: uuid, Admin: true, GroupID: 0, Socket: conn})
			break
		}

		if msg.Type == 0 {
			// Get
		} else if msg.Type == 1 {
			// Get All
			log.Println("WebSocket VM Get")
			resultNode := dbNode.GetAll()
			if resultNode.Err != nil {
				log.Println(resultNode.Err)
				vm.ClientBroadcast <- vm.WebSocketResult{
					UUID:      uuid,
					Type:      1,
					Err:       resultNode.Err.Error(),
					CreatedAt: time.Now(),
					Status:    false,
					Code:      0,
				}
			}

			var vms []vm.Detail

			for _, tmpNode := range resultNode.Node {
				for {

					conn, err := libvirt.NewConnect("qemu+ssh://" + tmpNode.User + "@" + tmpNode.IP + "/system")
					if err != nil {
						log.Println("failed to connect to qemu: " + err.Error())
						vm.ClientBroadcast <- vm.WebSocketResult{
							UUID:      uuid,
							Type:      1,
							Err:       err.Error(),
							CreatedAt: time.Now(),
							Status:    false,
							Code:      0,
						}
						return
					}
					defer conn.Close()

					doms, err := conn.ListAllDomains(libvirt.CONNECT_LIST_DOMAINS_ACTIVE | libvirt.CONNECT_LIST_DOMAINS_INACTIVE)
					if err != nil {
						log.Printf("ListAllDomains error: %s", err)
						vm.ClientBroadcast <- vm.WebSocketResult{
							UUID:      uuid,
							Type:      1,
							Err:       err.Error(),
							CreatedAt: time.Now(),
							Status:    false,
							Code:      0,
						}
						return
					}

					for _, dom := range doms {
						t := libVirtXml.Domain{}
						stat, _, _ := dom.GetState()
						xmlString, _ := dom.GetXMLDesc(libvirt.DOMAIN_XML_SECURE)
						xml.Unmarshal([]byte(xmlString), &t)

						checkSame := false
						for _, vm := range vms {
							vmUUID, _ := dom.GetUUIDString()
							if vm.VM.UUID == vmUUID {
								checkSame = true
							}
						}

						if !checkSame {
							vms = append(vms, vm.Detail{
								VM:   t,
								Stat: uint(stat),
							})
						}
					}

					vm.ClientBroadcast <- vm.WebSocketResult{
						UUID:      uuid,
						Type:      1,
						CreatedAt: time.Now(),
						Status:    true,
						Code:      0,
						VMDetail:  vms,
					}
				}
			}
		} else if msg.Type == 10 {
			// Create

			resultStorage := dbStorage.Get(storage.ID, &core.Storage{Model: gorm.Model{ID: msg.Create.Template.StorageID}})
			if resultStorage.Err != nil {
				log.Println(resultStorage.Err)
			}

			// nodeIDが存在するか確認
			node, conn, err := connectLibvirt(resultStorage.Storage[0].NodeID)
			if err != nil {
				log.Println(err)
				vm.ClientBroadcast <- vm.WebSocketResult{
					UUID:      uuid,
					Type:      10,
					Err:       err.Error(),
					CreatedAt: time.Now(),
					Status:    false,
					Code:      0,
				}
				return
			}

			if !msg.Create.TemplateApply {
				//手動作成時
				//VM作成用のデータ
				h := NewVMHandler(VMHandler{
					Conn: conn,
					VM:   msg.Create.VM,
					Node: *node,
				})

				err = h.CreateVM()
				if err != nil {
					log.Println(err)
					vm.ClientBroadcast <- vm.WebSocketResult{
						UUID:      uuid,
						Type:      10,
						Err:       err.Error(),
						CreatedAt: time.Now(),
						Status:    false,
						Code:      0,
					}
					//return
				}
			} else {
				log.Println("Template Apply")
				// storage
				vmBasePath := resultStorage.Storage[0]
				if vmBasePath.ID == 0 {
					log.Println("vmBasePath ID === 0")
					vm.ClientBroadcast <- vm.WebSocketResult{
						UUID:      uuid,
						Type:      10,
						Err:       "ID BasePath invalid...",
						CreatedAt: time.Now(),
						Status:    false,
						Code:      0,
					}
					//return
				}

				node.Storage = []core.Storage{vmBasePath}

				//----ベースイメージコピー処理----
				h := NewVMAdminTemplateHandler(VMTemplateHandler{
					uuid:     uuid,
					input:    msg.Create.VM,
					template: msg.Create.Template,
					node:     *node,
					storage:  resultStorage.Storage[0],
					conn:     conn,
				})

				log.Println("start template apply")

				err = h.templateApply()
				if err != nil {
					log.Println(err)
					vm.ClientBroadcast <- vm.WebSocketResult{
						UUID:      uuid,
						Type:      10,
						Err:       err.Error(),
						CreatedAt: time.Now(),
						Status:    false,
						Code:      0,
					}
					//return
				}
			}

		} else if msg.Type == 11 {
			// Delete
		}
	}
}

func GetWebSocket(c *gin.Context) {
	conn, err := vm.WsUpgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		log.Printf("Failed to set websocket upgrade: %+v\n", err)
		return
	}

	defer conn.Close()

	// /vm?user_token=accessID?access_token=token
	//  user_token = UserToken, access_token = AccessToken

	userToken := c.Query("user_token")
	accessToken := c.Query("access_token")

	uuid := gen.GenerateUUID()

	result := auth.GroupAuthentication(0, core.Token{UserToken: userToken, AccessToken: accessToken})
	if result.Err != nil {
		log.Printf("ws:// Auth error: %v\n", result.Err)
		conn.WriteMessage(websocket.TextMessage, []byte("error: auth error"))
		return
	}

	// WebSocket送信
	vm.Clients[&vm.WebSocket{
		UUID:    uuid,
		Admin:   false,
		UserID:  result.User.ID,
		GroupID: result.Group.ID,
		Socket:  conn,
	}] = true

	//WebSocket受信
	for {
		var msg vm.WebSocketInput
		err = conn.ReadJSON(&msg)
		if err != nil {
			log.Printf("error: %v", err)
			delete(vm.Clients, &vm.WebSocket{UUID: uuid, Admin: true, GroupID: 0, Socket: conn})
			break
		}

		log.Println(msg)

		if msg.Type == 0 {
			// Get
		} else if msg.Type == 1 {
			// Get All
			log.Println("WebSocket VM Get")
			resultGroup := dbGroup.Get(group.ID, &core.Group{Model: gorm.Model{ID: result.User.GroupID}})
			if resultGroup.Err != nil {
				log.Println(resultGroup.Err)
				vm.ClientBroadcast <- vm.WebSocketResult{
					UUID:      uuid,
					Type:      1,
					Err:       resultGroup.Err.Error(),
					CreatedAt: time.Now(),
					Status:    false,
					Code:      0,
				}
			}

			var vms []vm.Detail
			var nodes []core.Node
			for _, tmpVM := range resultGroup.Group[0].VMs {
				find := false
				for _, tmpNode := range nodes {
					if tmpNode.ID == tmpVM.NodeID {
						find = true
						break
					}
				}
				if !find {
					nodes = append(nodes, tmpVM.Node)
				}
			}

			for _, tmpNode := range nodes {
				conn, err := libvirt.NewConnect("qemu+ssh://" + tmpNode.User + "@" + tmpNode.IP + "/system")
				if err != nil {
					log.Println("failed to connect to qemu: " + err.Error())
					vm.ClientBroadcast <- vm.WebSocketResult{
						UUID:      uuid,
						Type:      1,
						Err:       err.Error(),
						CreatedAt: time.Now(),
						Status:    false,
						Code:      0,
					}
				}
				defer conn.Close()

				doms, err := conn.ListAllDomains(libvirt.CONNECT_LIST_DOMAINS_ACTIVE | libvirt.CONNECT_LIST_DOMAINS_INACTIVE)
				if err != nil {
					log.Printf("ListAllDomains error: %s", err)
					vm.ClientBroadcast <- vm.WebSocketResult{
						UUID:      uuid,
						Type:      1,
						Err:       err.Error(),
						CreatedAt: time.Now(),
						Status:    false,
						Code:      0,
					}
				}

				for _, dom := range doms {
					t := libVirtXml.Domain{}
					stat, _, _ := dom.GetState()
					xmlString, _ := dom.GetXMLDesc(libvirt.DOMAIN_XML_SECURE)
					xml.Unmarshal([]byte(xmlString), &t)

					checkSame := false
					groupVM := false

					for _, vm := range vms {
						vmUUID, _ := dom.GetUUIDString()
						if vm.VM.UUID == vmUUID {
							checkSame = true
						}
						for _, tmpVM := range resultGroup.Group[0].VMs {
							if tmpVM.UUID == vmUUID {
								groupVM = true
							}
						}
					}

					if !checkSame && groupVM {
						vms = append(vms, vm.Detail{
							VM:   t,
							Stat: uint(stat),
						})
					}
				}

				vm.ClientBroadcast <- vm.WebSocketResult{
					UUID:      uuid,
					Type:      1,
					CreatedAt: time.Now(),
					Status:    true,
					Code:      0,
					VMDetail:  vms,
				}
			}
		} else if msg.Type == 10 {
			for {
				// Create
				resultStorage := dbStorage.Get(storage.ID, &core.Storage{Model: gorm.Model{ID: msg.Create.Template.StorageID}})
				if resultStorage.Err != nil {
					log.Println(resultStorage.Err)
					break
				}

				resultVM := dbVM.Get(vm.GroupID, &core.VM{GroupID: result.Group.ID})
				if resultVM.Err != nil {
					log.Printf("error vm: %s\n", resultVM.Err)
					break
				}

				// nodeIDが存在するか確認
				node, conn, err := connectLibvirt(resultStorage.Storage[0].NodeID)
				if err != nil {
					log.Println(err)
					vm.ClientBroadcast <- vm.WebSocketResult{
						UUID:      uuid,
						Type:      10,
						Err:       err.Error(),
						CreatedAt: time.Now(),
						Status:    false,
						Code:      0,
					}
					break
				}

				log.Println("Template Apply")
				// storage
				vmBasePath := resultStorage.Storage[0]
				if vmBasePath.ID == 0 {
					log.Println("vmBasePath ID === 0")
					vm.ClientBroadcast <- vm.WebSocketResult{
						UUID:      uuid,
						Type:      10,
						Err:       "ID BasePath invalid...",
						CreatedAt: time.Now(),
						Status:    false,
						Code:      0,
					}
					break
				}

				node.Storage = []core.Storage{vmBasePath}

				resultIP, err := dbIP.Get(ip.GetUsed, &core.IP{})
				if err != nil {
					log.Println(err)
					vm.ClientBroadcast <- vm.WebSocketResult{
						UUID:      uuid,
						Type:      10,
						Err:       err.Error(),
						CreatedAt: time.Now(),
						Status:    false,
						Code:      0,
					}
					break
				}
				err = dbIP.Update(ip.UpdateReserved, core.IP{Model: gorm.Model{ID: resultIP[0].ID}})
				if err != nil {
					log.Println(err)
					vm.ClientBroadcast <- vm.WebSocketResult{
						UUID:      uuid,
						Type:      10,
						Err:       err.Error(),
						CreatedAt: time.Now(),
						Status:    false,
						Code:      0,
					}
					break
				}

				//----ベースイメージコピー処理----
				h := NewVMUserTemplateHandler(VMTemplateHandler{
					uuid:     uuid,
					input:    msg.Create.VM,
					template: msg.Create.Template,
					node:     *node,
					storage:  resultStorage.Storage[0],
					conn:     conn,
					groupID:  result.Group.ID,
					ipID:     resultIP[0].ID,
				})

				msg.Create.Template.IP = resultIP[0].IP
				msg.Create.Template.NetMask = resultIP[0].Subnet
				msg.Create.Template.Gateway = resultIP[0].Gateway
				msg.Create.Template.DNS = resultIP[0].DNS

				log.Println("start template apply")

				err = h.templateApply()
				if err != nil {
					log.Println(err)
					vm.ClientBroadcast <- vm.WebSocketResult{
						UUID:      uuid,
						Type:      10,
						Err:       err.Error(),
						CreatedAt: time.Now(),
						Status:    false,
						Code:      0,
					}
					break
				}

				dbVM.Create(&core.VM{
					NodeID:        node.ID,
					GroupID:       result.Group.ID,
					Name:          "",
					UUID:          "",
					VNCPort:       0,
					WebSocketPort: 0,
					Lock:          nil,
				})
				break
			}
		} else if msg.Type == 11 {
			// Delete
		}
	}
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
