package server

import (
	"encoding/json"
	"log"
	"net/http"
)

func _GET(next func(http.ResponseWriter, *http.Request)) func(http.ResponseWriter, *http.Request) {
	return _API(next, "GET")
}
func _POST(next func(http.ResponseWriter, *http.Request)) func(http.ResponseWriter, *http.Request) {
	return _API(next, "POST")
}
func _PUT(next func(http.ResponseWriter, *http.Request)) func(http.ResponseWriter, *http.Request) {
	return _API(next, "PUT")
}
func _DELETE(next func(http.ResponseWriter, *http.Request)) func(http.ResponseWriter, *http.Request) {
	return _API(next, "DELETE")
}

func _API(next func(http.ResponseWriter, *http.Request), method ...string) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		methodAllowed := false
		for _, v := range method {
			if r.Method == v {
				methodAllowed = true
				break
			}
		}
		if methodAllowed || len(method) == 0 {
			next(w, r)
		} else {
			w.WriteHeader(http.StatusMethodNotAllowed)
			_ERROR(w, 405, r.Method+"方法不允许")
		}
	}
}

type Response struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}

func _ERROR_DATA(w http.ResponseWriter, msg string, data interface{}) {
	// 创建一个Response实例
	resp := Response{
		Code:    500,
		Message: msg,
		Data:    data, // 可以是任意类型
	}
	out(w, resp)
}
func _ERROR(w http.ResponseWriter, code int, msg string) {
	// 创建一个Response实例
	resp := Response{
		Code:    code,
		Message: msg,
		Data:    nil, // 可以是任意类型
	}
	out(w, resp)
}

func _OK(w http.ResponseWriter, data interface{}) {
	// 创建一个Response实例
	resp := Response{
		Code:    200,
		Message: "Success",
		Data:    data, // 可以是任意类型
	}
	out(w, resp)
}

func out(w http.ResponseWriter, resp Response) {
	// 将Response实例序列化为JSON字符串
	jsonResp, err := json.Marshal(resp)
	if err != nil {
		log.Println("Error", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if resp.Code != 200 {
		http.Error(w, string(jsonResp), http.StatusInternalServerError)
		return
	}
	w.Write(jsonResp)
}
