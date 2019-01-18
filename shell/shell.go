package shell

import (
	"bufio"
	"fmt"
	"os/exec"
)

func Shell(inC chan []byte, outC chan []byte) {
	//函数返回一个*Cmd，用于使用给出的参数执行name指定的程序
	cmd := exec.Command("adb", "shell")

	//读取io.Writer类型的cmd.Stdout，再通过bytes.Buffer(缓冲byte类型的缓冲器)将byte类型转化为string类型(out.String():这是bytes类型提供的接口)
	//var out,in bytes.Buffer

	//cmd.Stdout = &out
	//cmd.Stdin = &in

	w, _ := cmd.StdinPipe()
	r, _ := cmd.StdoutPipe()

	go func() {
		s := bufio.NewScanner(r)
		for s.Scan() {
			fmt.Println(s.Text())
		}
	}()

	if err := cmd.Start(); err != nil {
		panic(err)
	}

	for {
		msg := <-inC
		fmt.Println("->" + string(msg))
		if _, err := w.Write(msg); err != nil {
			fmt.Println(err)
		}
	}
	go cmd.Wait()
}
