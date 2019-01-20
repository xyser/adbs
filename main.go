package main

import "adbs/api"

func main() {
	api.Init()

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
