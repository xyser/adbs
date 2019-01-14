package main

import (
	"adbs/shell"
	"adbs/shell/packages"
	"fmt"
)

func main() {
	var bo bool
	var err error
	if bo, err = shell.Connect("192.168.11.29"); err != nil {
		fmt.Println(err)
	}
	fmt.Printf("Connect: %s\n", bo)

	// 设备列表
	devices, err := shell.Lists()
	if err != nil {
		fmt.Printf("devices error: %s", err.Error())
	}
	fmt.Println(devices[0].No)

	// adb version
	version, err := shell.Version()
	fmt.Println(version)

	// 软件包列表
	pack, err := packages.List()
	fmt.Println(pack)

	// 清理软件包缓存
	if bo, err = packages.Clear("com.dingdayu.helloandriod"); err != nil {
		fmt.Println(err)
	}
	fmt.Printf("Packages.Clear: %s\n", bo)

	if bo, err = packages.StartByIntent("com.dingdayu.helloandroid/.MainActivity", map[string]string{}); err != nil {
		fmt.Println(err)
	}
	fmt.Printf("Packages.StartByIntent: %s\n", bo)

	//阻塞直到该命令执行完成，该命令必须是被Start方法开始执行的
	//cmd.Wait()

	//stdin.Write([]byte("go text for grep\n"))
	//stdin.Write([]byte("go test text for grep\n"))
	//stdin.Close()
	//
	//out_bytes, _ := ioutil.ReadAll(stdout)
	//
	//if err := cmd.Wait(); err != nil {
	//	fmt.Println("Execute failed when Wait:" + err.Error())
	//	return
	//}
	//
	//fmt.Println("Execute finished:" + string(out_bytes))
}
