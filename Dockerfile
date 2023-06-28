# 设置基础镜像
FROM golang:latest

# 设置工作目录
WORKDIR /app

# 将项目文件复制到容器中
COPY . .

# 下载依赖并编译项目
RUN go mod download
RUN go build -o main .

# 设置容器启动命令
CMD ["./main"]
