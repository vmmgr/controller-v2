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
		var msg vm.WebSocketAdminInput
		err = conn.ReadJSON(&msg)
		if err != nil {
			log.Printf("error: %v", err)
			delete(vm.Clients, &vm.WebSocket{UUID: uuid, Admin: true, GroupID: 0, Socket: conn})
			break
		}

		if msg.Type == 0 {
			// Get
			log.Println("WebSocket VM Get " + msg.UUID)
			_, conn, err := connectLibvirt(msg.NodeID)
			if err != nil {
				log.Println(err)
				continue
			}

			dom, err := conn.LookupDomainByUUIDString(msg.UUID)
			if err != nil {
				log.Println(err)
				continue
			}

			t := libVirtXml.Domain{}
			stat, _, _ := dom.GetState()
			xmlString, _ := dom.GetXMLDesc(libvirt.DOMAIN_XML_SECURE)
			xml.Unmarshal([]byte(xmlString), &t)

			vm.ClientBroadcast <- vm.WebSocketResult{
				UUID:      uuid,
				Type:      0,
				CreatedAt: time.Now(),
				Status:    true,
				Code:      0,
				VMDetail:  []vm.Detail{{VM: t, Stat: uint(stat)}},
			}

		} else if msg.Type == 1 {
			// Get All
			log.Println("WebSocket VM GetAll")
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
				continue
			}

			var vms []vm.Detail

			for _, tmpNode := range resultNode.Node {
				log.Printf("[%s] %s\n", tmpNode.IP, tmpNode.User)
				conn, err := libvirt.NewConnect("qemu+ssh://" + tmpNode.User + "@" + tmpNode.IP + "/system")
				if err != nil {
					log.Println("failed to connect to qemu: " + err.Error())
					//vm.ClientBroadcast <- vm.WebSocketResult{
					//	UUID:      uuid,
					//	Type:      1,
					//	Err:       err.Error(),
					//	CreatedAt: time.Now(),
					//	Status:    false,
					//	Code:      0,
					//}
					continue
				}
				defer conn.Close()

				net, _ := conn.ListNetworks()
				log.Println(net)
				doms, err := conn.ListAllDomains(libvirt.CONNECT_LIST_DOMAINS_ACTIVE | libvirt.CONNECT_LIST_DOMAINS_INACTIVE)
				if err != nil {
					log.Printf("ListAllDomains error: %s", err)
					//vm.ClientBroadcast <- vm.WebSocketResult{
					//	UUID:      uuid,
					//	Type:      1,
					//	Err:       err.Error(),
					//	CreatedAt: time.Now(),
					//	Status:    false,
					//	Code:      0,
					//}
					continue
				} else {
					for _, dom := range doms {
						t := libVirtXml.Domain{}
						stat, _, _ := dom.GetState()
						xmlString, _ := dom.GetXMLDesc(libvirt.DOMAIN_XML_SECURE)
						xml.Unmarshal([]byte(xmlString), &t)

						vms = append(vms, vm.Detail{
							Node: tmpNode.ID,
							VM:   t,
							Stat: uint(stat),
						})
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
				continue
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
				continue
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
					continue
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
					continue
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
					continue
				}
			}
		} else if msg.Type == 11 {
			// Delete
		} else if msg.Type == 20 {
			// Start
			detail, err := Startup(msg.NodeID, msg.UUID)
			if err != nil {
				log.Println(err)
				continue
			}
			vm.ClientBroadcast <- vm.WebSocketResult{
				UUID:      msg.UUID,
				Type:      20,
				Err:       "",
				CreatedAt: time.Now(),
				Status:    true,
				VMDetail:  []vm.Detail{*detail},
			}
		} else if msg.Type == 21 {
			// Force Shutdown
			detail, err := Shutdown(msg.NodeID, msg.UUID, true)
			if err != nil {
				log.Println(err)
				continue
			}
			vm.ClientBroadcast <- vm.WebSocketResult{
				UUID:      msg.UUID,
				Type:      21,
				Err:       "",
				CreatedAt: time.Now(),
				Status:    true,
				VMDetail:  []vm.Detail{*detail},
			}
		} else if msg.Type == 22 {
			// Shutdown
			detail, err := Shutdown(msg.NodeID, msg.UUID, false)
			if err != nil {
				log.Println(err)
				continue
			}
			vm.ClientBroadcast <- vm.WebSocketResult{
				UUID:      msg.UUID,
				Type:      22,
				Err:       "",
				CreatedAt: time.Now(),
				Status:    true,
				VMDetail:  []vm.Detail{*detail},
			}
		} else if msg.Type == 23 {
			// Reset
			detail, err := Reset(msg.NodeID, msg.UUID)
			if err != nil {
				log.Println(err)
				continue
			}
			vm.ClientBroadcast <- vm.WebSocketResult{
				UUID:      msg.UUID,
				Type:      23,
				Err:       "",
				CreatedAt: time.Now(),
				Status:    true,
				VMDetail:  []vm.Detail{*detail},
			}
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
			log.Println("WebSocket VMID: " + strconv.Itoa(int(msg.ID)) + " Get")

			var vmData *core.VM = nil
			for _, tmpVM := range result.User.Group.VMs {
				if tmpVM.ID == msg.ID {
					vmData = tmpVM
					continue
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
				continue
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
				continue
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
				continue
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
				continue
			}
		} else if msg.Type == 1 {
			// Get All
			log.Println("WebSocket VM Get")

			var vms []vm.Detail
			var nodes []core.Node
			for _, tmpVM := range result.User.Group.VMs {
				find := false
				for _, tmpNode := range nodes {
					if tmpNode.ID == tmpVM.NodeID {
						find = true
						continue
					}
				}
				if !find {
					nodes = append(nodes, tmpVM.Node)
				}
			}

			for _, tmpNode := range nodes {
				log.Printf("[%s] %s\n", tmpNode.IP, tmpNode.User)
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
								Node: tmpNode.ID,
								VM:   t,
								Stat: uint(stat),
							})
							continue
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
		} else if msg.Type == 10 {
			for {
				// Create
				resultStorage := dbStorage.Get(storage.ID, &core.Storage{Model: gorm.Model{ID: msg.Create.Template.StorageID}})
				if resultStorage.Err != nil {
					log.Println(resultStorage.Err)
					continue
				}

				resultVM := dbVM.Get(vm.GroupID, &core.VM{GroupID: &result.Group.ID})
				if resultVM.Err != nil {
					log.Printf("error vm: %s\n", resultVM.Err)
					continue
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
					continue
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
					continue
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
					continue
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
					continue
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
					continue
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
					groupID:  *result.User.GroupID,
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
					continue
				}
				continue
			}
		} else if msg.Type == 15 {
			// Delete(Unde)
		} else if msg.Type == 16 {
			// Delete(Force)
		} else if msg.Type == 20 {
			// Start
			var vmData *core.VM = nil
			for _, tmpVM := range result.User.Group.VMs {
				if tmpVM.ID == msg.ID {
					vmData = tmpVM
					continue
				}
			}
			if vmData == nil {
				log.Printf("VM ID mismatch: %s", err)
				vm.ClientBroadcast <- vm.WebSocketResult{
					ID:        msg.ID,
					Type:      20,
					Err:       err.Error(),
					CreatedAt: time.Now(),
					Status:    false,
					Code:      0,
				}
				continue
			}

			detail, err := Startup(vmData.Node.ID, vmData.UUID)
			if err != nil {
				log.Println(err)
				continue
			}

			vm.ClientBroadcast <- vm.WebSocketResult{
				ID:        msg.ID,
				Type:      20,
				Err:       "",
				CreatedAt: time.Now(),
				Status:    true,
				VMDetail:  []vm.Detail{*detail},
			}
		} else if msg.Type == 21 {
			// Force Shutdown
			var vmData *core.VM = nil
			for _, tmpVM := range result.User.Group.VMs {
				if tmpVM.ID == msg.ID {
					vmData = tmpVM
					continue
				}
			}
			if vmData == nil {
				log.Printf("VM ID mismatch: %s", err)
				vm.ClientBroadcast <- vm.WebSocketResult{
					ID:        msg.ID,
					Type:      21,
					Err:       err.Error(),
					CreatedAt: time.Now(),
					Status:    false,
					Code:      0,
				}
				continue
			}

			detail, err := Shutdown(vmData.Node.ID, vmData.UUID, true)
			if err != nil {
				log.Println(err)
				continue
			}
			vm.ClientBroadcast <- vm.WebSocketResult{
				ID:        msg.ID,
				Type:      21,
				Err:       "",
				CreatedAt: time.Now(),
				Status:    true,
				VMDetail:  []vm.Detail{*detail},
			}
		} else if msg.Type == 22 {
			// Shutdown
			var vmData *core.VM = nil
			for _, tmpVM := range result.User.Group.VMs {
				if tmpVM.ID == msg.ID {
					vmData = tmpVM
					continue
				}
			}
			if vmData == nil {
				log.Printf("VM ID mismatch: %s", err)
				vm.ClientBroadcast <- vm.WebSocketResult{
					ID:        msg.ID,
					Type:      22,
					Err:       err.Error(),
					CreatedAt: time.Now(),
					Status:    false,
					Code:      0,
				}
				continue
			}

			detail, err := Shutdown(vmData.Node.ID, vmData.UUID, false)
			if err != nil {
				log.Println(err)
				continue
			}
			vm.ClientBroadcast <- vm.WebSocketResult{
				ID:        msg.ID,
				Type:      22,
				Err:       "",
				CreatedAt: time.Now(),
				Status:    true,
				VMDetail:  []vm.Detail{*detail},
			}
		} else if msg.Type == 23 {
			// Reset
			var vmData *core.VM = nil
			for _, tmpVM := range result.User.Group.VMs {
				if tmpVM.ID == msg.ID {
					vmData = tmpVM
					continue
				}
			}
			if vmData == nil {
				log.Printf("VM ID mismatch: %s", err)
				vm.ClientBroadcast <- vm.WebSocketResult{
					ID:        msg.ID,
					Type:      23,
					Err:       err.Error(),
					CreatedAt: time.Now(),
					Status:    false,
					Code:      0,
				}
				continue
			}

			detail, err := Reset(vmData.Node.ID, vmData.UUID)
			if err != nil {
				log.Println(err)
				continue
			}
			vm.ClientBroadcast <- vm.WebSocketResult{
				ID:        msg.ID,
				Type:      23,
				Err:       "",
				CreatedAt: time.Now(),
				Status:    true,
				VMDetail:  []vm.Detail{*detail},
			}
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
