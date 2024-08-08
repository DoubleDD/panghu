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
	// 上传完后删除临时文件
	defer func() {
		log.Println("删除视频临时文件", filePath)
		os.Remove(filePath)
	}()

	hlsDir := os.Getenv("HOME") + "/tmp/hls"

	// 视频转码
	m3u8File, m3u8Dir, err := video.Convert(filePath, hlsDir)
	if err != nil {
		log.Fatal("视频转码失败", err)
		common.ERROR_DATA(w, "视频转码失败", err)
		return
	}
	// 上传完后删除临时目录
	defer func() {
		log.Println("删除 m3u8 临时文件夹", m3u8Dir)
		os.RemoveAll(m3u8Dir)
	}()

	// 将m3u8 文件上传至minio
	upload.MultipartUpload(m3u8Dir, m3u8Dir)

	// minioPath
	url := "/" + upload.BucketName + "/" + m3u8File

	common.OK(w, url)
}
