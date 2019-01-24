package shell

import (
	"bufio"
	"bytes"
	"errors"
	"io"
	"os/exec"
	"strings"
)

type Device struct {
	No    string `json:"no"`
	State string `json:"state"`
}

// 连接一个IP的设备
func Connect(ip string) (bool, error) {
	cmd := exec.Command("adb", "connect", ip)

	//读取io.Writer类型的cmd.Stdout，再通过bytes.Buffer(缓冲byte类型的缓冲器)将byte类型转化为string类型(out.String():这是bytes类型提供的接口)
	var out bytes.Buffer
	cmd.Stdout = &out

	//Run执行c包含的命令，并阻塞直到完成。  这里stdout被取出，cmd.Wait()无法正确获取stdin,stdout,stderr，则阻塞在那了
	if err := cmd.Run(); err != nil {
		return false, err
	}

	var res = strings.Trim(out.String(), "\n")
	if strings.Contains(res, "failed to connect") {
		return false, errors.New(res)
	}
	if strings.Contains(res, "already connected to") || strings.Contains(res, "connected to") {
		return true, nil
	}
	return false, errors.New(out.String())
}

// 解除设备连接
func Disconnect(ip string) (bool, error) {
	var cmd *exec.Cmd
	if ip == "all" {
		cmd = exec.Command("adb", "disconnect")
	} else {
		cmd = exec.Command("adb", "disconnect", ip)
	}

	//读取io.Writer类型的cmd.Stdout，再通过bytes.Buffer(缓冲byte类型的缓冲器)将byte类型转化为string类型(out.String():这是bytes类型提供的接口)
	var out bytes.Buffer
	cmd.Stdout = &out

	//Run执行c包含的命令，并阻塞直到完成。  这里stdout被取出，cmd.Wait()无法正确获取stdin,stdout,stderr，则阻塞在那了
	if err := cmd.Run(); err != nil {
		return false, err
	}

	var res = strings.Trim(out.String(), "\n")
	// 没有找到设备连接
	if strings.Contains(res, "no such device") {
		return false, errors.New(res)
	}
	// 成功接触连接
	if strings.Contains(res, "disconnected") {
		return true, nil
	}
	return false, errors.New(out.String())
}

// 获取连接的设备列表
// Deprecated: 会逐渐采用 adbkit 代替
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

// 内部字符串转设备类型
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
