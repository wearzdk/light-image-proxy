# Light Image Proxy

一个超轻量级的图像代理服务，使用 Go 语言开发。用于解决由于防盗链导致的图片无法加载问题，特别适用于需要在本地渲染第三方图片资源的场景。

## 功能特点

- 🚀 **超轻量**：内存占用低，启动速度快
- 🔄 **纯转发**：无缓存，直接转发请求到源站
- 🛠️ **可定制**：支持自定义请求头，灵活应对不同防盗链策略
- 🐳 **容器化**：提供 Dockerfile，支持 Docker 部署
- 🌐 **通用性**：适用于各种图片加载场景

## 快速开始

### 直接运行

```bash
# 默认配置启动（端口16524）
go run main.go

# 自定义端口
go run main.go -port=8080

# 自定义超时时间（秒）
go run main.go -timeout=30

# 自定义User-Agent
go run main.go -ua="Custom User Agent"
```

### Docker 运行

```bash
# 运行容器
docker run -d -p 16524:16524 --name light-image-proxy wearzdk/light-image-proxy

# 自定义参数
docker run -d -p 8080:8080 --name light-image-proxy wearzdk/light-image-proxy -port=8080
```

## 使用方法

### 基本请求

```
http://localhost:16524/get?url=https://example.com/image.jpg
```

### 自定义请求头

```
http://localhost:16524/get?url=https://example.com/image.jpg&headers=Referer:https://example.com,Cache-Control:no-cache
```

请求头格式：`key1:value1,key2:value2`

## 应用场景

- 在使用 Remotion 等工具生成视频时，加载受防盗链限制的第三方图片
- 前端开发过程中，解决跨域或防盗链问题
- 各种需要访问受限制图片资源的场景

## 配置参数

| 参数    | 默认值     | 说明               |
| ------- | ---------- | ------------------ |
| port    | 16524      | 服务监听端口       |
| timeout | 10         | 请求超时时间（秒） |
| ua      | Mozilla... | 默认 User-Agent    |

## 构建项目

```bash
# 克隆仓库
git clone https://github.com/your-username/light-image-proxy.git
cd light-image-proxy

# 构建
go build -o light-image-proxy

# 运行
./light-image-proxy
```

## 许可证

MIT
