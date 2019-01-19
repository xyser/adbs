package handlers

import (
	"adbs/shell"
	"fmt"
	"github.com/gorilla/websocket"
	"net/http"
	"time"
)

type SockerBuff struct {
	Cmd string `json:"cmd"`
}

var wsUpgrade = websocket.Upgrader{
	ReadBufferSize:   1024,
	WriteBufferSize:  1024,
	HandshakeTimeout: 5 * time.Second,
	// 取消ws跨域校验
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

// 处理ws请求
func WsHandler(w http.ResponseWriter, r *http.Request) {
	conn, err := wsUpgrade.Upgrade(w, r, nil)
	if err != nil {
		fmt.Println("Failed to set websocket upgrade: %+v", err)
		return
	}
	defer conn.Close()
	// 剩下的SHELL处理
	shell.Shell(conn)
}
