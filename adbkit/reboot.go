package adbkit

import (
	"errors"
	"fmt"
	"io/ioutil"
	"strings"
)

func (c Client) Reboot(serial string) error {
	conn, err := c.Transport(serial)
	if err != nil {
		return err
	}
	// 准备读取返回
	readChan := make(chan []byte)
	go func() {
		buf, _ := ioutil.ReadAll(conn)
		readChan <- buf
	}()

	command := "reboot:"
	prefix := strings.ToUpper("0000" + fmt.Sprintf("%X", len(command)))
	length := len(prefix)
	prefix = prefix[length-4 : length]

	// 写入命令
	_, err = conn.Write([]byte(prefix + command))
	if err != err {
		return err
	}

	resp := <-readChan

	if string(resp[0:4]) == OKAY {
		return nil
	} else if string(resp[0:4]) == FAIL {
		return errors.New("adb response: Fail")
	}

	return errors.New("adb response: " + string(resp))
}
