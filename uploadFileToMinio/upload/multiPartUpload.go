package upload

import (
	"context"
	"fmt"
	"io"
	"log"
	"os"

	"github.com/minio/minio-go/v7"
)

func MultipartUpload(filePath, objectName string) {
	// 分片大小 (例如：5MB)
	chunkSize := int64(5 * 1024 * 1024)

	// 读取文件并分割
	file, err := os.Open(filePath)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()
	putOptions := createPutOptions()

	// 分片上传
	uploadID, err := coreClient.NewMultipartUpload(context.Background(), bucketName, objectName, putOptions)
	if err != nil {
		log.Fatal(err)
	}
	defer func() {
		if err != nil {
			coreClient.RemoveIncompleteUpload(context.Background(), bucketName, objectName)
		}
	}()

	// 创建一个通道来收集所有的 part 信息
	partsChan := make(chan minio.CompletePart, 100)
	var parts []minio.CompletePart

	// 开始 goroutines
	go func() {
		for part := range partsChan {
			parts = append(parts, part)
		}
	}()

	// 读取文件并上传分片
	uploadPart(file, chunkSize, bucketName, objectName, uploadID, partsChan)

	// 完成 multipart upload
	_, err = coreClient.CompleteMultipartUpload(context.Background(), bucketName, objectName, uploadID, parts, putOptions)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("File uploaded successfully.")
}

// 上传分片
func uploadPart(file *os.File, chunkSize int64, bucketName, objectName, uploadID string, parts chan<- minio.CompletePart) {
	partNumber := 1
	for {
		// 创建一个新的 reader 用于读取文件的一部分
		reader := io.LimitReader(file, chunkSize)

		// 如果已经到达文件末尾，则停止循环
		if _, ok := reader.(*io.LimitedReader); !ok {
			break
		}

		// 上传分片
		result, err := coreClient.PutObjectPart(context.Background(), bucketName, objectName, uploadID, partNumber, reader, chunkSize, createPutPartOptions())
		if err != nil {
			log.Fatal(err)
		}

		// 发送 part 到通道
		parts <- minio.CompletePart{ETag: result.ETag, PartNumber: partNumber}

		partNumber++
	}
}

func createPutOptions() minio.PutObjectOptions {
	return minio.PutObjectOptions{}
}
func createPutPartOptions() minio.PutObjectPartOptions {
	return minio.PutObjectPartOptions{}
}
