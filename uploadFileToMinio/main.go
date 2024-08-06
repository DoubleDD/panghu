package main

import (
	"minioUploadFile/upload"
	"os"
)

func main() {
	if len(os.Args) == 1 {
		//fmt.Println("缺少参数")
		// return
	}
	// localFile := os.Args[1]
	// minioPath := os.Args[2]
	localFile := "./upload/upload.go"
	minioPath := "/abc"

	// 上传错误信息
	// upload.UploadError(localFile, minioPath)

	upload.MultipartUpload(localFile, minioPath)

}
