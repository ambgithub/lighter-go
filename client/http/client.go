package http

import (
	"crypto/tls"
	"net"
	"net/http"
	"net/url"
	"time"

	core "github.com/ambgithub/lighter-go/client"
)

var (
	dialer = &net.Dialer{
		Timeout:   10 * time.Second,
		KeepAlive: 60 * time.Second,
	}
	transport = &http.Transport{
		DialContext:         dialer.DialContext,
		MaxConnsPerHost:     1000,
		MaxIdleConnsPerHost: 100,
		IdleConnTimeout:     10 * time.Second,
		TLSClientConfig:     &tls.Config{InsecureSkipVerify: false},
	}

	httpClient = &http.Client{
		Timeout:   time.Second * 30,
		Transport: transport,
	}
)

var _ core.MinimalHTTPClient = (*client)(nil)

type client struct {
	endpoint string
	// === 新增以下字段 ===
	proxyURL      *url.URL          // 代理地址
	localAddr     net.Addr          // 本地网卡地址 (用于多IP服务器)
	customHeaders map[string]string // 自定义请求头
	httpClient    *http.Client      // 内部使用的 http client，我们需要深度定制它
}

func NewClient(baseUrl string) core.MinimalHTTPClient {
	if baseUrl == "" {
		return nil
	}

	return &client{
		endpoint: baseUrl,
	}
}
