# 导入基础镜像golang:alpine
FROM golang:alpine AS builder

# 设置环境变量
ENV GO111MODULE=auto \
    CGO_ENABLED=0 \
    GOOS=linux \
    GOARCH=amd64 \
    GOPROXY="https://goproxy.cn,direct"
	
# 创建并移动到工作目录（可自定义路径）
WORKDIR /build

# 将代码复制到容器中
COPY . .

# 将代码编译成二进制可执行文件,文件名为 WebApp
RUN go build -o WebApp .

# 利用scratch创建一个小镜像
FROM scratch

WORKDIR /app
# 从builder镜像中把/WebApp 拷贝到当前目录
COPY --from=builder /build/WebApp /app/

# 将项目用到的静态文件拷贝到镜像（如果没有可忽略该步骤）
COPY ./config-k8s.toml /app/config.toml

# 声明服务端口
EXPOSE 7700

# 启动容器时运行的命令
CMD ["/app/WebApp","server"]
