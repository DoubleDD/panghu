package upload

import (
	"context"
	"fmt"
	"io"
	"log"
	"os"
	"regexp"
	"sort"
	"sync"
	"time"

	"github.com/minio/minio-go/v7"
)

func MultipartUpload(filePath, prefixPath string) {
	// 开始计时
	timeStart := time.Now()

	// 读取文件并分割
	file, err := os.Open(filePath)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	// 文件信息
	fileInfo, err := file.Stat()
	if err != nil {
		fmt.Println("Error getting file info:", err)
		return
	}
	if fileInfo.IsDir() {
		// 遍历目录
		fileInfos, err := file.Readdir(-1)
		if err != nil {
			fmt.Println("Error getting file info:", err)
			return
		}
		// 遍历目录内容并打印文件名
		for _, fileInfo := range fileInfos {
			fileName := normalizeSlashes(filePath + "/" + fileInfo.Name())
			objectName := normalizeSlashes(prefixPath + "/" + fileInfo.Name())
			MultipartUpload(fileName, objectName)
		}
	} else {
		// 上传文件,采用分片上传
		uploadFileWithM(filePath, normalizeSlashes(prefixPath+"/"+fileInfo.Name()))
	}

	// 计时结束
	timeEnd := time.Now()
	// 计算时间差
	elapsed := timeEnd.Sub(timeStart)
	fmt.Printf("文件上传总耗时 %vms\n", elapsed.Milliseconds())
	fmt.Println("File uploaded successfully.")

}
func uploadFileWithM(fileName, objectName string) {
	// 分片大小 (例如：5MB)
	chunkSize := int64(5 * 1024 * 1024)

	// 读取文件并分割
	file, err := os.Open(fileName)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()
	// 文件信息
	fileInfo, err := file.Stat()
	if err != nil {
		fmt.Println("Error getting file info:", err)
		return
	}

	// 文件大小
	fileSize := fileInfo.Size()
	log.Println(fileInfo.Name()+"\t", fileSize)

	// 根据分片大小计算需要几个分片, 需要向上取整
	chunkNumber := (fileSize + chunkSize - 1) / chunkSize
	log.Printf("上传文件 %v  分配 %d 个任务", fileName, chunkNumber)
	// 分片上传可选参数
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

	// 创建一个通道来收集所有的 part(分片) 信息
	partsChan := make(chan minio.CompletePart, 100)
	var parts []minio.CompletePart
	var wg sync.WaitGroup
	wg.Add(int(chunkNumber))
	// 开始 goroutines
	go func() {
		for part := range partsChan {
			log.Println("任务完成", part)
			parts = append(parts, part)
			wg.Done()
		}
	}()

	// 每个分片都分配一个 goroutine
	for i := int64(0); i < chunkNumber; i++ {
		// 上传分片, 开始位置，结束位置
		start := i * chunkSize
		end := (i + 1) * chunkSize
		if start > fileSize {
			break
		}
		if end > fileSize {
			end = fileSize - 1
		}
		go uploadPart(file, i+1, start, end, bucketName, objectName, uploadID, partsChan)
	}

	wg.Wait()

	// parts 排序
	sort.Sort(ByPartNumber(parts))

	// 完成 multipart upload
	_, err = coreClient.CompleteMultipartUpload(context.Background(), bucketName, objectName, uploadID, parts, putOptions)
	if err != nil {
		log.Fatal(err)
	}
}

// 上传分片
func uploadPart(file *os.File, partNumber, start, end int64, bucketName, objectName, uploadID string, parts chan<- minio.CompletePart) {
	log.Printf("任务-%d，start: %d, end: %d", partNumber, start, end)
	// 创建一个新的 reader 用于读取文件的一部分
	sectionReader := io.NewSectionReader(file, start, end-start+1)

	// 上传分片
	result, err := coreClient.PutObjectPart(context.Background(), bucketName, objectName, uploadID, int(partNumber), sectionReader, end-start+1, createPutPartOptions())
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("任务-%d，完成。发送消息", partNumber)
	// 发送 part 到通道
	parts <- minio.CompletePart{ETag: result.ETag, PartNumber: int(partNumber)}
}

func createPutOptions() minio.PutObjectOptions {
	return minio.PutObjectOptions{}
}
func createPutPartOptions() minio.PutObjectPartOptions {
	return minio.PutObjectPartOptions{}
}
func normalizeSlashes(str string) string {
	re := regexp.MustCompile(`/{2,}`)
	return re.ReplaceAllString(str, "/")
}

// ByPartNumber 是一个自定义类型，用于按照 PartNumber 排序
type ByPartNumber []minio.CompletePart

// Len 实现 sort.Interface 的 Len 方法
func (a ByPartNumber) Len() int { return len(a) }

// Less 实现 sort.Interface 的 Less 方法
func (a ByPartNumber) Less(i, j int) bool { return a[i].PartNumber < a[j].PartNumber }

// Swap 实现 sort.Interface 的 Swap 方法
func (a ByPartNumber) Swap(i, j int) { a[i], a[j] = a[j], a[i] }
