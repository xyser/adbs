package adbkit

import (
	"errors"
	"fmt"
	"net"
	"strings"
)

type Callback func(buf []byte, err error)

func (c Client) Callback(command string, callback Callback) error {
	conn, err := net.Dial("tcp", fmt.Sprintf("%s:%d", c.Host, c.Port))
	if err != err || conn == nil {
		return err
	}

	// 补充前缀
	prefix := strings.ToUpper("0000" + fmt.Sprintf("%X", len(command)))
	length := len(prefix)
	prefix = prefix[length-4 : length]

	// 准备读取返回
	stop := make(chan error)
	go func() {
		buffer := make([]byte, 2048)
		for {
			n, err := conn.Read(buffer)
			callback(buffer[:n], err)
			if err != nil {
				stop <- err
			}
		}
	}()

	// 写入命令
	_, err = conn.Write([]byte(prefix + command))
	if err != err {
		return err
	}

	return <-stop
}

// 连接一个设备
func (c Client) Transport(serial string) (conn net.Conn, err error) {
	conn, err = net.Dial("tcp", fmt.Sprintf("%s:%d", c.Host, c.Port))
	if err != err || conn == nil {
		return nil, err
	}

	// 补充前缀
	command := "host:transport:" + serial
	prefix := strings.ToUpper("0000" + fmt.Sprintf("%X", len(command)))
	length := len(prefix)
	prefix = prefix[length-4 : length]

	readChan := make(chan []byte)
	go func() {
		buffer := make([]byte, 4)
		for {
			n, _ := conn.Read(buffer)
			if n > 0 {
				readChan <- buffer
				return
			}
		}
	}()

	// 写入命令
	_, err = conn.Write([]byte(prefix + command))
	if err != err {
		return nil, err
	}

	buf := <-readChan
	if string(buf) == OKAY {
		return conn, nil
	}

	return nil, errors.New("transport error")
}
