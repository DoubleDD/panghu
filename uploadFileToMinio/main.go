package main

import (
	"fmt"
	"minioUploadFile/server"
	"minioUploadFile/upload"
	"minioUploadFile/video"
	"os"
)

func main() {
	if len(os.Args) == 1 {
		help()
		return
	}

	switch os.Args[1] {
	case "server":
		server.Start(7700)
	case "m3u8":
		video.Convert(os.Args[2], os.Args[3])
	case "upload":
		upload.ParallelUpload(os.Args[2], os.Args[3])
	default:
		fmt.Println("参数错误")
	}
}

func help() {
	// 显示使用帮助
	message := `
minioUploadFile [server|m3u8|upload] [input] [output]

1. 第一个参数表示工具类型，有三个取值：server、m3u8、upload
	server: 启动一个web服务，端口7700。
	m3u8:	将视频转成m3u8格式，后面需再加两个参数
		1. 视频文件路径，相对/绝对路径都可以
		2. 输出目录
	upload: 将文件上传至minio，minio相关配置可在当前目录的config.toml文件将修改,使用多协程+分片的方式加速上传,也需要两个参数：
		1. 要上传的目录/文件
		2. minio上面的路径
		`
	println(message)
}
