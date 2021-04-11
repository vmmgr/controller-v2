package nic

import (
	"github.com/gorilla/websocket"
	"github.com/vmmgr/controller/pkg/api/core"
	"net/http"
	"time"
)

const (
	ID        = 0
	NodeID    = 1
	GroupID   = 2
	AdminOnly = 3
	Name      = 4
	Enable    = 5
	Virtual   = 6
	UpdateAll = 110
)

// channel定義(websocketで使用)
var ListClients = make(map[*WebSocketList]bool)
var ListBroadcast = make(chan WebSocketListResult)
var ListClientBroadcast = make(chan WebSocketListResult)

// websocket用
type WebSocketListResult struct {
	NodeID      uint      `json:"node_id"`
	Name        string    `json:"name"`
	FilePath    string    `json:"file_path"`
	Size        int64     `json:"size"`
	Time        string    `json:"time"`
	Err         string    `json:"error"`
	CreatedAt   time.Time `json:"created_at"`
	UserToken   string    `json:"user_token"`
	AccessToken string    `json:"access_token"`
	Admin       bool      `json:"admin"`
	Message     string    `json:"message"`
}

type WebSocketList struct {
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

type Post struct {
	Name string `json:"name"`
	URL  string `json:"url"`
}

type Delete struct {
	Name string `json:"name"`
}

type Result struct {
	Status bool       `json:"status"`
	Error  string     `json:"error"`
	NIC    []core.NIC `json:"nic"`
}

type ResultOne struct {
	Status bool     `json:"status"`
	Error  string   `json:"error"`
	NIC    core.NIC `json:"nic"`
}

type ResultDatabase struct {
	Err error
	NIC []core.NIC
}
