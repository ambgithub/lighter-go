package client

import (
	"crypto/tls"
	"fmt"
	"net"
	"net/http"
	"net/url"
	"time"
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
		TLSClientConfig:     &tls.Config{InsecureSkipVerify: true},
	}

	httpClient = &http.Client{
		Timeout:   time.Second * 30,
		Transport: transport,
	}
)

// HTTPClientOption 是一个函数类型，用于以非侵入的方式配置HTTP客户端。
// 这种"Option模式"是Go中非常推荐的一种设计，让配置变得非常灵
type HTTPClientOption func(*HTTPClient)

// HTTPClient 是一个经过增强的、支持高级网络配置的API客户端。
type HTTPClient struct {
	Client              *http.Client // 核心http客户端，现在是实例的一部分，而不是全局变量
	endpoint            string       // API的基础URL
	defaultHeaders      http.Header  // 存储所有请求都将携带的默认Header
	channelName         string       // 旧版库中的字段，予以保留
	fatFingerProtection bool         // 旧版库中的字段，予以保留
}

// NewHTTPClient 创建一个新的、可配置的HTTPClient实例。
// 它接受一个基础URL和一系列可选的配置函数（Options）。
func NewHTTPClient(baseUrl string, options ...HTTPClientOption) (*HTTPClient, error) {
	if baseUrl == "" {
		return nil, fmt.Errorf("baseUrl不能为空")
	}

	// 初始化客户端，并设置默认值
	c := &HTTPClient{
		endpoint:            baseUrl,
		defaultHeaders:      make(http.Header),
		channelName:         "",
		fatFingerProtection: true,
	}

	// 应用所有传入的配置选项
	for _, option := range options {
		option(c)
	}

	// 如果在应用选项后，Client仍未被自定义，则创建一个默认的
	if c.Client == nil {
		c.Client = &http.Client{
			Timeout: time.Second * 30,
			Transport: &http.Transport{
				Proxy: http.ProxyFromEnvironment, // 默认支持系统环境变量中的代理
				DialContext: (&net.Dialer{
					Timeout:   10 * time.Second,
					KeepAlive: 60 * time.Second,
				}).DialContext,
				MaxConnsPerHost:     1000,
				MaxIdleConnsPerHost: 100,
				IdleConnTimeout:     10 * time.Second,
				TLSClientConfig:     &tls.Config{InsecureSkipVerify: true}, // 注意：跳过TLS验证，与原版保持一致
			},
		}
	}

	return c, nil
}

// WithProxy 是一个配置选项，用于为客户端设置HTTP或SOCKS5代理。
func WithProxy(proxyURL string) HTTPClientOption {
	return func(c *HTTPClient) {
		if proxyURL == "" {
			return
		}
		p, err := url.Parse(proxyURL)
		if err != nil {
			// 在实际应用中，您可能希望返回一个错误而不是panic
			panic(fmt.Sprintf("无效的代理URL: %v", err))
		}

		// 获取或创建Transport，并设置代理
		NewTransport, _ := c.getOrCreateTransport()
		NewTransport.Proxy = http.ProxyURL(p)
	}
}

// WithCustomHeaders 是一个配置选项，用于在客户端初始化时设置默认的Header。
func WithCustomHeaders(headers map[string]string) HTTPClientOption {
	return func(c *HTTPClient) {
		for key, value := range headers {
			c.defaultHeaders.Set(key, value)
		}
	}
}

// WithLocalAddr 是一个配置选项，用于指定客户端发出请求时绑定的本地出口IP地址。
func WithLocalAddr(localIP string) HTTPClientOption {
	return func(c *HTTPClient) {
		if localIP == "" {
			return
		}
		addr, err := net.ResolveTCPAddr("tcp", localIP+":0")
		if err != nil {
			panic(fmt.Sprintf("无效的本地IP地址: %v", err))
		}
		newDialer := &net.Dialer{
			Timeout:   10 * time.Second,
			KeepAlive: 60 * time.Second,
			LocalAddr: addr, // 在这里绑定本地地址
		}

		newTransport, _ := c.getOrCreateTransport()
		newTransport.DialContext = newDialer.DialContext
	}
}

// getOrCreateTransport 是一个辅助方法，用于安全地获取或创建一个http.Transport实例。
func (c *HTTPClient) getOrCreateTransport() (*http.Transport, bool) {
	if c.Client == nil {
		c.Client = &http.Client{
			Transport: &http.Transport{},
		}
	}
	newTransport, ok := c.Client.Transport.(*http.Transport)
	return newTransport, ok
}

func (c *HTTPClient) SetFatFingerProtection(enabled bool) {
	c.fatFingerProtection = enabled
}
