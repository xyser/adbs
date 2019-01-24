package adbkit

import (
	"errors"
	"fmt"
	"io/ioutil"
	"net"
)

type Client struct {
	Host string
	Port int
}

func New(host string, port int) Client {
	if !CheckTcp(host, port) {
		panic(errors.New("network fail"))
	}
	return Client{Host: host, Port: port}
}

// TODO:: 当adb service 未启动时调用存在问题
func (c Client) Command(command string) (response []byte, err error) {
	conn, err := net.Dial("tcp", fmt.Sprintf("%s:%d", c.Host, c.Port))
	if err != err || conn == nil {
		return nil, err
	}

	// 准备读取返回
	readChan := make(chan []byte)
	go func() {
		buf, _ := ioutil.ReadAll(conn)
		readChan <- buf
	}()

	// 写入命令
	_, err = conn.Write(EncodeCommend(command))
	if err != err {
		return nil, err
	}

	return <-readChan, nil
}

func CheckTcp(host string, port int) bool {
	if _, err := net.Dial("tcp", fmt.Sprintf("%s:%d", host, port)); err != nil {
		return false
	}
	return true
}
