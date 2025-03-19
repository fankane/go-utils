package http

import (
	"context"
	"fmt"
	"net"
	"net/http"
	"net/http/httputil"
	"net/url"
)

// ProxyRequest 转发请求到指定路径，并将响应写回客户端
func ProxyRequest(ctx context.Context, targetURL string, w http.ResponseWriter, r *http.Request) error {
	target, err := url.Parse(targetURL)
	if err != nil {
		return fmt.Errorf("parse target url failed:%s", err)
	}
	// 保留原始查询参数
	target.RawQuery = r.URL.RawQuery
	// 创建一个反向代理
	proxy := httputil.NewSingleHostReverseProxy(target)
	// 自定义 Director 函数
	proxy.Director = func(req *http.Request) {
		// 修改请求的目标地址
		req.URL.Scheme = target.Scheme
		req.URL.Host = target.Host
		req.URL.Path = target.Path

		// 保留源请求的 Body
		if r.Body != nil {
			req.Body = r.Body
			req.ContentLength = r.ContentLength
		}

		// 复制源请求的 Header
		req.Header = make(http.Header)
		for key, values := range r.Header {
			for _, value := range values {
				req.Header.Add(key, value)
			}
		}

		// 设置 X-Forwarded-For 头
		if clientIP, _, err := net.SplitHostPort(r.RemoteAddr); err == nil {
			req.Header.Set("X-Forwarded-For", clientIP)
		}
	}
	// 转发请求并将响应写回客户端
	proxy.ServeHTTP(w, r)
	return nil
}
