package adbkit

import (
	"errors"
	"io/ioutil"
)

// Reboot 重启设备
// fmt.Println(adbkit.New("127.0.0.1", 5037).Reboot("emulator-5554"))
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

	// 写入命令
	_, err = conn.Write(EncodeCommend("reboot:"))
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

// Remount 重新挂载磁盘
func (c Client) Remount(serial string) error {
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

	// 写入命令
	_, err = conn.Write(EncodeCommend("remount:"))
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
