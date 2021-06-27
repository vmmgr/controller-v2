package v0

import (
	"bytes"
	"encoding/json"
	"github.com/gin-gonic/gin"
	auth "github.com/vmmgr/controller/pkg/api/core/auth/v0"
	"github.com/vmmgr/controller/pkg/api/core/controller"
	"github.com/vmmgr/controller/pkg/api/core/support"
	"github.com/vmmgr/controller/pkg/api/core/tool/config"
	"github.com/vmmgr/controller/pkg/api/core/tool/hash"
	vmInt "github.com/vmmgr/controller/pkg/api/core/vm"
	vm "github.com/vmmgr/controller/pkg/api/core/vm/v2"
	imaConController "github.com/vmmgr/imacon/pkg/api/core/controller"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
	"time"
)

func SendChatAdmin(data controller.Chat) {
	client := &http.Client{}
	client.Timeout = time.Second * 5

	body, _ := json.Marshal(controller.Chat{
		Err:       data.Err,
		CreatedAt: data.CreatedAt,
		UserID:    data.UserID,
		UserName:  data.UserName,
		GroupID:   data.GroupID,
		Admin:     data.Admin,
		Message:   data.Message,
	})

	//Header部分
	header := http.Header{}
	header.Set("Content-Length", "10000")
	header.Add("Content-Type", "application/json")
	header.Add("TOKEN_1", config.Conf.Controller.Auth.Token1)
	header.Add("TOKEN_2", hash.Generate(config.Conf.Controller.Auth.Token2+config.Conf.Controller.Auth.Token3))

	//リクエストの作成
	req, err := http.NewRequest("POST", "http://"+config.Conf.Controller.User.IP+":"+
		strconv.Itoa(config.Conf.Controller.User.Port)+"/api/v1/controller/chat", bytes.NewBuffer(body))
	if err != nil {
		return
	}

	req.Header = header

	resp, err := client.Do(req)
	if err != nil {
		log.Println(err)
		return
	}
	defer resp.Body.Close()

	_, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Println(err)
		return
	}
}

func ReceiveChatAdmin(c *gin.Context) {
	token1 := c.Request.Header.Get("TOKEN_1")
	token2 := c.Request.Header.Get("TOKEN_2")

	if err := auth.ControllerAuthentication(controller.Controller{Token1: token1, Token2: token2}); err != nil {
		log.Println(err)
	}

	var input controller.Chat
	log.Println(c.BindJSON(&input))

	support.Broadcast <- support.WebSocketResult{
		CreatedAt: time.Now(),
		UserID:    input.UserID,
		GroupID:   input.GroupID,
		Admin:     input.Admin,
		Message:   input.Message,
	}
}

func ReceiveImaConAdmin(c *gin.Context) {
	//token1 := c.Request.Header.Get("TOKEN_1")
	//token2 := c.Request.Header.Get("TOKEN_2")
	//
	//if err := auth.ControllerAuthentication(controller.Controller{Token1: token1, Token2: token2}); err != nil {
	//	log.Println(err)
	//}

	var input imaConController.Controller
	log.Println(c.BindJSON(&input))

	controller.ImaConBroadcast <- controller.WebSocketImaConResult{
		CreatedAt: time.Now(),
		UUID:      input.UUID,
		Progress:  input.Progress,
		Finish:    input.Finish,
	}
}

func HandleMessages() {
	for {
		msg := <-controller.ImaConBroadcast
		log.Println(msg)

		//登録されているクライアント宛にデータ送信する
		//コントローラが管理者モードの場合
		if msg.Finish {
			log.Println(msg.Progress)
			vmInt.ClientBroadcast <- vmInt.WebSocketResult{
				UUID:       msg.UUID,
				Type:       10,
				Status:     true,
				Processing: true,
				CreatedAt:  time.Now(),
				Message:    "[finish] storage copy",
				Progress:   100,
			}
			vmh := vm.NewVMHandler(vm.LoadCopy[msg.UUID])
			err := vmh.CreateVM()
			if err != nil {
				log.Println(err)
				vmInt.ClientBroadcast <- vmInt.WebSocketResult{
					UUID:      msg.UUID,
					Type:      10,
					CreatedAt: time.Now(),
					Status:    false,
					Err:       err.Error(),
					Message:   "[error] storage copy",
					Progress:  100,
				}
			} else {
				log.Println("No Error")
				vmInt.ClientBroadcast <- vmInt.WebSocketResult{
					UUID:       msg.UUID,
					Type:       10,
					Status:     true,
					Processing: true,
					CreatedAt:  time.Now(),
					Message:    "[finish] VM Create",
					Progress:   100,
					Finish:     true,
				}
			}

			delete(vm.LoadCopy, msg.UUID)
		} else {
			vmInt.ClientBroadcast <- vmInt.WebSocketResult{
				UUID:       msg.UUID,
				Type:       10,
				Status:     true,
				Processing: true,
				CreatedAt:  time.Now(),
				Message:    "[copy] storage copy",
				Progress:   msg.Progress,
			}
		}
	}
}
