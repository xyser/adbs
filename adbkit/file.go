package adbkit

import (
	"bufio"
	"bytes"
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
		fmt.Println("host:transport: ok")
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
		fmt.Println("sync: ok")
		return m
	}
	return m
}

func (m Machine) Push(fh *multipart.FileHeader, remote string) {
	m = m.Sync()
	go func() {
		buffer := make([]byte, 1024)
		for {
			n, _ := m.Conn.Read(buffer)
			if n > 0 {
				fmt.Println("resp: " + string(buffer[0:n]))
			}
		}
	}()

	// 写入命令
	path := remote + ",0644"

	buf := new(bytes.Buffer)
	buf.WriteString("SEND")
	buf.Write(Uint32ToBytes(uint32(len(path))))
	buf.WriteString(path)
	// 写入发送命令
	n, err := m.Conn.Write(buf.Bytes())
	if err != nil {
		fmt.Println("write send error: " + err.Error())
	}
	fmt.Printf("bytes length: %d, %s\n", n, buf.Bytes())

	file, _ := fh.Open()
	inputReader := bufio.NewReader(file)

	for {

		track := make([]byte, 65536)
		n, err := inputReader.Read(track)
		if err == io.EOF {
			// 读取完毕
			done := bytes.NewBuffer([]byte{})
			done.WriteString("DONE")
			done.Write(Uint32ToBytes(uint32(n)))
			_, err = m.Conn.Write(done.Bytes())
			if err != nil {
				fmt.Println("write done error: " + err.Error())
			}

			fmt.Printf("upload done: %d\n", n)
			break
		}

		// 发送长度
		data := bytes.NewBuffer([]byte{})
		data.WriteString("DATA")
		data.Write(Uint32ToBytes(uint32(n)))
		_, _ = m.Conn.Write(data.Bytes())

		// 发送内容
		_, _ = m.Conn.Write(track[0:n])
	}

}

func (m Machine) Pull(fh *multipart.FileHeader, remote string) {
	m = m.Sync()

	go func() {
		buffer := make([]byte, 1024)
		for {
			n, _ := m.Conn.Read(buffer)
			if n > 0 {
				fmt.Println("resp: " + string(buffer[0:n]))
			}
		}
	}()

	// 写入命令
	buf := new(bytes.Buffer)
	buf.WriteString("RECV")
	buf.Write(Uint32ToBytes(uint32(len(remote))))
	buf.WriteString(remote)

}
