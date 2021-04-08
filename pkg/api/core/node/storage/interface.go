package storage

import (
	"github.com/gorilla/websocket"
	"github.com/vmmgr/controller/pkg/api/core"
	"net/http"
	"time"
)

const (
	ID        = 0
	NodeID    = 1
	AdminOnly = 2
	Name      = 3
	UpdateAll = 110
)

type Result struct {
	Storage []core.Storage `json:"storage"`
}

type ResultOne struct {
	Status  bool         `json:"status"`
	Error   string       `json:"error"`
	Storage core.Storage `json:"storage"`
}

type ResultDatabase struct {
	Err     error
	Storage []core.Storage
}

// channel定義(websocketで使用)
var Clients = make(map[*WebSocket]bool)
var Broadcast = make(chan WebSocketResult)
var ClientBroadcast = make(chan WebSocketResult)

// websocket用
type WebSocketResult struct {
	NodeID    uint      `json:"node_id"`
	Name      string    `json:"name"`
	Err       string    `json:"error"`
	CreatedAt time.Time `json:"created_at"`
	Status    int       `json:"status"`
	Code      uint      `json:"code"`
	FilePath  string    `json:"file_path"`
	Admin     bool      `json:"admin"`
	Message   string    `json:"message"`
	Progress  uint      `json:"progress"`
	UUID      string    `json:"uuid"`
}

type WebSocket struct {
	UUID    string
	GroupID uint
	UserID  uint
	Admin   bool
	Error   error
	Socket  *websocket.Conn
}

var WsUpgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}
