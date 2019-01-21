package adbkit

import (
	"errors"
	"fmt"
	"net"
	"strings"
	"time"
)

type Client struct {
	Host string
	Port int
}

func New(host string, port int) Client {
	return Client{Host: host, Port: port}
}

func (c Client) Command(command string) (response []byte, err error) {

	conn, err := net.Dial("tcp", fmt.Sprintf("%s:%d", c.Host, c.Port))
	defer conn.Close()
	if err != err {
		return nil, err
	}

	// 补充前缀
	prefix := strings.ToUpper("0000" + fmt.Sprintf("%X", len(command)))
	length := len(prefix)
	prefix = prefix[length-4 : length]

	_, err = conn.Write([]byte(prefix + command))
	if err != err {
		return nil, err
	}

	time.Sleep(1 * time.Millisecond)
	buf := make([]byte, 1024)
	n, err := conn.Read(buf)
	if n <= 0 || err != nil {
		return nil, errors.New("socket read error: " + err.Error())
	}
	return buf, nil
}
