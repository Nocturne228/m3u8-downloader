package m3u8

import (
	"net/url"
	"path/filepath"
	"strings"

	"m3u8-downloader/internal/errors"
	"m3u8-downloader/internal/http"
	"m3u8-downloader/internal/logger"
)

// Fetcher M3U8 获取接口
type Fetcher interface {
	FetchManifest(m3u8URL string, cookie string) (*Manifest, error)
}

// M3U8Fetcher M3U8 获取器实现
type M3U8Fetcher struct {
	httpClient http.Client
	logger     logger.Logger
}

// NewFetcher 创建新的 M3U8 获取器
func NewFetcher(httpClient http.Client, logger logger.Logger) Fetcher {
	return &M3U8Fetcher{
		httpClient: httpClient,
		logger:     logger,
	}
}

// FetchManifest 获取 M3U8 清单文件
func (f *M3U8Fetcher) FetchManifest(m3u8URL string, cookie string) (*Manifest, error) {
	// 验证 URL
	if !strings.HasPrefix(m3u8URL, "http") {
		return nil, errors.New(errors.InvalidURL, "M3U8 URL 必须以 http 或 https 开头", nil)
	}

	f.logger.Info("获取 M3U8 清单: %s", m3u8URL)

	// 获取 M3U8 内容
	var content []byte
	var err error

	if cookie != "" {
		content, err = f.httpClient.GetWithCookie(m3u8URL, cookie)
	} else {
		content, err = f.httpClient.Get(m3u8URL)
	}

	if err != nil {
		return nil, errors.New(errors.M3U8Parse, "获取 M3U8 文件失败", err)
	}

	// 提取主机
	hostURL, err := ExtractHost(m3u8URL, "v1")
	if err != nil {
		return nil, err
	}

	// 创建解析器并解析
	parser := NewParser(hostURL, f.httpClient, f.logger)
	manifest, err := parser.Parse(string(content))
	if err != nil {
		return nil, err
	}

	return manifest, nil
}

// ExtractHost 从 M3U8 URL 提取主机
func ExtractHost(m3u8URL string, hostType string) (string, error) {
	if m3u8URL == "" {
		return "", errors.New(errors.InvalidURL, "URL 为空", nil)
	}

	u, err := url.Parse(m3u8URL)
	if err != nil {
		return "", errors.New(errors.InvalidURL, "URL 解析失败", err)
	}

	var host string
	switch hostType {
	case "v1":
		// 协议 + 主机 + 路径目录
		host = u.Scheme + "://" + u.Host + filepath.Dir(u.EscapedPath())
	case "v2":
		// 协议 + 主机
		host = u.Scheme + "://" + u.Host
	default:
		host = u.Scheme + "://" + u.Host + filepath.Dir(u.EscapedPath())
	}

	return host, nil
}
