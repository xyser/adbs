package packages

import (
	"bufio"
	"errors"
	"io"
	"os/exec"
	"regexp"
)

type Package struct {
	ApplicationId string
	ApkPath       string
}

func List() (packages []Package, err error) {
	cmd := exec.Command("adb", "shell", "pm", "list", "packages", "-f")

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
		if pack, err := strToPackage(line); err == nil {
			packages = append(packages, pack)
		}
	}
	return packages, nil
}

func strToPackage(str string) (pack Package, err error) {
	var re = regexp.MustCompile(`(?msi)(\S+):(\S+)=(\S+)`)
	arr := re.FindStringSubmatch(str)
	if len(arr) >= 4 {
		pack.ApkPath = arr[2]
		pack.ApplicationId = arr[3]
	} else {
		return pack, errors.New("package nil")
	}
	return pack, nil
}
