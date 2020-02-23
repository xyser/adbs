package shell

import (
	"bytes"
	"os/exec"
	"regexp"
	"strconv"
)

type Ver struct {
	Version     string
	VersionCode int
	CommandPath string
}

// Version 通过 shell 获取版本
func Version() (version Ver, err error) {
	//函数返回一个*Cmd，用于使用给出的参数执行name指定的程序
	cmd := exec.Command("adb", "version")

	//读取io.Writer类型的cmd.Stdout，再通过bytes.Buffer(缓冲byte类型的缓冲器)将byte类型转化为string类型(out.String():这是bytes类型提供的接口)
	var out bytes.Buffer
	cmd.Stdout = &out

	//Run执行c包含的命令，并阻塞直到完成。  这里stdout被取出，cmd.Wait()无法正确获取stdin,stdout,stderr，则阻塞在那了
	if err := cmd.Run(); err != nil {
		return version, err
	}

	var re = regexp.MustCompile(`(?msi)^Android Debug Bridge version (\S+)\nVersion (\d+)\nInstalled as (\S+)$`)
	arr := re.FindStringSubmatch(out.String())
	if len(arr) >= 4 {
		version.Version = arr[1]
		version.VersionCode, _ = strconv.Atoi(arr[2])
		version.CommandPath = arr[3]
	}

	return version, nil
}
