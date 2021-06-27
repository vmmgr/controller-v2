package v0

import (
	"encoding/xml"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	"github.com/koding/websocketproxy"
	"github.com/libvirt/libvirt-go"
	libVirtXml "github.com/libvirt/libvirt-go-xml"
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

// comment
//?host=" + ip + "&port=" + port + "&path=api/" + group.UUID + "/" + r.GetVmname() + "/vnc
//http://localhost/noVNC/vnc.html?host=127.0.0.1&port=8081&path=ws/v1/vnc/ws/v1/vnc/[access_token]/[node]/[uuid]

func GetByAdmin(c *gin.Context) {
	//tokenData := c.Param("request")
	accessToken := c.Param("access_token")
	nodeID, _ := strconv.Atoi(c.Param("node"))
	vmUUID := c.Query("uuid")

	var tokenResult token.ResultDatabase

	tokenResult = dbToken.Get(token.AccessToken, &core.Token{AccessToken: accessToken})
	if tokenResult.Err != nil {
		log.Println(tokenResult.Err)
		return
	}

	log.Println(tokenResult)

	resultNode := dbNode.Get(node.ID, &core.Node{Model: gorm.Model{ID: uint(nodeID)}})
	if resultNode.Err != nil {
		log.Println(resultNode.Err)
		return
	}

	conn, err := libvirt.NewConnect("qemu+ssh://" + resultNode.Node[0].User + "@" + resultNode.Node[0].IP + "/system")
	if err != nil {
		log.Println("failed to connect to qemu: " + err.Error())
		return
	}
	defer conn.Close()

	doms, err := conn.ListAllDomains(libvirt.CONNECT_LIST_DOMAINS_ACTIVE | libvirt.CONNECT_LIST_DOMAINS_INACTIVE)
	if err != nil {
		log.Printf("ListAllDomains error: %s", err)
		return
	}

	for _, dom := range doms {
		t := libVirtXml.Domain{}
		//stat, _, _ := dom.GetState()
		xmlString, _ := dom.GetXMLDesc(libvirt.DOMAIN_XML_SECURE)
		xml.Unmarshal([]byte(xmlString), &t)

		domVMUUID, _ := dom.GetUUIDString()
		if vmUUID == domVMUUID {
			if len(t.Devices.Graphics) < 1 {
				log.Println("Error: No Graphics...")
				return
			}

			webSocketPort := t.Devices.Graphics[0].VNC.WebSocket

			u := &url.URL{
				Scheme: "ws",
				Host:   fmt.Sprintf("%s:%d", resultNode.Node[0].IP, webSocketPort),
				Path:   "/",
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

		}
	}
}

// comment
//?host=" + ip + "&port=" + port + "&path=api/" + group.UUID + "/" + r.GetVmname() + "/vnc
//http://localhost/noVNC/vnc.html?host=127.0.0.1&port=8081&path=ws/v1/vnc/ws/v1/vnc/[user_token]/[access_token]/[vm]

func Get(c *gin.Context) {
	//tokenData := c.Param("request")
	userToken := c.Param("user_token")
	accessToken := c.Param("access_token")
	vmID, _ := strconv.Atoi(c.Param("vm"))
	//vmUUID := c.Query("uuid")

	resultToken := dbToken.Get(token.UserTokenAndAccessToken, &core.Token{UserToken: userToken, AccessToken: accessToken})
	if resultToken.Err != nil {
		log.Println(resultToken.Err)
		return
	}

	log.Println(resultToken)

	var vmData *core.VM = nil
	for _, tmpVM := range resultToken.Token[0].User.Group.VMs {
		if tmpVM.ID == uint(vmID) {
			vmData = tmpVM
			break
		}
	}
	if vmData == nil {
		log.Printf("VM ID mismatch")
		return
	}

	conn, err := libvirt.NewConnect("qemu+ssh://" + vmData.Node.User + "@" + vmData.Node.IP + "/system")
	if err != nil {
		log.Println("failed to connect to qemu: " + err.Error())
		return
	}
	defer conn.Close()

	doms, err := conn.ListAllDomains(libvirt.CONNECT_LIST_DOMAINS_ACTIVE | libvirt.CONNECT_LIST_DOMAINS_INACTIVE)
	if err != nil {
		log.Printf("ListAllDomains error: %s", err)
		return
	}

	for _, dom := range doms {
		t := libVirtXml.Domain{}
		//stat, _, _ := dom.GetState()
		xmlString, _ := dom.GetXMLDesc(libvirt.DOMAIN_XML_SECURE)
		xml.Unmarshal([]byte(xmlString), &t)

		vmUUID, _ := dom.GetUUIDString()
		if vmData.UUID == vmUUID {
			if len(t.Devices.Graphics) < 1 {
				log.Println("Error: No Graphics...")
				return
			}

			webSocketPort := t.Devices.Graphics[0].VNC.WebSocket

			u := &url.URL{
				Scheme: "ws",
				Host:   fmt.Sprintf("%s:%d", vmData.Node.IP, webSocketPort),
				Path:   "/",
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
		}
	}

}
