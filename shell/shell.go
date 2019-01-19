package shell

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/websocket"
	"github.com/kr/pty"
	"io"
	"log"
	"os"
	"os/exec"
	"syscall"
	"unsafe"
)

type windowSize struct {
	Rows uint16 `json:"rows"`
	Cols uint16 `json:"cols"`
	X    uint16
	Y    uint16
}

func Shell(conn *websocket.Conn) {
	//函数返回一个*Cmd，用于使用给出的参数执行name指定的程序
	cmd := exec.Command("adb", "shell")
	cmd.Env = append(os.Environ(), "TERM=xterm")

	tty, err := pty.Start(cmd)
	if err != nil {
		log.Println("Unable to start pty/cmd")
		_ = conn.WriteMessage(websocket.TextMessage, []byte(err.Error()))
		return
	}
	defer func() {
		_ = cmd.Process.Kill()
		_, _ = cmd.Process.Wait()
		_ = tty.Close()
		_ = conn.Close()
	}()

	go func() {
		for {
			buf := make([]byte, 1024)
			read, err := tty.Read(buf)
			if err != nil {
				if err == io.EOF {
					_ = conn.WriteMessage(websocket.TextMessage, []byte("Warn: ADB 断开连接"))
					fmt.Println("Warn: ADB 断开连接")
				} else {
					_ = conn.WriteMessage(websocket.TextMessage, []byte(err.Error()))
					log.Println("Error: Unable to read from pty/cmd")
				}
				return
			}
			_ = conn.WriteMessage(websocket.BinaryMessage, buf[:read])
		}
	}()

	for {
		messageType, reader, err := conn.NextReader()
		if err != nil {
			log.Println("Error: Unable to grab next reader")
			return
		}

		if messageType == websocket.TextMessage {
			log.Println("Warn: Unexpected text message")
			_ = conn.WriteMessage(websocket.TextMessage, []byte("Unexpected text message"))
			continue
		}

		dataTypeBuf := make([]byte, 1)
		read, err := reader.Read(dataTypeBuf)
		if err != nil {
			log.Println("Error: Unable to read message type from reader")
			_ = conn.WriteMessage(websocket.TextMessage, []byte("Unable to read message type from reader"))
			return
		}

		if read != 1 {
			log.Println("Error: Unexpected number of bytes read")
			return
		}

		switch dataTypeBuf[0] {
		case 0:
			copied, err := io.Copy(tty, reader)
			if err != nil {
				log.Printf("Error: Error after copying %d bytes\n", copied)
			}
		case 1:
			decoder := json.NewDecoder(reader)
			resizeMessage := windowSize{}
			err := decoder.Decode(&resizeMessage)
			if err != nil {
				_ = conn.WriteMessage(websocket.TextMessage, []byte("Error: 解析窗口信息失败: "+err.Error()))
				continue
			}
			log.Printf("Info: Resizing terminal [%v]\n", resizeMessage)
			_, _, errno := syscall.Syscall(
				syscall.SYS_IOCTL,
				tty.Fd(),
				syscall.TIOCSWINSZ,
				uintptr(unsafe.Pointer(&resizeMessage)),
			)
			if errno != 0 {
				log.Printf("Error: 未能成功重置窗口大小[%s]\n", errno.Error())
			}
		default:
			log.Printf("Error: Unknown data type[%d]\n", dataTypeBuf[0])
		}
	}
}
