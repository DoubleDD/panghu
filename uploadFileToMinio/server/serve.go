package server

import (
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

func Start(port int) {
	mux := &http.ServeMux{}
	initHandler(mux)
	log.Println("Starting server at port", port)
	if err := http.ListenAndServe(":"+strconv.Itoa(port), mux); err != nil {
		log.Printf("Server failed to start: %v\n", err)
	}
}

func initHandler(mux *http.ServeMux) {
	mux.HandleFunc("/video/2m3u8", mp42m3u8)
	mux.HandleFunc("/file/upload", _API(fileUpload, "POST", "PUT"))
}

// 将视频转成m3u8格式，并上传至minio
func mp42m3u8(w http.ResponseWriter, r *http.Request) {
	log.Println("接口测试OK")
	log.Panic("Panic 测试")
	_OK(w, "接口测试OK")
}

// 文件上传，支持断点续传，主要通过 content-range 实现
func fileUpload(w http.ResponseWriter, r *http.Request) {
	contentRange := r.Header.Get("content-range")
	file, fileHeader, err := r.FormFile("file")
	if err != nil {
		log.Panic("获取文件出错", err)
		_ERROR_DATA(w, "获取文件出错", err)
		return
	}
	defer file.Close()

	log.Printf("文件名：%v", fileHeader.Filename)
	log.Printf("文件大小：%v", fileHeader.Size)

	if strings.TrimSpace(contentRange) != "" {
		// 断点续传部分
		log.Println("断点续传")
	} else {
		// 创建目标文件路径
		uploadPath := "/Users/kedong/tmp/go-upload"
		err := os.MkdirAll(uploadPath, os.ModePerm)
		if err != nil {
			_ERROR_DATA(w, "创建目录失败", err)
			return
		}

		dstPath := filepath.Join(uploadPath, fileHeader.Filename)
		dst, err := os.Create(dstPath)
		if err != nil {
			_ERROR_DATA(w, "创建文件失败", err)
			return
		}
		defer dst.Close()

		// 将上传的文件内容复制到目标文件
		_, err = io.Copy(dst, file)
		if err != nil {
			_ERROR_DATA(w, "copy文件失败", err)
			return
		}
	}
	_OK(w, fileHeader.Filename)
}
