package server

import (
	"log"
	"minioUploadFile/server/common"
	rh "minioUploadFile/server/route_handle"
	"minioUploadFile/upload"
	"net/http"
	"strconv"
)

func Start(port int) {
	mux := &http.ServeMux{}
	initHandler(mux)

	// 初始化其他东西
	serverInit()

	log.Println("Starting server at port", port)
	if err := http.ListenAndServe(":"+strconv.Itoa(port), mux); err != nil {
		log.Printf("Server failed to start: %v\n", err)
	}
}

func initHandler(mux *http.ServeMux) {
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		common.OK(w, "ok")
	})
	mux.HandleFunc("/k8s/pod/info", rh.K8sPodInfo)
	mux.HandleFunc("/video/2m3u8", rh.Mp42m3u8)
	mux.HandleFunc("/file/upload", common.API(rh.FileUpload, "POST", "PUT"))
	mux.HandleFunc("/minio/upload", common.POST(rh.MinioUpload))
}

func serverInit() {
	upload.MinioInit()
}
