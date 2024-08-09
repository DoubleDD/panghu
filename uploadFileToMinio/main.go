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
		fmt.Println("缺少参数")
		return
	}
	fnType := os.Args[1]
	localFile := os.Args[2]
	minioPath := os.Args[3]

	switch fnType {
	case "server":
		server.Start(7700)
	case "m3u8":
		video.Convert(localFile, minioPath)
	case "upload":
		upload.ParallelUpload(localFile, minioPath)
	default:
		fmt.Println("参数错误")

	}
}
