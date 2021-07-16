package v2

import (
	"encoding/xml"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"github.com/libvirt/libvirt-go"
	libVirtXml "github.com/libvirt/libvirt-go-xml"
	"github.com/vmmgr/controller/pkg/api/core"
	auth "github.com/vmmgr/controller/pkg/api/core/auth/v0"
	"github.com/vmmgr/controller/pkg/api/core/tool/gen"
	"github.com/vmmgr/controller/pkg/api/core/vm"
	"github.com/vmmgr/controller/pkg/api/core/vm/storage"
	"github.com/vmmgr/controller/pkg/api/store/ip"
	dbIP "github.com/vmmgr/controller/pkg/api/store/ip/v0"
	dbStorage "github.com/vmmgr/controller/pkg/api/store/node/storage/v0"
	dbNode "github.com/vmmgr/controller/pkg/api/store/node/v0"
	dbVM "github.com/vmmgr/controller/pkg/api/store/vm/v0"
	"gorm.io/gorm"
	"log"
	"strconv"
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
		conn.WriteMessage(websocket.ClosePolicyViolation, []byte("error: auth error"))
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
			for {
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
					break
				}

				var vms []vm.Detail

				for _, tmpNode := range resultNode.Node {
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

					if err == nil {
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
							break
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
				break
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
					groupID:  0,
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
		} else if msg.Type == 20 {
			// Start

		} else if msg.Type == 21 {
			// Shutdown

		} else if msg.Type == 22 {
			// Reset

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
		conn.WriteMessage(websocket.ClosePolicyViolation, []byte("error: auth error"))
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
			for {
				log.Println("WebSocket VM " + strconv.Itoa(int(msg.ID)) + " Get")

				var vmData *core.VM = nil
				for _, tmpVM := range result.User.Group.VMs {
					if tmpVM.ID == msg.ID {
						vmData = tmpVM
						break
					}
				}
				if vmData == nil {
					log.Printf("VM ID mismatch: %s", err)
					vm.ClientBroadcast <- vm.WebSocketResult{
						UUID:      uuid,
						Type:      0,
						Err:       err.Error(),
						CreatedAt: time.Now(),
						Status:    false,
						Code:      0,
					}
					break
				}

				conn, err := libvirt.NewConnect("qemu+ssh://" + vmData.Node.User + "@" + vmData.Node.IP + "/system")
				if err != nil {
					log.Println("failed to connect to qemu: " + err.Error())
					vm.ClientBroadcast <- vm.WebSocketResult{
						UUID:      uuid,
						Type:      0,
						Err:       err.Error(),
						CreatedAt: time.Now(),
						Status:    false,
						Code:      0,
					}
					break
				}
				defer conn.Close()

				doms, err := conn.ListAllDomains(libvirt.CONNECT_LIST_DOMAINS_ACTIVE | libvirt.CONNECT_LIST_DOMAINS_INACTIVE)
				if err != nil {
					log.Printf("ListAllDomains error: %s", err)
					vm.ClientBroadcast <- vm.WebSocketResult{
						UUID:      uuid,
						Type:      0,
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

					vmUUID, _ := dom.GetUUIDString()
					if vmData.UUID == vmUUID {
						vm.ClientBroadcast <- vm.WebSocketResult{
							UUID:      uuid,
							Type:      0,
							CreatedAt: time.Now(),
							Status:    true,
							Code:      0,
							VMDetail:  []vm.Detail{{VM: t, Stat: uint(stat)}},
						}
					}
				}
				break
			}
		} else if msg.Type == 1 {
			// Get All
			for {
				log.Println("WebSocket VM Get")

				var vms []vm.Detail
				var nodes []core.Node
				for _, tmpVM := range result.User.Group.VMs {
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

						for _, tmpVM := range result.User.Group.VMs {
							vmUUID, _ := dom.GetUUIDString()
							if tmpVM.UUID == vmUUID {
								vms = append(vms, vm.Detail{
									ID:   tmpVM.ID,
									VM:   t,
									Stat: uint(stat),
								})
							}
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
				break
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

				resultIP, err := dbIP.Get(ip.GetUnused, &core.IP{})
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

				if len(resultIP) == 0 {
					log.Println("ip data is not found...")
					vm.ClientBroadcast <- vm.WebSocketResult{
						UUID:      uuid,
						Type:      10,
						Err:       "ip data is not found...",
						CreatedAt: time.Now(),
						Status:    false,
						Code:      0,
					}
					break
				}

				err = dbIP.Update(ip.UpdateReserved, core.IP{
					Model:    gorm.Model{ID: resultIP[0].ID},
					Reserved: &[]bool{true}[0],
				})
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

				msg.Create.Template.IP = resultIP[0].IP
				msg.Create.Template.NetMask = resultIP[0].Subnet
				msg.Create.Template.Gateway = resultIP[0].Gateway
				msg.Create.Template.DNS = resultIP[0].DNS

				//----ベースイメージコピー処理----
				h := NewVMUserTemplateHandler(VMTemplateHandler{
					uuid:     uuid,
					input:    msg.Create.VM,
					template: msg.Create.Template,
					node:     *node,
					storage:  resultStorage.Storage[0],
					conn:     conn,
					groupID:  result.User.GroupID,
					ipID:     resultIP[0].ID,
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
					break
				}

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
