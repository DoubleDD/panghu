package server

import (
	"log"
	"minioUploadFile/server/common"
	routehandle "minioUploadFile/server/route_handle"
	"net/http"
	"strconv"
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
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		common.OK(w, "ok")
	})
	mux.HandleFunc("/video/2m3u8", routehandle.Mp42m3u8)
	mux.HandleFunc("/file/upload", common.API(routehandle.FileUpload, "POST", "PUT"))
}
