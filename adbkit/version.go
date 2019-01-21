package adbkit

import (
	"errors"
	"strconv"
)

// 获取 ADB 版本
func (c Client) Version() (int, error) {
	resp, err := c.Command("host:version")
	if err != nil {
		return 0, err
	}

	if string(resp[0:4]) == OKAY {
		length, _ := strconv.Atoi(string(resp[4:8]))
		version, _ := strconv.Atoi(string(resp[8 : 8+length]))
		return version, nil
	} else if string(resp[0:4]) == FAIL {
		return 0, errors.New("adb response: Fail")
	}
	return 0, errors.New("error response: " + string(resp))
}
