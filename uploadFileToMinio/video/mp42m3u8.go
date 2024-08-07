package video

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strings"
)

/*
mp4 -> m3u8

通过ffmpeg 实现
*/
func Convert(inputFile, outputDir string) (string, string, error) {
	fileName := filepath.Base(inputFile)
	dotIndex := strings.LastIndex(fileName, ".")
	if dotIndex > 0 {
		fileName = fileName[:dotIndex]
	}
	targetDir := outputDir + "/" + fileName
	// 创建输出目录如果不存在
	err := os.MkdirAll(targetDir, os.ModePerm)
	if err != nil {
		log.Fatalf("Error creating output directory: %s", err)
	}
	reg := regexp.MustCompile(`/{2,}`)
	// 将输出写入文件
	m3u8FileName := "index.m3u8"
	m3u8File := reg.ReplaceAllString(targetDir+"/"+m3u8FileName, "/")
	segmentFile := reg.ReplaceAllString(targetDir+"/index%05d.ts", "/")

	// 重定向标准输出到文件
	run(exec.Command("ffmpeg",
		"-i", inputFile, // 从标准输入读取
		"-c:v", "h264",
		"-flags", "+cgop",
		"-g", "30",
		"-map", "0",
		"-f", "hls",
		"-hls_time", "10",
		"-hls_list_size", "0",
		"-hls_segment_filename", segmentFile,
		m3u8File,
	))

	run(exec.Command("ls", "-lh", outputDir))

	fmt.Println("Conversion completed successfully.")
	fmt.Println("\n", m3u8File)
	return m3u8FileName, targetDir, nil
}

func run(cmd *exec.Cmd) {
	// 打印将要执行的命令
	fmt.Printf("执行命令: \n%s\n", cmd.String())

	// 将命令的标准输出和标准错误输出连接到标准输出和标准错误输出
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	// 启动命令
	if err := cmd.Start(); err != nil {
		fmt.Printf("Error starting command: %v\n", err)
		return
	}

	// 等待命令完成
	if err := cmd.Wait(); err != nil {
		fmt.Printf("Command finished with error: %v\n", err)
	}

}
