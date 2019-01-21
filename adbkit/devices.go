package adbkit

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
)

const OKAY = "OKAY"
const FAIL = "FAIL"

type Device struct {
	No    string `json:"no"`
	State string `json:"state"`
}

// 获取设备列表
func (c Client) Devices(state bool) ([]Device, error) {
	command := "host:devices"
	if state {
		command = "host:devices-l"
	}
	resp, err := c.Command(command)
	if err != nil {
		return nil, err
	}
	if string(resp[0:4]) == OKAY {
		var devices []Device
		for _, line := range strings.Split(string(resp[8:]), "\n") {
			device := strings.Split(line, "\t")
			if len(device) > 1 {
				devices = append(devices, Device{No: device[0], State: device[1]})
			}
		}
		return devices, nil
	} else if string(resp[0:4]) == FAIL {
		return nil, errors.New("adb response: Fail")
	}
	return nil, errors.New("error response: " + string(resp))
}

// 连接一个设备
func (c Client) Connect(ip string, port int) (bool, error) {
	resp, err := c.Command(fmt.Sprintf("host:connect:#%s:#%d", ip, port))
	if err != nil {
		return false, err
	}
	if string(resp[0:4]) == OKAY {
		length, _ := strconv.Atoi(string(resp[4:8]))

		var res = strings.Trim(string(resp[8:8+length]), "\n")
		if strings.Contains(res, "failed to connect") || strings.Contains(res, "unable to connect to") {
			return false, errors.New("failed to connect device")
		}
		if strings.Contains(res, "already connected to") || strings.Contains(res, "connected to") {
			return true, nil
		}
	} else if string(resp[0:4]) == FAIL {
		return false, errors.New("adb response: Fail")
	}
	return false, errors.New("error response: " + string(resp))
}

// 断开设备
func (c Client) Disconnect(ip string, port int) (bool, error) {
	resp, err := c.Command(fmt.Sprintf("host:disconnect:#%s:#%d", ip, port))
	if err != nil {
		return false, err
	}
	if string(resp[0:4]) == OKAY {
		length, _ := strconv.Atoi(string(resp[4:8]))
		fmt.Println(length)
		fmt.Println(resp)
		var res = strings.Trim(string(resp[8:8+length]), "\n")
		if strings.Contains(res, "No such device") {
			return false, errors.New("no such device")
		}
		return true, nil
	} else if string(resp[0:4]) == FAIL {
		return false, errors.New("adb response: Fail")
	}
	return false, errors.New("error response: " + string(resp))
}
