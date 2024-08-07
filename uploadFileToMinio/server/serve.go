package server

import (
	"log"
	"net/http"
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
	if strings.TrimSpace(contentRange) != "" {
		// 断点续传部分
		log.Println("")

	} else {
		// 非断点续传
		log.Println("")
	}
	_OK(w, nil)
}
