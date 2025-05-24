# Light Image Proxy

A super lightweight image proxy service developed in Go. Designed to solve image loading issues caused by hotlink protection, especially useful when rendering third-party image resources locally.

[中文](./README.md)

## Features

- 🚀 **Ultra-lightweight**: Low memory footprint, small image size
- 🔄 **Pure forwarding**: No caching, direct request forwarding to source
- 🛠️ **Customizable**: Support for custom request headers to handle various hotlink protection strategies
- 🐳 **Containerized**: Dockerfile provided, supports Docker deployment
- 🌐 **Versatile**: Suitable for various image loading scenarios

## Quick Start

### Docker Deployment (Recommended)

```bash
# Run container
docker run -d -p 16524:16524 --name light-image-proxy wearzdk/light-image-proxy

# With custom parameters
docker run -d -p 8080:8080 --name light-image-proxy wearzdk/light-image-proxy -port=8080
```

### Direct Execution

```bash
# Start with default configuration (port 16524)
go run main.go

# Custom port
go run main.go -port=8080

# Custom timeout (seconds)
go run main.go -timeout=30

# Custom User-Agent
go run main.go -ua="Custom User Agent"
```

## Usage

### Basic Request

```
http://localhost:16524/get?url=https://example.com/image.jpg
```

### Custom Headers

```
http://localhost:16524/get?url=https://example.com/image.jpg&headers=Referer:https://example.com,Cache-Control:no-cache
```

Header format: `key1:value1,key2:value2`

## Use Cases

- Loading third-party images restricted by hotlink protection when generating videos with tools like Remotion
- Solving CORS or hotlink issues during frontend development
- Various scenarios requiring access to restricted image resources

## Configuration Parameters

| Parameter | Default    | Description               |
| --------- | ---------- | ------------------------- |
| port      | 16524      | Service listening port    |
| timeout   | 10         | Request timeout (seconds) |
| ua        | Mozilla... | Default User-Agent        |

## Building the Project

```bash
# Clone repository
git clone https://github.com/your-username/light-image-proxy.git
cd light-image-proxy

# Build
go build -o light-image-proxy

# Run
./light-image-proxy
```

## License

MIT
