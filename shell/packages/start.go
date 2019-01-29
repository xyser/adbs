package packages

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"os/exec"
)

// StartByIntent 启动默认设备包
func StartByIntent(intent string, data map[string]string) (bool, error) {

	arg := []string{"shell", "am", "start", "-n"}
	arg = append(arg, intent)
	if data, err := toData(data); err == nil {
		arg = append(arg, data)
	}

	cmd := exec.Command("adb", arg...)

	var err error
	var stdout io.ReadCloser
	if stdout, err = cmd.StdoutPipe(); err != nil {
		return false, err
	}

	if err = cmd.Start(); err != nil {
		return false, err
	}

	//创建一个流来读取管道内内容，这里逻辑是通过一行一行的读取的
	reader := bufio.NewReader(stdout)

	//实时循环读取输出流中的一行内容
	for {
		line, err := reader.ReadString('\n')
		if err != nil || io.EOF == err {
			break
		}
		fmt.Println(line)
	}
	//return packages,nil
	return false, nil
}

// StartByAction 通过 Action 启动一个包
func StartByAction(action string, data map[string]string) (bool, error) {

	arg := []string{"shell", "am", "start", "-a"}
	arg = append(arg, action)
	if data, err := toData(data); err == nil {
		arg = append(arg, data)
	}

	cmd := exec.Command("adb", arg...)

	var err error
	var stdout io.ReadCloser
	if stdout, err = cmd.StdoutPipe(); err != nil {
		return false, err
	}

	if err = cmd.Start(); err != nil {
		return false, err
	}

	//创建一个流来读取管道内内容，这里逻辑是通过一行一行的读取的
	reader := bufio.NewReader(stdout)

	//实时循环读取输出流中的一行内容
	for {
		line, err := reader.ReadString('\n')
		if err != nil || io.EOF == err {
			break
		}
		fmt.Println(line)
	}
	//return packages,nil
	return false, nil
}

func toData(data map[string]string) (string, error) {
	if len(data) <= 0 {
		return "", errors.New("nil")
	}
	var ret string
	for k, v := range data {
		ret = ret + " -d " + k + ":" + v
	}
	return ret, nil
}
