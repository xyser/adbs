package adbkit

import (
	"bufio"
	"bytes"
	"encoding/binary"
	"fmt"
	"io"
	"mime/multipart"
	"net"
)

type Machine struct {
	Client Client
	Serial string
	Conn   net.Conn
}

func (c Client) Select(serial string) Machine {
	conn, _ := net.Dial("tcp", fmt.Sprintf("%s:%d", c.Host, c.Port))

	// 写入命令
	command := "host:transport:" + serial
	_, _ = conn.Write(EncodeCommend(command))

	buf := make([]byte, 4)
	_, _ = conn.Read(buf)

	if string(buf) == OKAY {
		return Machine{Client: c, Serial: serial, Conn: conn}
	}
	return Machine{Client: c, Serial: serial}
}

func (m Machine) Sync() Machine {
	// 写入命令
	command := "sync:"
	_, _ = m.Conn.Write(EncodeCommend(command))

	buf := make([]byte, 4)
	_, _ = m.Conn.Read(buf)

	if string(buf) == OKAY {
		return m
	}
	return m
}

func (m Machine) Push(fh *multipart.FileHeader, remote string) {
	// 写入命令
	path := remote + ",0644"

	buf := bytes.NewBuffer([]byte{})
	buf.WriteString("SEND")
	_ = binary.Write(buf, binary.BigEndian, len(path))
	buf.WriteString(path)
	// 写入发送命令
	_, _ = m.Conn.Write(buf.Bytes())

	file, _ := fh.Open()
	inputReader := bufio.NewReader(file)

	for {

		track := make([]byte, 65536)
		n, err := inputReader.Read(track)
		if err == io.EOF {
			// 读取完毕
			done := bytes.NewBuffer([]byte{})
			done.WriteString("DONE")
			_ = binary.Write(done, binary.BigEndian, n)
			_, _ = m.Conn.Write(done.Bytes())
			break
		}

		// 发送长度
		data := bytes.NewBuffer([]byte{})
		data.WriteString("DATA")
		_ = binary.Write(data, binary.BigEndian, n)
		_, _ = m.Conn.Write(data.Bytes())

		// 发送内容
		_, _ = m.Conn.Write(track[0:n])
	}

}
