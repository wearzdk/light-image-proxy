FROM golang:alpine AS builder

WORKDIR /app
COPY . .

# 构建可执行文件
RUN go mod tidy && \
    CGO_ENABLED=0 GOOS=linux go build -a -ldflags="-s -w" -o light-image-proxy .

FROM alpine:latest
RUN apk add --no-cache ca-certificates

# 从builder阶段复制可执行文件
COPY --from=builder /app/light-image-proxy /light-image-proxy

# 暴露端口
EXPOSE 16524

# 运行程序
ENTRYPOINT ["/light-image-proxy"] 