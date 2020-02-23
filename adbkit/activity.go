package adbkit

import (
	"errors"
	"fmt"
	"io/ioutil"
	"regexp"
	"strings"
)

// ActivityManager Activity 管理
func (c Client) ActivityManager(serial, command string, args []string) (bool, error) {
	conn, err := c.Transport(serial)
	if err != nil {
		return false, err
	}
	// 准备读取返回
	readChan := make(chan []byte)
	go func() {
		buf, _ := ioutil.ReadAll(conn)
		readChan <- buf
	}()

	// 写入命令
	if len(args) > 0 {
		command = command + strings.Join(args, " ")
	}
	command = fmt.Sprintf("shell:am %s", command)
	_, err = conn.Write(EncodeCommend(command))
	if err != err {
		return false, err
	}

	resp := <-readChan
	if string(resp[0:4]) == OKAY {
		var re = regexp.MustCompile(`(?m)^Error: (.*)$`)
		match := re.FindStringSubmatch(string(resp[0:4]))
		if len(match) > 1 {
			return false, errors.New(match[1])
		}
		return true, nil
	} else if string(resp[0:4]) == FAIL {
		return false, errors.New("adb fail response: " + string(resp[4:]))
	}

	return false, errors.New("adb response: " + string(resp))
}
