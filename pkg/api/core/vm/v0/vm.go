package v0

import (
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"github.com/jinzhu/gorm"
	auth "github.com/vmmgr/controller/pkg/api/core/auth/v0"
	"github.com/vmmgr/controller/pkg/api/core/node"
	"github.com/vmmgr/controller/pkg/api/core/token"
	"github.com/vmmgr/controller/pkg/api/core/tool/client"
	"github.com/vmmgr/controller/pkg/api/core/tool/gen"
	"github.com/vmmgr/controller/pkg/api/core/vm"
	template "github.com/vmmgr/controller/pkg/api/core/vm/template/v0"
	dbNode "github.com/vmmgr/controller/pkg/api/store/node/v0"
	nodeNIC "github.com/vmmgr/node/pkg/api/core/nic"
	"github.com/vmmgr/node/pkg/api/core/storage"
	nodeVM "github.com/vmmgr/node/pkg/api/core/vm"
	"log"
	"net/http"
	"strconv"
	"time"
)

func Create(c *gin.Context) {
	var input vm.Template
	//userToken := c.Request.Header.Get("USER_TOKEN")
	//accessToken := c.Request.Header.Get("ACCESS_TOKEN")

	if err := c.BindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, vm.Result{Status: false, Error: err.Error()})
		return
	}

	//result := auth.GroupAuthentication(token.Token{UserToken: userToken, AccessToken: accessToken})
	//if result.Err != nil {
	//	c.JSON(http.StatusUnauthorized, vm.Result{Status: false, Error: result.Err.Error()})
	//	return
	//}

	// nodeIDが存在するか確認
	resultNode := dbNode.Get(node.ID, &node.Node{Model: gorm.Model{ID: input.NodeID}})
	if resultNode.Err != nil {
		c.JSON(http.StatusBadRequest, vm.Result{Status: false, Error: resultNode.Err.Error()})
		return
	}

	//nodeIDの数が0である場合は
	if len(resultNode.Node) == 0 {
		c.JSON(http.StatusNotFound, vm.Result{Status: false, Error: "node id is not found..."})
		return
	}

	////NodeのGroupIDが0かつAdminOnlyがfalseの時の以外である場合、
	//if !(resultNode.Node[0].GroupID == 0 && resultNode.Node[0].AdminOnly == &[]bool{false}[0]) {
	//	c.JSON(http.StatusForbidden, vm.Result{Status: false, Error: "You can't use this node..."})
	//	return
	//}

	// NodeIDとStoragePathTypeがGroupで使用可能か確認

	//----ベースイメージコピー時は以下に記述----
	// Templateを検索
	vmTemplate, vmTemplatePlan, err := template.GetTemplate(input.TemplateID, input.TemplatePlanID)
	if err != nil {
		c.JSON(http.StatusNotFound, vm.Result{Status: false, Error: "template is not found..."})
		return
	}

	log.Println(vmTemplate, vmTemplatePlan)

	go func() {
		imaConResult, imageResult, err := extractImaCon(vmTemplate)
		if err != nil {
			c.JSON(http.StatusNotFound, vm.Result{Status: false, Error: err.Error()})
		}

		log.Println(imaConResult)
		uuid := gen.GenerateUUID()
		path := strconv.Itoa(1) + gen.GenerateUUID() + "-1.img"
		// path := strconv.Itoa(int(result.Group.ID)) + gen.GenerateUUID() + "-1.img"
		gid := uint(0)
		// Storage作成用にbodyを作成する
		createStorageBody, _ := json.Marshal(storage.Storage{
			Mode: 1,
			FromImaCon: storage.ImaCon{
				IP:   imaConResult.IP,
				Path: imageResult.Data.Path,
			},
			Type:     10, // BootDisk(virtIO)
			FileType: 0,  // qcow2
			GroupID:  gid,
			UUID:     uuid,
			PathType: input.StoragePathType,
			Capacity: input.StorageCapacity,
			ReadOnly: false,
			Path:     path,
		})

		resultStorageProcess, err := client.Post(
			"http://"+resultNode.Node[0].IP+":"+strconv.Itoa(int(resultNode.Node[0].Port))+"/api/v1/storage",
			createStorageBody)
		if err != nil {
			log.Println(err)
			return
		}

		log.Println(resultStorageProcess)

		t := time.NewTimer(20 * time.Minute)
		defer t.Stop()

		//Todo 取りこぼす可能性があるので、要調査
	L:
		for {
			select {
			//20分以上かかる場合はタイムアウトさせる
			case <-t.C:
				log.Println("Error: timeout")
				err = fmt.Errorf("Error: timeout ")
				break L
				//UUIDとGroupIDがMatchし、Progressが100の場合、storage転送処理が終了
			case msg := <-vm.Broadcast:
				if msg.UUID == uuid && msg.GroupID == gid && msg.Progress == 100 {
					//path変数にnode側のストレージをフルパスで代入する
					path = msg.FilePath
					err = nil
					break L
				}
			}
		}

		// Errorが発生した場合
		if err != nil {
			log.Println(err)
			return
		}

		log.Println("End: copy storage")

		//VM作成用のデータ
		body, _ := json.Marshal(nodeVM.VirtualMachine{
			Name:    input.Name,
			Memory:  vmTemplatePlan.Mem,
			CPUMode: 0,
			VCPU:    vmTemplatePlan.CPU,
			NIC:     []nodeNIC.NIC{},
			VNCPort: 0, //VNCポートをNode側で自動生成
			Storage: []storage.VMStorage{
				{
					Type:     10, // BootDisk(virtIO)
					FileType: 0,  //qcow2
					Path:     path,
					ReadOnly: false,
					Boot:     0,
				},
			},
		})
		resultVMCreateProcess, err := client.Post(
			"http://"+resultNode.Node[0].IP+":"+strconv.Itoa(int(resultNode.Node[0].Port))+"/api/v1/vm",
			body)
		if err != nil {
			log.Println(err)
			return
		}
		log.Println(resultVMCreateProcess)
	}()

	c.JSON(http.StatusOK, vm.Result{Status: true})

	////DB追加
	//dbVM.Create(&vm.VM{NodeID: input.NodeID, GroupID: result.Group.ID, Name: input.Name,
	//	UUID: input.UUID, VNCPort: input.VM.VNCPort})

}

func GetWebSocket(c *gin.Context) {
	//
	// /support?uuid=0?user_token=accessID?access_token=token
	// uuid = processID, user_token = UserToken, access_token = AccessToken

	userToken := c.Query("user_token")
	accessToken := c.Query("access_token")

	id, err := strconv.Atoi(c.Query("uuid"))
	if err != nil {
		log.Println("id wrong: ", err)
		return
	}
	//wsHandle(c.Writer, c.Request)
	conn, err := vm.WsUpgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		log.Printf("Failed to set websocket upgrade: %+v\n", err)
		return
	}

	defer conn.Close()

	result := auth.GroupAuthentication(token.Token{UserToken: userToken, AccessToken: accessToken})
	if result.Err != nil {
		log.Println("ws:// support error:Auth error")
		conn.WriteMessage(websocket.TextMessage, []byte("error: auth error"))
		return
	}

	// WebSocket送信
	vm.Clients[&vm.WebSocket{TicketID: uint(id), Admin: false,
		UserID: result.User.ID, GroupID: result.Group.ID, Socket: conn}] = true
}

func HandleMessages(admin bool) {
	for {
		msg := <-vm.Broadcast

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
