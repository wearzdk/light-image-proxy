package main

import (
	"context"
	"flag"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"
)

var (
	port        = flag.Int("port", 16524, "服务端口")
	timeout     = flag.Int("timeout", 10, "请求超时时间(秒)")
	userAgent   = flag.String("ua", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/123.0.0.0 Safari/537.36", "默认User-Agent")
	logRequests = flag.Bool("log", false, "是否记录请求日志")
)

// 代理处理函数
func proxyHandler(w http.ResponseWriter, r *http.Request) {
	startTime := time.Now()

	// 获取URL参数
	targetURL := r.URL.Query().Get("url")
	if targetURL == "" {
		http.Error(w, "缺少url参数", http.StatusBadRequest)
		return
	}

	// 创建HTTP客户端
	client := &http.Client{
		Timeout: time.Duration(*timeout) * time.Second,
	}

	// 创建请求
	req, err := http.NewRequest(r.Method, targetURL, r.Body)
	if err != nil {
		http.Error(w, fmt.Sprintf("创建请求失败: %v", err), http.StatusInternalServerError)
		return
	}

	// 设置默认User-Agent
	req.Header.Set("User-Agent", *userAgent)

	// 传递部分原始请求的header
	if r.Header.Get("Referer") != "" {
		req.Header.Set("Referer", r.Header.Get("Referer"))
	}

	// 处理自定义headers (格式: key1:value1,key2:value2)
	customHeaders := r.URL.Query().Get("headers")
	if customHeaders != "" {
		headerPairs := strings.Split(customHeaders, ",")
		for _, pair := range headerPairs {
			parts := strings.SplitN(pair, ":", 2)
			if len(parts) == 2 {
				key := strings.TrimSpace(parts[0])
				value := strings.TrimSpace(parts[1])
				if key != "" && value != "" {
					req.Header.Set(key, value)
				}
			}
		}
	}

	// 执行请求
	resp, err := client.Do(req)
	if err != nil {
		http.Error(w, fmt.Sprintf("请求失败: %v", err), http.StatusInternalServerError)
		return
	}
	defer resp.Body.Close()

	// 复制响应头
	for key, values := range resp.Header {
		for _, value := range values {
			w.Header().Add(key, value)
		}
	}

	// 设置响应状态码
	w.WriteHeader(resp.StatusCode)

	// 传输响应体
	bytesCopied, err := io.Copy(w, resp.Body)
	if err != nil {
		log.Printf("传输响应失败: %v", err)
		return
	}

	// 记录请求日志
	if *logRequests {
		duration := time.Since(startTime)
		log.Printf("请求: %s [%s] 状态: %d 大小: %.2f KB 耗时: %.2f ms",
			targetURL, r.Method, resp.StatusCode, float64(bytesCopied)/1024, float64(duration.Microseconds())/1000)
	}
}

// 简化的代理处理函数 - 支持URL路径形式请求
func simplifiedProxyHandler(w http.ResponseWriter, r *http.Request) {
	// 从路径中提取URL
	urlPath := r.URL.Path[len("/"):]
	if urlPath == "" {
		http.Error(w, "无效的URL路径", http.StatusBadRequest)
		return
	}

	// 确保URL有协议前缀
	if !strings.HasPrefix(urlPath, "http://") && !strings.HasPrefix(urlPath, "https://") {
		urlPath = "https://" + urlPath
	}

	// 设置url参数并转发到标准处理函数
	q := r.URL.Query()
	q.Set("url", urlPath)
	r.URL.RawQuery = q.Encode()

	proxyHandler(w, r)
}

func main() {
	flag.Parse()

	// 注册标准代理路由
	http.HandleFunc("/get", proxyHandler)

	// 注册简化路由 - 支持直接在路径中指定URL
	http.HandleFunc("/http://", simplifiedProxyHandler)
	http.HandleFunc("/https://", simplifiedProxyHandler)

	// 添加健康检查路由
	http.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("ok"))
	})

	// 添加首页简单说明
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		// 如果不是根路径，尝试作为简化形式的代理
		if r.URL.Path != "/" {
			simplifiedProxyHandler(w, r)
			return
		}

		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		w.Write([]byte(`
		<html>
			<head>
				<title>轻量图像代理服务</title>
				<style>
					body { font-family: Arial, sans-serif; max-width: 800px; margin: 0 auto; padding: 20px; }
					code { background-color: #f5f5f5; padding: 2px 5px; border-radius: 3px; }
					pre { background-color: #f5f5f5; padding: 10px; border-radius: 5px; overflow-x: auto; }
					.example { margin: 10px 0; }
				</style>
			</head>
			<body>
				<h1>轻量图像代理服务</h1>
				
				<h2>基本用法</h2>
				<div class="example">
					<p>1. 标准格式: <code>/get?url=图片地址</code></p>
					<p>示例: <a href="/get?url=https://example.com/image.jpg" target="_blank">/get?url=https://example.com/image.jpg</a></p>
				</div>
				
				<div class="example">
					<p>2. 简化格式: <code>/https://example.com/image.jpg</code></p>
					<p>示例: <a href="/https://example.com/image.jpg" target="_blank">/https://example.com/image.jpg</a></p>
				</div>
				
				<h2>自定义请求头</h2>
				<div class="example">
					<p>格式: <code>/get?url=图片地址&headers=Header1:Value1,Header2:Value2</code></p>
					<p>示例: <a href="/get?url=https://example.com/image.jpg&headers=Referer:https://example.com,Cache-Control:no-cache" target="_blank">/get?url=...&headers=Referer:https://example.com,Cache-Control:no-cache</a></p>
				</div>
				
				<h2>支持的请求方法</h2>
				<p>支持GET、POST等多种HTTP请求方法，可用于不同图片加载场景</p>
			</body>
		</html>
		`))
	})

	// 优雅关闭服务
	serverAddr := fmt.Sprintf(":%d", *port)
	server := &http.Server{
		Addr: serverAddr,
	}

	// 启动服务
	go func() {
		log.Printf("图像代理服务已启动 http://localhost:%d", *port)
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("服务启动失败: %v", err)
		}
	}()

	// 等待中断信号
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("正在关闭服务...")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		log.Fatalf("服务关闭失败: %v", err)
	}

	log.Println("服务已安全关闭")
}
