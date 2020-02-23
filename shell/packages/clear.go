package packages

import (
	"bytes"
	"errors"
	"os/exec"
	"strings"
)

// Clear 删除与软件包关联的所有数据
func Clear(pack string) (bool, error) {
	//函数返回一个*Cmd，用于使用给出的参数执行name指定的程序
	cmd := exec.Command("adb", "shell", "pm", "clear", pack)

	//读取io.Writer类型的cmd.Stdout，再通过bytes.Buffer(缓冲byte类型的缓冲器)将byte类型转化为string类型(out.String():这是bytes类型提供的接口)
	var out bytes.Buffer
	cmd.Stdout = &out

	//Run执行c包含的命令，并阻塞直到完成。  这里stdout被取出，cmd.Wait()无法正确获取stdin,stdout,stderr，则阻塞在那了
	if err := cmd.Run(); err != nil {
		return false, err
	}

	var res = strings.Trim(out.String(), "\n")
	if strings.Contains(res, "Failed") {
		return false, errors.New("Failed")
	}
	if strings.Contains(res, "Success") {
		return true, nil
	}
	return false, errors.New(out.String())
}
