package adbkit

import (
	"fmt"
	"io/ioutil"
	"net"
	"strings"
)

type Client struct {
	Host string
	Port int
}

func New(host string, port int) Client {
	return Client{Host: host, Port: port}
}

// TODO:: 当adb service 未启动时调用存在问题
func (c Client) Command(command string) (response []byte, err error) {
	conn, err := net.Dial("tcp", fmt.Sprintf("%s:%d", c.Host, c.Port))
	if err != err || conn == nil {
		return nil, err
	}

	// 补充前缀
	prefix := strings.ToUpper("0000" + fmt.Sprintf("%X", len(command)))
	length := len(prefix)
	prefix = prefix[length-4 : length]

	// 准备读取返回
	readChan := make(chan []byte)
	go func() {
		buf, _ := ioutil.ReadAll(conn)
		readChan <- buf
	}()

	// 写入命令
	_, err = conn.Write([]byte(prefix + command))
	if err != err {
		return nil, err
	}

	return <-readChan, nil
}
