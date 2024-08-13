package routehandle

import (
	"io"
	"log"
	"mime/multipart"
	"minioUploadFile/server/common"
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

// 文件上传，支持断点续传，主要通过 content-range 实现
func FileUpload(w http.ResponseWriter, r *http.Request) {
	contentRange := r.Header.Get("content-range")
	file, fileHeader, err := r.FormFile("file")
	if err != nil {
		log.Panic("获取文件出错", err)
		common.ERROR_DATA(w, "获取文件出错", err)
		return
	}
	defer file.Close()

	log.Printf("文件名：%v", fileHeader.Filename)
	log.Printf("文件大小：%v", fileHeader.Size)

	if strings.TrimSpace(contentRange) != "" {
		// TODO 断点续传部分
		log.Println("断点续传")
	} else {
		SaveFile(file, fileHeader)
	}
	common.OK(w, fileHeader.Filename)
}

func SaveFile(file multipart.File, fileHeader *multipart.FileHeader) (string, error) {
	// 创建目标文件路径
	uploadPath := os.Getenv("HOME") + "/tmp/go-upload"
	err := os.MkdirAll(uploadPath, os.ModePerm)
	if err != nil {
		log.Fatal("创建目录失败", err)
		return "", err
	}

	dstPath := filepath.Join(uploadPath, fileHeader.Filename)
	dst, err := os.Create(dstPath)
	if err != nil {
		log.Fatal("创建文件失败", err)
		return "", err
	}
	defer dst.Close()

	// 将上传的文件内容复制到目标文件
	_, err = io.Copy(dst, file)
	if err != nil {
		log.Fatal("copy文件失败", err)
		return "", err
	}
	return uploadPath + "/" + fileHeader.Filename, nil
}


