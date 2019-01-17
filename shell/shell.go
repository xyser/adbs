package shell

import (
	"bufio"
	"fmt"
	"os/exec"
)

func Shell(in chan []byte, out chan []byte) error {
	//函数返回一个*Cmd，用于使用给出的参数执行name指定的程序
	cmd := exec.Command("adb", "shell")

	//读取io.Writer类型的cmd.Stdout，再通过bytes.Buffer(缓冲byte类型的缓冲器)将byte类型转化为string类型(out.String():这是bytes类型提供的接口)
	//var out bytes.Buffer
	//var in bytes.Buffer
	//cmd.Stdout = &out
	//cmd.Stdin = &in
	outRead, _ := cmd.StdoutPipe()
	inRead, _ := cmd.StdinPipe()

	outScanner := bufio.NewScanner(outRead)
	outScanner.Split(bufio.ScanWords)
	for outScanner.Scan() {
		out <- outScanner.Bytes()
	}

	for {
		msg := <-in
		if _, err := inRead.Write(msg); err != nil {
			fmt.Println(err)
		}
	}

	//Run执行c包含的命令，并阻塞直到完成。  这里stdout被取出，cmd.Wait()无法正确获取stdin,stdout,stderr，则阻塞在那了
	return cmd.Run()
}
