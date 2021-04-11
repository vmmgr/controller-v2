package v0

import (
	"github.com/gin-gonic/gin"
	"github.com/pkg/sftp"
	"github.com/vmmgr/controller/pkg/api/core"
	"github.com/vmmgr/controller/pkg/api/core/node/storage/vmImage"
	"github.com/vmmgr/controller/pkg/api/core/tool/gen"
	"github.com/vmmgr/controller/pkg/api/core/tool/ssh"
	dbNode "github.com/vmmgr/controller/pkg/api/store/node/v0"
	"log"
	"regexp"
	"strings"
	"time"
)

func GetVMListWebSocketAdmin(c *gin.Context) {
	conn, err := vmImage.WsUpgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		log.Printf("Failed to set websocket upgrade: %+v\n", err)
		return
	}

	defer conn.Close()

	resultNode := dbNode.GetAll()
	if resultNode.Err != nil {
		log.Println(resultNode.Err)
		delete(vmImage.ListClients, &vmImage.WebSocketList{Admin: true, GroupID: 0, Socket: conn, Error: resultNode.Err})
		return
	}

	uuid := gen.GenerateUUID()

	// WebSocket送信
	vmImage.ListClients[&vmImage.WebSocketList{Admin: true, GroupID: 0, Socket: conn}] = true

	for _, tmpNode := range resultNode.Node {
		go GetWebSocketAdminVM(tmpNode, uuid)
	}

	//WebSocket受信
	for {
		var msg vmImage.WebSocketListResult
		err = conn.ReadJSON(&msg)
		if err != nil {
			log.Printf("error: %v", err)
			delete(vmImage.ListClients, &vmImage.WebSocketList{Admin: true, GroupID: 0, Socket: conn})
			break
		}
	}
}

func VMListHandleMessages(admin bool) {
	for {
		msg := <-vmImage.ListClientBroadcast

		//登録されているクライアント宛にデータ送信する
		//コントローラが管理者モードの場合
		for client := range vmImage.ListClients {
			if admin {
				err := client.Socket.WriteJSON(msg)
				if err != nil {
					log.Printf("error: %v", err)
					client.Socket.Close()
					delete(vmImage.ListClients, client)
				}
			}
		}
	}
}

func GetWebSocketAdminVM(node core.Node, uuid string) {

	conn, err := ssh.ConnectSSH(node)
	if err != nil {
		vmImage.ListClientBroadcast <- vmImage.WebSocketListResult{
			NodeID: node.ID,
			Err:    err.Error(),
		}
		return
	}

	defer conn.Close()

	log.Println(node.Storage)

	// SFTP Client
	client, err := sftp.NewClient(conn)
	if err != nil {
		vmImage.ListClientBroadcast <- vmImage.WebSocketListResult{
			NodeID: node.ID,
			Err:    err.Error(),
		}
		return
	}
	defer client.Close()
	for _, tmp := range node.Storage {
		log.Println(tmp)
		if *tmp.VMImage {
			log.Println("BasePath: " + tmp.Path + "/template")
			w := client.Walk(tmp.Path + "/template")
			for w.Step() {
				if w.Err() != nil {
					continue
				}
				match, _ := regexp.MatchString("^"+tmp.Path+"/template"+"$", w.Path())
				if !match {
					vmImage.ListClientBroadcast <- vmImage.WebSocketListResult{
						NodeID:      node.ID,
						Name:        w.Path()[len(tmp.Path)+len("/template")+1:],
						Err:         "",
						CreatedAt:   time.Time{},
						UserToken:   "",
						AccessToken: "",
						Size:        w.Stat().Size(),
						Time:        w.Stat().ModTime().String(),
						FilePath:    w.Path(),
						CloudInit:   strings.Contains(w.Path(), vmImage.CloudInitString),
						Admin:       false,
						Message:     "",
					}
				}
			}
		}
	}
}
