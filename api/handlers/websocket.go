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

	in := make(chan []byte)
	out := make(chan []byte)

	go func() {
		for {
			msg := <-out
			err = conn.WriteMessage(1, msg)
			if err != nil {
				fmt.Println(err)
				break
			}
		}
	}()

	go shell.Shell(in, out)

	// 必须死循环，gin通过协程调用该handler函数，一旦退出函数，ws会被主动销毁
	for {
		_, p, err := conn.ReadMessage()
		if err != nil {
			break
		}
		in <- p

		//if err := conn.WriteMessage(1, <-out); err != nil {
		//	fmt.Println(out)
		//	fmt.Println("Failed to set websocket upgrade: %+v", err)
		//	break
		//}
	}
}
