package routehandle

import (
	"log"
	"minioUploadFile/server/common"
	"minioUploadFile/upload"
	"minioUploadFile/video"
	"net/http"
	"os"
)

// 将视频转成m3u8格式，并上传至minio
func Mp42m3u8(w http.ResponseWriter, r *http.Request) {
	file, handle, err := r.FormFile("file")
	if err != nil {
		log.Fatal("获取文件失败", err)
		common.ERROR_DATA(w, "获取文件失败", err)
		return
	}

	// 先将文件保存至本地
	filePath, err := SaveFile(file, handle)
	if err != nil {
		log.Fatal("保存文件失败", err)
		common.ERROR_DATA(w, "保存文件失败", err)
		return
	}

	hlsDir := os.Getenv("HOME") + "/tmp/hls"

	// 视频转码
	m3u8File, dir, err := video.Convert(filePath, hlsDir)
	if err != nil {
		log.Fatal("视频转码失败", err)
		common.ERROR_DATA(w, "视频转码失败", err)
		return
	}

	// 将m3u8 文件上传至minio
	upload.MultipartUpload(dir, dir)

	// minioPath
	url := "/" + upload.BucketName + "/" + m3u8File

	common.OK(w, url)
}
