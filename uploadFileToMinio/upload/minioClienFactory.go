package upload

import (
	"context"
	"fmt"
	"log"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
	"github.com/spf13/viper"
)

var bucketName, endpoint, accessKeyID, secretAccessKey string
var minioClient *minio.Client
var coreClient *minio.Core

func init() {
	viper.SetConfigName("config")
	viper.AddConfigPath(".")
	err := viper.ReadInConfig()
	if err != nil {
		fmt.Println(err)
		return
	}
	bucketName = viper.GetString("bucket")
	endpoint = viper.GetString("endpoint")
	accessKeyID = viper.GetString("accessKeyID")
	secretAccessKey = viper.GetString("secretAccessKey")

	bucketCheck(bucketName)

	// 初始化minioClient
	coreClient = minioCore(endpoint, accessKeyID, secretAccessKey)
	minioClient = coreClient.Client
}

func minioCore(endpoint, accessKeyID, secretAccessKey string) *minio.Core {
	coreClient, err := minio.NewCore(endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(accessKeyID, secretAccessKey, ""),
		Secure: false,
	})
	if err != nil {
		log.Fatalln(err)
	}

	return coreClient
}

// 检查 bucket 是否存在
func bucketCheck(bucketName string) {
	found, err := minioClient.BucketExists(context.Background(), bucketName)
	if err != nil {
		log.Fatal(err)
	}
	if !found {
		log.Fatalf("Bucket '%s' not found", bucketName)
	}
}
