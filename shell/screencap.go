package shell

import (
	"bytes"
	"os/exec"
)

func Screencap() ([]byte, error) {
	cmd := exec.Command("adb", "exec-out", "screencap -p")

	//读取io.Writer类型的cmd.Stdout，再通过bytes.Buffer(缓冲byte类型的缓冲器)将byte类型转化为string类型(out.String():这是bytes类型提供的接口)
	var out bytes.Buffer
	cmd.Stdout = &out

	//Run执行c包含的命令，并阻塞直到完成。  这里stdout被取出，cmd.Wait()无法正确获取stdin,stdout,stderr，则阻塞在那了
	if err := cmd.Run(); err != nil {
		return nil, err
	}

	return out.Bytes(), nil
}
