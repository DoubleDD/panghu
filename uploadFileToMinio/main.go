package main

import (
	"fmt"
	"minioUploadFile/upload"
	"os"
)

func main() {
	if len(os.Args) == 1 {
		fmt.Println("缺少参数")
		return
	}
	localFile := os.Args[1]
	minioPath := os.Args[2]

	// 上传错误信息
	// upload.UploadError(localFile, minioPath)

	upload.MultipartUpload(localFile, minioPath)

}
