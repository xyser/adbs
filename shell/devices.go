package shell

import (
	"bufio"
	"errors"
	"io"
	"os/exec"
	"strings"
)

type Device struct {
	No    string
	State string
}

func Lists() (devices []Device, err error) {
	cmd := exec.Command("adb", "devices")

	var stdout io.ReadCloser
	if stdout, err = cmd.StdoutPipe(); err != nil {
		return nil, err
	}

	if err = cmd.Start(); err != nil {
		return nil, err
	}

	//创建一个流来读取管道内内容，这里逻辑是通过一行一行的读取的
	reader := bufio.NewReader(stdout)

	//实时循环读取输出流中的一行内容
	for {
		line, err := reader.ReadString('\n')
		if err != nil || io.EOF == err {
			break
		}
		if device, err := strToDevice(line); err == nil {
			devices = append(devices, device)
		}
	}
	return devices, nil
}

func strToDevice(s string) (device Device, err error) {
	if strings.HasPrefix(s, "*") || s == "List of devices attached" {
		return device, errors.New("not device")
	}
	strArr := strings.Fields(strings.TrimSpace(s))
	if len(strArr) == 2 {
		device.No = strArr[0]
		device.State = strArr[1]
		return device, nil
	}
	return device, errors.New("string format error")
}
