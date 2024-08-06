package main

import (
	"fmt"
	"net/http"
)

func main() {
	mux := &http.ServeMux{}
	mux.HandleFunc("/auth/", auth)
	fmt.Println("Starting server at port 8080")
	if err := http.ListenAndServe(":8080", mux); err != nil {
		fmt.Printf("Server failed to start: %v\n", err)
	}
}

func auth(w http.ResponseWriter, request *http.Request) {
	// 获取URL路径
	url := request.URL.Path
	// file url
	fileUrl := url[len("/auth"):]
	fmt.Println("\nurl: ", fileUrl)

	// 获取参数
	query := request.URL.Query()
	fmt.Println("查询参数: ", query)
	token := query["token"]
	if token == nil || token[0] != "123" {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	w.Header().Set("X-Accel-Redirect", "/oss"+fileUrl)
	w.WriteHeader(http.StatusOK)
}
