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
	"os"
	"time"
)

const TEMP_PATH = "/data/local/tmp"

type Machine struct {
	Client Client
	Serial string
	Conn   net.Conn
}

// 选择一个设备
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

// 上传一个文件
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
	_, err := m.Conn.Write(buf.Bytes())
	if err != nil {
		fmt.Println("write send error: " + err.Error())
	}

	file, _ := fh.Open()
	inputReader := bufio.NewReader(file)

	for {
		track := make([]byte, 65536)
		n, err := inputReader.Read(track)
		if err == io.EOF {
			// 读取完毕
			done := bytes.NewBuffer([]byte{})
			done.WriteString(DONE)
			done.WriteString(string(time.Now().Unix()))
			_, err = m.Conn.Write(done.Bytes())
			if err != nil {
				fmt.Println("write done error: " + err.Error())
			}
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

// 拉取一个文件
func (m Machine) Pull(remote string) ([]byte, error) {
	m = m.Sync()

	readChan := make(chan []byte)
	errChan := make(chan error)
	go func() {
		buf := bytes.NewBuffer([]byte{})

		for {
			stat := make([]byte, 8)
			n, err := m.Conn.Read(stat)
			if err != nil {
				errChan <- err
			}

			switch string(stat[0:4]) {
			case DATA:
				// 读长度
				length := binary.LittleEndian.Uint32(stat[4:n])

				// 循环读取文件
				for length > 0 {
					content := make([]byte, length)
					n, err := m.Conn.Read(content)
					if err != nil {
						errChan <- err
					}
					if n > 0 {
						buf.Write(content[0:n])
						length = length - uint32(n)
					}
				}
			case FAIL:
				errChan <- errors.New("adb response: FAIL")
			case DONE:
				readChan <- buf.Bytes()
			default:
				errChan <- errors.New("adb response: " + string(stat))
			}
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

// 文件信息结构体
type Stat struct {
	Name    string
	Size    int64
	Mode    os.FileMode
	ModTime time.Time
}

// 获取路径文件信息
func (m Machine) Stat(remote string) (Stat, error) {
	m = m.Sync()

	statChan := make(chan Stat)
	errChan := make(chan error)
	go func() {
		buffer := make([]byte, 16)

		n, err := m.Conn.Read(buffer)
		if err != nil {
			errChan <- err
		}
		if n > 0 {
			var stat Stat
			if string(buffer[0:4]) == STAT {
				fmt.Println(buffer[4:8])
				mode := binary.LittleEndian.Uint32(buffer[4:8])
				if mode == 0 {
					errChan <- errors.New("no find path")
				}
				fmt.Println(mode)
				stat.Mode = os.FileMode(mode) // 文件权限
				stat.Size = int64(binary.LittleEndian.Uint32(buffer[8:12]))
				stat.ModTime = time.Unix(int64(binary.LittleEndian.Uint32(buffer[12:n])), 0)
				statChan <- stat
			} else {
				errChan <- errors.New("resp error: " + string(buffer))
			}
		}
		errChan <- errors.New("socket read error")
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

	select {
	case stat := <-statChan:
		return stat, nil
	case err := <-errChan:
		return Stat{}, err
	}
}

// 获取路径目录
func (m Machine) Dir(path string) ([]Stat, error) {
	m = m.Sync()

	statChan := make(chan []Stat)
	errChan := make(chan error)
	go func() {
		var stats []Stat
		for {
			buffer := make([]byte, 4)

			_, err := m.Conn.Read(buffer)
			if err != nil {
				errChan <- err
			}
			switch string(buffer) {
			case DENT:
				buffer = make([]byte, 16)
				n, err := m.Conn.Read(buffer)
				if err != nil {
					errChan <- err
				}
				var stat Stat
				stat.Mode = os.FileMode(binary.LittleEndian.Uint32(buffer[0:4])) // 文件权限
				stat.Size = int64(binary.LittleEndian.Uint32(buffer[4:8]))
				stat.ModTime = time.Unix(int64(binary.LittleEndian.Uint32(buffer[8:12])), 0)

				// 读文件名
				nameLen := binary.LittleEndian.Uint32(buffer[12:n])
				nameBuf := make([]byte, nameLen)
				_, err = m.Conn.Read(nameBuf)
				if err != nil {
					errChan <- err
				}
				if string(nameBuf) != "." && string(nameBuf) != ".." {
					stat.Name = string(nameBuf)
					stats = append(stats, stat)
				}
			case DONE:
				statChan <- stats
				break
			default:
				fmt.Println(string(buffer))
			}
		}
	}()

	buf := new(bytes.Buffer)
	buf.WriteString(LIST)
	buf.Write(Uint32ToBytes(uint32(len(path))))
	buf.WriteString(path)
	// 写入发送命令
	_, err := m.Conn.Write(buf.Bytes())
	if err != nil {
		fmt.Println("write send error: " + err.Error())
	}

	select {
	case stat := <-statChan:
		return stat, nil
	case err := <-errChan:
		// `var t []string` 声明了一个nil slice
		// `t := []string{}` 声明了一个长度为0的非nil的slice。
		var stats []Stat
		return stats, err
	}
}
