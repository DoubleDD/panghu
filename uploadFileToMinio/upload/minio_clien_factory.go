package upload

import (
	"context"
	"fmt"
	"log"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
	"github.com/spf13/viper"
)

var BucketName, Endpoint, accessKeyID, secretAccessKey string
var minioClient *minio.Client
var coreClient *minio.Core

func GetConfig()  {
	
}

func MinioInit() {
	log.Println("初始化minio配置")
	viper.SetConfigName("config")
	viper.AddConfigPath(".")
	err := viper.ReadInConfig()
	if err != nil {
		fmt.Println(err)
		return
	}
	BucketName = viper.GetString("minio.bucket")
	Endpoint = viper.GetString("minio.endpoint")
	accessKeyID = viper.GetString("minio.accessKeyID")
	secretAccessKey = viper.GetString("minio.secretAccessKey")

	log.Println("minio配置:")
	log.Println("\tbucket: \t", BucketName)
	log.Println("\tendpoint: \t", Endpoint)
	log.Println("\taccessKeyID: \t", accessKeyID)
	log.Println("\tsecretAccessKey:", secretAccessKey)

	// 初始化minioClient
	coreClient = minioCore(Endpoint, accessKeyID, secretAccessKey)
	minioClient = coreClient.Client

	bucketCheck(BucketName)
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
