package upload

import (
	"context"
	"log"
	"path/filepath"
	"time"

	"github.com/minio/minio-go/v7"
)

func UploadError(gc_path, hp_path string) {

	// 日期
	timeStr := time.Now().Format("2006-01-02_15:04:05")

	// 上传gc日志文件
	// gc_path = "/Users/kedong/code/com.gitee/uploadFileToMinio/test.log"
	objectName := filepath.Base(gc_path)
	uploadFile("jvm_oom/"+timeStr+"/gc.log/"+objectName, gc_path)

	// 上传heapdump文件
	uploadFile("jvm_oom/"+timeStr+"/heapdump.bin", hp_path)
}

func uploadFile(objectName string, filePath string) {
	ctx := context.Background()
	err := minioClient.MakeBucket(ctx, bucketName, minio.MakeBucketOptions{})
	if err != nil {
		// Check to see if we already own this bucket (which happens if you run this twice)
		exists, errBucketExists := minioClient.BucketExists(ctx, bucketName)
		if errBucketExists == nil && exists {
			log.Printf("We already own %s\n", bucketName)
		} else {
			log.Fatalln(err)
		}
	} else {
		log.Printf("Successfully created %s\n", bucketName)
	}

	contentType := "application/octet-stream"

	// Upload the test file with FPutObject
	info, err := minioClient.FPutObject(ctx, bucketName, objectName, filePath, minio.PutObjectOptions{ContentType: contentType})
	if err != nil {
		log.Fatalln(err)
	}

	log.Printf("Successfully uploaded %s of size %d\n", objectName, info.Size)
}
