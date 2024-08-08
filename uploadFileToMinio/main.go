package main

import (
	"fmt"
	"minioUploadFile/server"
	"minioUploadFile/upload"
	"minioUploadFile/video"
	"os"
)

func main() {
	server.Start(7700)
}

func localTool() {
	if len(os.Args) == 1 {
		fmt.Println("缺少参数")
		return
	}
	localFile := os.Args[1]
	minioPath := os.Args[2]
	// localFile := "./upload/upload.go"
	// minioPath := "/abc"

	// 上传错误信息
	// upload.UploadError(localFile, minioPath)

	// 视频转码
	video.Convert(localFile, minioPath)

	// 将转码后的上传到minio
	upload.ParallelUpload(minioPath, minioPath)
}
