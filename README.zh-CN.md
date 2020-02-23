# adbs

Android Debug Bridge (adb) 远程调试服务。

其他语言请阅读： [English](README.md)， [简体中文](README.zh-CN.md)。

## 功能列表

* ADB 服务
* WEB 控制
* 远程控制

## 扩展包

项目采用 `go mod` 方案，  引用了以下第三方包：

- https://github.com/gin-gonic/gin
- https://github.com/gorilla/websocket
- https://github.com/shogo82148/androidbinary
- https://github.com/kr/pty

## Roadmap

* 支持多人 `web shell`.
* 输出 `screencap` 实现截图显示.
* 连接和断开设备连接.
* 获取设备得软件包列表.
* 上传和下载设备上得文件.
* 控制设备输入.
* 改进设备截图方案

## 快速开始

### ADB

请先下载新版本得 `platform-tools`, [platform-tools](https://developer.android.com/studio/releases/platform-tools).

下载完成后，请将解压所得得路径，添加到系统 `PATH` 里面.

使其可以在控制台里，直接成功运行 `adb devices`.

### Download

```shell
git clone https://github.com/dingdayu/adbs

cd adbs

go mod tidy

go run .
```

## TODO

- [X] 写文件时的时间
- [X] 推文件写入协议优化
- [ ] 获取文件信息时的文件类型问题（目录/连接）
- [ ] 设备列表等API接口完成设备选择
- [ ] 提供编译版本
- [ ] 提供 `docker` 镜像


## 参考

该项目受以下项目或文章的影响: 

- [7.0上截图的问题](https://github.com/mzlogin/awesome-adb/issues/33)
- [Android之高效率截图](https://juejin.im/post/5bab409ef265da0afc2c032e)
- [Read binary stdout data from adb shell?](https://stackoverflow.com/questions/13578416/read-binary-stdout-data-from-adb-shell)
- [Go 内嵌静态资源](http://fuxiaohei.me/2016/10/1/go-binary-embed-asset.html)