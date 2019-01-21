package main

import (
	"adbs/adbkit"
	"fmt"
)

func main() {
	//api.Init()

	//response, _ := adbkit.New("127.0.0.1", 5037).Command("host:devices")
	//fmt.Println(string(response))

	fmt.Println(adbkit.New("127.0.0.1", 5037).Devices())

	//// adb version
	//version, err := shell.Version()
	//fmt.Println(version)
	//
	//// 软件包列表
	//pack, err := packages.List()
	//fmt.Println(pack)
	//
	//// 清理软件包缓存
	//if bo, err = packages.Clear("com.dingdayu.helloandriod"); err != nil {
	//	fmt.Println(err)
	//}
	//fmt.Printf("Packages.Clear: %s\n", bo)
	//
	//if bo, err = packages.StartByIntent("com.dingdayu.helloandroid/.MainActivity", map[string]string{}); err != nil {
	//	fmt.Println(err)
	//}
	//fmt.Printf("Packages.StartByIntent: %s\n", bo)
}
