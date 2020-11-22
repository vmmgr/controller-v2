package vm

import (
	"github.com/gorilla/websocket"
	"github.com/jinzhu/gorm"
	"github.com/vmmgr/controller/pkg/api/core/vm/nic"
	"github.com/vmmgr/controller/pkg/api/core/vm/storage"
	"net/http"
	"time"
)

const (
	ID             = 0
	NodeID         = 1
	GroupID        = 2
	UUID           = 3
	NodeAndVNCPort = 4
	Lock           = 5
	UpdateName     = 100
	UpdateNode     = 101
	UpdateGroup    = 102
	UpdateBoot     = 103
	UpdateInfo     = 104
	UpdateAll      = 110
)

// channel定義(websocketで使用)
var Clients = make(map[*WebSocket]bool)
var Broadcast = make(chan WebSocketResult)

// websocket用
type WebSocketResult struct {
	ID          uint      `json:"id"`
	Err         string    `json:"error"`
	CreatedAt   time.Time `json:"created_at"`
	UserToken   string    `json:"user_token"`
	AccessToken string    `json:"access_token"`
	UUID        string    `json:"uuid"`
	GroupID     uint      `json:"group_id"`
	FilePath    string    `json:"file_path"`
	Admin       bool      `json:"admin"`
	Message     string    `json:"message"`
	Progress    uint      `json:"progress"`
}

type WebSocket struct {
	TicketID uint
	GroupID  uint
	UserID   uint
	Admin    bool
	Socket   *websocket.Conn
}

type VM struct {
	gorm.Model
	NodeID  uint   `json:"node_id"`
	GroupID uint   `json:"group_id"`
	Name    string `json:"name"`
	UUID    string `json:"uuid"`
	VNCPort uint   `json:"vnc_port"`
	Lock    *bool  `json:"lock"`
}

type Input struct {
	NodeID   uint   `json:"node_id"`
	GroupID  uint   `json:"group_id"`
	Name     string `json:"name"`
	UUID     string `json:"uuid"`
	OS       uint   `json:"os"` //32bit=> 32 64bit=> 64
	CPU      uint   `json:"cpu"`
	CPUModel string `json:"cpu_mode"`
	Memory   uint   `json:"memory"`
	VNCPort  uint   `json:"vnc_port"`
	Boot     uint   `json:"boot"` //0: hd 1:cdrom 2:floppy
	Lock     *bool  `json:"lock"`
}

type CreateAdmin struct {
	VM            VM              `json:"vm"`
	Storage       storage.Storage `json:"storage"`
	NIC           nic.NIC         `json:"nic"`
	TemplateApply bool            `json:"template_apply"`
	Template      Template        `json:"template"`
}

type DeleteAdmin struct {
	DeleteStorage bool `json:"delete_storage"`
}

type UserInput struct {
}

type GetImaCon struct {
	Status int `json:"status"`
	Data   struct {
		Path string `json:"path"`
	} `json:"data"`
}

type VMAll struct {
	VM      VM
	Storage storage.Storage
	NIC     nic.NIC
}

type Template struct {
	Name            string `json:"name"`
	NodeID          uint   `json:"node_id"`
	TemplateID      uint   `json:"template_id"`
	TemplatePlanID  uint   `json:"template_plan_id"`
	StorageCapacity uint   `json:"storage_capacity"`
	StoragePathType uint   `json:"storage_path_type"`
	NICType         uint   `json:"nic_type"` //0:default 1~:custom
}

type Result struct {
	Status bool   `json:"status"`
	Error  string `json:"error"`
	VM     []VM   `json:"vm"`
}

type ResultOne struct {
	Status bool   `json:"status"`
	Error  string `json:"error"`
	VM     VM     `json:"vm"`
}

type ResultDatabase struct {
	Err error
	VMs []VM
}

var WsUpgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}
