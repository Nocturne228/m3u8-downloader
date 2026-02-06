package http

import (
	"fmt"
	"net"
	"time"

	"m3u8-downloader/internal/errors"
	"m3u8-downloader/internal/logger"

	"github.com/levigross/grequests"
)

// Client HTTP 客户端接口
type Client interface {
	Get(url string) ([]byte, error)
	GetWithHeaders(url string, headers map[string]string) ([]byte, error)
	GetWithCookie(url string, cookie string) ([]byte, error)
}

// HTTPClient HTTP 客户端实现
type HTTPClient struct {
	timeout        time.Duration
	maxRetries     int
	userAgent      string
	insecureVerify bool
	logger         logger.Logger
}

// NewClient 创建新的 HTTP 客户端
func NewClient(timeout time.Duration, maxRetries int, userAgent string, logger logger.Logger) Client {
	return &HTTPClient{
		timeout:    timeout,
		maxRetries: maxRetries,
		userAgent:  userAgent,
		logger:     logger,
	}
}

// SetInsecureVerify 设置是否跳过 SSL 验证
func (c *HTTPClient) SetInsecureVerify(insecure bool) {
	c.insecureVerify = insecure
}

// Get 获取 URL 内容
func (c *HTTPClient) Get(url string) ([]byte, error) {
	return c.getWithOptions(url, map[string]string{
		"User-Agent":      c.userAgent,
		"Connection":      "keep-alive",
		"Accept":          "*/*",
		"Accept-Encoding": "*",
		"Accept-Language": "zh-CN,zh;q=0.9, en;q=0.8, de;q=0.7, *;q=0.5",
	})
}

// GetWithHeaders 使用自定义请求头获取 URL 内容
func (c *HTTPClient) GetWithHeaders(url string, headers map[string]string) ([]byte, error) {
	// 合并默认请求头
	finalHeaders := map[string]string{
		"User-Agent":      c.userAgent,
		"Connection":      "keep-alive",
		"Accept":          "*/*",
		"Accept-Encoding": "*",
		"Accept-Language": "zh-CN,zh;q=0.9, en;q=0.8, de;q=0.7, *;q=0.5",
	}

	for k, v := range headers {
		finalHeaders[k] = v
	}

	return c.getWithOptions(url, finalHeaders)
}

// GetWithCookie 使用自定义 Cookie 获取 URL 内容
func (c *HTTPClient) GetWithCookie(url string, cookie string) ([]byte, error) {
	headers := map[string]string{}
	if cookie != "" {
		headers["Cookie"] = cookie
	}
	return c.GetWithHeaders(url, headers)
}

func (c *HTTPClient) getWithOptions(url string, headers map[string]string) ([]byte, error) {
	var lastErr error

	for attempt := 1; attempt <= c.maxRetries; attempt++ {
		ro := &grequests.RequestOptions{
			UserAgent:      c.userAgent,
			RequestTimeout: c.timeout,
			Headers:        headers,
		}

		if c.insecureVerify {
			ro.InsecureSkipVerify = true
		}

		resp, err := grequests.Get(url, ro)
		if err != nil {
			lastErr = err
			c.logger.Warn("HTTP 请求失败 (尝试 %d/%d): %v", attempt, c.maxRetries, err)

			// 如果是网络错误，可能需要重试
			if isNetworkError(err) && attempt < c.maxRetries {
				backoffDuration := time.Duration((attempt - 1)) * time.Second
				time.Sleep(backoffDuration)
				continue
			}
		}

		if resp == nil || !resp.Ok {
			statusCode := 0
			if resp != nil {
				statusCode = resp.StatusCode
			}
			lastErr = fmt.Errorf("HTTP 状态异常: %d", statusCode)
			c.logger.Warn("HTTP 响应失败 (尝试 %d/%d): 状态码 %d", attempt, c.maxRetries, statusCode)

			if attempt < c.maxRetries {
				backoffDuration := time.Duration((attempt - 1)) * time.Second
				time.Sleep(backoffDuration)
				continue
			}
		}

		if resp != nil && resp.Ok {
			return resp.Bytes(), nil
		}
	}

	return nil, errors.New(errors.HTTPRequest, "HTTP 请求失败", lastErr)
}

// isNetworkError 检查是否为网络错误
func isNetworkError(err error) bool {
	if err == nil {
		return false
	}

	// 检查超时错误
	if netErr, ok := err.(net.Error); ok {
		return netErr.Timeout()
	}

	return true
}
