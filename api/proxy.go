package api

import (
	"context"
	"fmt"
	"github.com/koding/websocketproxy"
	"github.com/labstack/echo"
	"github.com/vmmgr/controller/data"
	"github.com/vmmgr/controller/db"
	pb "github.com/vmmgr/node/proto/proto-go"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
	"log"
	"net/http"
	"net/url"
	"strconv"
	"time"
)

const vncWebSocketPort = 7000

func webSocketProxy() func(echo.Context) error {
	return func(c echo.Context) error {
		token := c.Param("token")
		id, _ := strconv.Atoi(c.Param("id"))
		group, _ := strconv.Atoi(c.Param("group"))

		for _, dataNode := range db.GetAllDBNode() {
			if dataNode.ID == uint(id/100000) {
				conn, err := grpc.Dial(dataNode.IP, grpc.WithInsecure(), grpc.WithBlock(), grpc.WithTimeout(2*time.Second))
				if err != nil {
					log.Fatalf("Not connect; %v", err)
				}
				defer conn.Close()

				client := pb.NewNodeClient(conn)
				header := metadata.New(map[string]string{"node": "true"})
				ctx := metadata.NewOutgoingContext(context.Background(), header)

				r, err := client.GetVM(ctx, &pb.VMData{ID: int64(id % 100000)})
				if err != nil {
					log.Fatal(err)
				}

				if 2 < data.VerifySameGroup(token, group) || int(r.GroupID) != group {
					return nil
				}

				u := &url.URL{
					Scheme: "ws",
					Host:   fmt.Sprintf("%s:%d", dataNode.IP, (id%100000)+vncWebSocketPort),
					Path:   "/",
				}
				log.Println(u)

				ws := &websocketproxy.WebsocketProxy{
					Backend: func(r *http.Request) *url.URL {
						return u
					},
				}

				delete(c.Request().Header, "Origin")
				log.Printf("[DEBUG] websocket proxy requesting to backend '%s'\n", ws.Backend(c.Request()))
				ws.ServeHTTP(c.Response(), c.Request())
			}
		}
		return nil
	}
}
