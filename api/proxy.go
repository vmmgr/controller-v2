package api

import (
	"fmt"
	"github.com/koding/websocketproxy"
	"github.com/labstack/echo"
	"github.com/vmmgr/controller/data"
	"github.com/vmmgr/controller/db"
	"log"
	"net/http"
	"net/url"
)

const vncwebsocketport = 7000

func WebSocketProxy() func(echo.Context) error {
	return func(c echo.Context) error {
		uuid := c.Param("uuid")
		vmname := c.Param("vmname")

		group, result := db.GetDBGroupToken(uuid)
		if result == false {
			log.Printf("Get Error: GroupDBGroupToken")
			return nil
		}
		nodeip, port := data.GetVNCVMData(group.ID, vmname)
		if nodeip == "" && port == 0 {
			log.Printf("Get Error: VNC Port")
			return nil
		}

		u := &url.URL{
			Scheme: "ws",
			Host:   fmt.Sprintf("%s:%d", nodeip, port+vncwebsocketport),
			Path:   "/",
		}
		fmt.Println(u)

		ws := &websocketproxy.WebsocketProxy{
			Backend: func(r *http.Request) *url.URL {
				return u
			},
		}
		delete(c.Request().Header, "Origin")
		log.Printf("[DEBUG] websocket proxy requesting to backend '%s'", ws.Backend(c.Request()))
		ws.ServeHTTP(c.Response(), c.Request())

		return nil
	}
}
