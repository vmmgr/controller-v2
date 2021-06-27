package v0

import (
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	"github.com/koding/websocketproxy"
	"github.com/vmmgr/controller/pkg/api/core"
	"github.com/vmmgr/controller/pkg/api/core/node"
	"github.com/vmmgr/controller/pkg/api/core/token"
	dbNode "github.com/vmmgr/controller/pkg/api/store/node/v0"
	dbToken "github.com/vmmgr/controller/pkg/api/store/token/v0"
	"log"
	"net/http"
	"net/url"
	"strconv"
)

//
// comment
//?host=" + ip + "&port=" + port + "&path=api/" + group.UUID + "/" + r.GetVmname() + "/vnc
//http://localhost/noVNC/vnc.html?host=127.0.0.1&port=8081&path=ws/v1/vnc
// ws/v1/vnc/[user_token]/[access_token]/[node]?uuid=[uuid]

func Get(c *gin.Context) {
	//tokenData := c.Param("request")
	userToken := c.Param("user_token")
	accessToken := c.Param("access_token")
	nodeID, _ := strconv.Atoi(c.Param("node"))
	//vmUUID := c.Query("uuid")

	var tokenResult token.ResultDatabase

	if userToken == "0" {
		tokenResult = dbToken.Get(token.AccessToken, &core.Token{AccessToken: accessToken})
		if tokenResult.Err != nil {
			log.Println(tokenResult.Err)
			return
		}
	} else {
		tokenResult = dbToken.Get(token.UserTokenAndAccessToken, &core.Token{UserToken: userToken, AccessToken: accessToken})
		if tokenResult.Err != nil {
			log.Println(tokenResult.Err)
			return
		}
	}

	log.Println(tokenResult)

	//管理者ではない場合
	if !*tokenResult.Token[0].Admin {
		//管理者ではない場合はVM Tableより検証を行う必要あり
	}

	nodeResult := dbNode.Get(node.ID, &core.Node{Model: gorm.Model{ID: uint(nodeID)}})
	if nodeResult.Err != nil {
		log.Println(nodeResult.Err)
		return
	}

	//res, err := vm.Get(nodeResult.Node[0].IP, nodeResult.Node[0].Port, vmUUID)
	//if err != nil {
	//	log.Println(err)
	//	return
	//}

	//if len(res.Data.VM.VM.Devices.Graphics) < 1 {
	//	log.Println("Error: No Graphics...")
	//	return
	//}

	//webSocketPort := res.Data.VM.VM.Devices.Graphics[0].VNC.WebSocket

	u := &url.URL{
		Scheme: "ws",
		//Host:   fmt.Sprintf("%s:%d", nodeResult.Node[0].IP, webSocketPort),
		Path: "/",
	}
	log.Println(u)

	ws := &websocketproxy.WebsocketProxy{
		Backend: func(r *http.Request) *url.URL {
			return u
		},
	}

	delete(c.Request.Header, "Origin")
	log.Printf("[DEBUG] websocket proxy requesting to backend '%s'\n", ws.Backend(c.Request))
	ws.ServeHTTP(c.Writer, c.Request)

	//return func(c echo.Context) error {
	//	token := c.Param("request")
	//	id, _ := strconv.Atoi(c.Param("uuid"))
	//	group, _ := strconv.Atoi(c.Param("group"))
	//
	//	for _, dataNode := range db.GetAllDBNode() {
	//		if dataNode.ID == uint(id/100000) {
	//			conn, err := grpc.Dial(dataNode.IP, grpc.WithInsecure(), grpc.WithBlock(), grpc.WithTimeout(2*time.Second))
	//			if err != nil {
	//				log.Fatalf("Not connect; %v", err)
	//			}
	//			defer conn.Close()
	//
	//			client := pb.NewNodeClient(conn)
	//			header := metadata.New(map[string]string{"node": "true"})
	//			ctx := metadata.NewOutgoingContext(context.Background(), header)
	//
	//			r, err := client.GetVM(ctx, &pb.VMData{ID: int64(id % 100000)})
	//			if err != nil {
	//				log.Fatal(err)
	//			}
	//
	//			if 2 < data.VerifySameGroup(token, group) || int(r.GroupID) != group {
	//				return nil
	//			}
	//
	//			u := &url.URL{
	//				Scheme: "ws",
	//				Host:   fmt.Sprintf("%s:%d", dataNode.IP, (id%100000)+vncWebSocketPort),
	//				Path:   "/",
	//			}
	//			log.Println(u)
	//
	//			ws := &websocketproxy.WebsocketProxy{
	//				Backend: func(r *http.Request) *url.URL {
	//					return u
	//				},
	//			}
	//
	//			delete(c.Request().Header, "Origin")
	//			log.Printf("[DEBUG] websocket proxy requesting to backend '%s'\n", ws.Backend(c.Request()))
	//			ws.ServeHTTP(c.Response(), c.Request())
	//		}
	//	}
	//	return nil
	//}
}
