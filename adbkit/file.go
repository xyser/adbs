package adbkit

import (
	"bufio"
	"bytes"
	"encoding/binary"
	"errors"
	"fmt"
	"io"
	"mime/multipart"
	"net"
	"time"
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
	buf.WriteString(SEND)
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
			done.WriteString(DONE)
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
		data.WriteString(DATA)
		data.Write(Uint32ToBytes(uint32(n)))
		_, _ = m.Conn.Write(data.Bytes())

		// 发送内容
		_, _ = m.Conn.Write(track[0:n])
	}

}

func (m Machine) Pull(remote string) ([]byte, error) {
	m = m.Sync()

	readChan := make(chan []byte)
	errChan := make(chan error)
	go func() {
		stat := make([]byte, 4)
		_, _ = m.Conn.Read(stat)
		switch string(stat) {
		case DATA:
			// 读长度
			leng := make([]byte, 4)
			_, _ = m.Conn.Read(leng)
			length := binary.LittleEndian.Uint32(leng)

			if length > 0 {
				content := make([]byte, length)
				n, _ := m.Conn.Read(content)
				if n > 0 {
					done := make([]byte, 4)
					_, _ = m.Conn.Read(done)
					if string(done) == DONE {
						readChan <- content
					}
				}
			}
		case FAIL:
			errChan <- errors.New("adb response: FAIL")
		}
	}()

	// 写入命令
	buf := new(bytes.Buffer)
	buf.WriteString(RECV)
	buf.Write(Uint32ToBytes(uint32(len(remote))))
	buf.WriteString(remote)
	_, err := m.Conn.Write(buf.Bytes())
	if err != nil {
		return nil, err
	}

	select {
	case content := <-readChan:
		return content, nil
	case err := <-errChan:
		return nil, err
	}
}

func (m Machine) Stat(remote string) {
	m = m.Sync()

	go func() {
		buffer := make([]byte, 1024)
		for {
			n, err := m.Conn.Read(buffer)
			if err == io.EOF {
				fmt.Println("read finish")
			}
			if n > 0 {
				fmt.Println("resp: " + string(buffer[0:n]))
				fmt.Println(string(buffer[0:4]))
				fmt.Println(binary.LittleEndian.Uint32(buffer[4:8]))
				fmt.Println(binary.LittleEndian.Uint32(buffer[8:12]))
				fmt.Println(binary.LittleEndian.Uint32(buffer[12:n]))
				fmt.Println(n)

			}
		}
	}()

	buf := new(bytes.Buffer)
	buf.WriteString(STAT)
	buf.Write(Uint32ToBytes(uint32(len(remote))))
	buf.WriteString(remote)
	// 写入发送命令
	_, err := m.Conn.Write(buf.Bytes())
	if err != nil {
		fmt.Println("write send error: " + err.Error())
	}

	time.Sleep(10 * time.Second)
}
