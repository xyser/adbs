package adbkit

import (
	"errors"
	"io/ioutil"
	"regexp"
)

func (c Client) GetProperties(serial string) (prop map[string]string, err error) {
	conn, err := c.Transport(serial)
	if err != nil {
		return nil, err
	}
	// 准备读取返回
	readChan := make(chan []byte)
	go func() {
		buf, _ := ioutil.ReadAll(conn)
		readChan <- buf
	}()

	// 写入命令
	command := "shell:getprop"
	_, err = conn.Write(EncodeCommend(command))
	if err != err {
		return nil, err
	}

	resp := <-readChan
	if string(resp[0:4]) == OKAY {
		var re = regexp.MustCompile(`(?m)\[(.+)\]: \[(.+)\]`)
		for _, match := range re.FindAllStringSubmatch(string(resp[0:4]), -1) {
			if len(match) > 2 {
				prop[match[1]] = match[2]
			}
		}
		return prop, nil
	} else if string(resp[0:4]) == FAIL {
		return nil, errors.New("adb fail response: " + string(resp[4:]))
	}

	return nil, errors.New("adb response: " + string(resp))
}
