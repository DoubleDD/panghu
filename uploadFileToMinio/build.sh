#!/bin/bash

# 定义要编译的目标操作系统和架构
targets=(
    # "linux amd64"
    "linux arm64"
    # "darwin arm64"
)

# 循环遍历每个目标
for target in "${targets[@]}"; do
    # 解析目标
    IFS=' ' read -r os arch <<< "$target"
    
    # 设置环境变量
    export GOOS=$os
    export GOARCH=$arch
    
    # 构建输出文件名
    output="minioUpload-$os-$arch"
    
    # 编译
    go build -o "$output"
done

# 清理环境变量
unset GOOS
unset GOARCH
