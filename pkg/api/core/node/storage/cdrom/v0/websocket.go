package v0

import (
	"github.com/gin-gonic/gin"
	"github.com/pkg/sftp"
	"github.com/vmmgr/controller/pkg/api/core"
	cdrom "github.com/vmmgr/controller/pkg/api/core/node/storage/cdrom"
	"github.com/vmmgr/controller/pkg/api/core/tool/gen"
	"github.com/vmmgr/controller/pkg/api/core/tool/ssh"
	dbNode "github.com/vmmgr/controller/pkg/api/store/node/v0"
	"log"
	"regexp"
	"time"
)

func GetCDROMListWebSocketAdmin(c *gin.Context) {
	conn, err := cdrom.WsUpgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		log.Printf("Failed to set websocket upgrade: %+v\n", err)
		return
	}

	defer conn.Close()

	resultNode := dbNode.GetAll()
	if resultNode.Err != nil {
		log.Println(resultNode.Err)
		delete(cdrom.ListClients, &cdrom.WebSocketList{Admin: true, GroupID: 0, Socket: conn, Error: resultNode.Err})
		return
	}

	uuid := gen.GenerateUUID()

	// WebSocket送信
	cdrom.ListClients[&cdrom.WebSocketList{Admin: true, GroupID: 0, Socket: conn}] = true

	for _, tmpNode := range resultNode.Node {
		go GetWebSocketAdminVM(tmpNode, uuid)
	}

	//WebSocket受信
	for {
		var msg cdrom.WebSocketListResult
		err = conn.ReadJSON(&msg)
		if err != nil {
			log.Printf("error: %v", err)
			delete(cdrom.ListClients, &cdrom.WebSocketList{Admin: true, GroupID: 0, Socket: conn})
			break
		}
	}
}

func CDROMListHandleMessages(admin bool) {
	for {
		msg := <-cdrom.ListClientBroadcast

		//登録されているクライアント宛にデータ送信する
		//コントローラが管理者モードの場合
		for client := range cdrom.ListClients {
			if admin {
				err := client.Socket.WriteJSON(msg)
				if err != nil {
					log.Printf("error: %v", err)
					client.Socket.Close()
					delete(cdrom.ListClients, client)
				}
			}
		}
	}
}

func GetWebSocketAdminVM(node core.Node, uuid string) {

	conn, err := ssh.ConnectSSH(node)
	if err != nil {
		cdrom.ListClientBroadcast <- cdrom.WebSocketListResult{
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
		cdrom.ListClientBroadcast <- cdrom.WebSocketListResult{
			NodeID: node.ID,
			Err:    err.Error(),
		}
		return
	}
	defer client.Close()
	for _, tmp := range node.Storage {
		log.Println(tmp)
		if !*tmp.VMImage {
			log.Println("BasePath: " + tmp.Path)
			w := client.Walk(tmp.Path)
			for w.Step() {
				if w.Err() != nil {
					continue
				}
				match, _ := regexp.MatchString("^"+tmp.Path+"$", w.Path())
				if !match {
					cdrom.ListClientBroadcast <- cdrom.WebSocketListResult{
						NodeID:      node.ID,
						Name:        w.Path()[len(tmp.Path)+1:],
						Err:         "",
						CreatedAt:   time.Time{},
						UserToken:   "",
						AccessToken: "",
						Size:        w.Stat().Size(),
						Time:        w.Stat().ModTime().String(),
						FilePath:    w.Path(),
						Admin:       false,
						Message:     "",
					}
				}
			}
		}
	}
}
