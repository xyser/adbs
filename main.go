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
