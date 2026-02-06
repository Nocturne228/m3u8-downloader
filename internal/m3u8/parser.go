package m3u8

import (
	"fmt"
	"strings"

	"m3u8-downloader/internal/errors"
	"m3u8-downloader/internal/http"
	"m3u8-downloader/internal/logger"
)

// TsSegment 表示一个 TS 文件段
type TsSegment struct {
	Name     string
	URL      string
	Duration float64
}

// EncryptionKey 加密密钥信息
type EncryptionKey struct {
	Method string
	URL    string
	IV     string
	Data   []byte
}

// Manifest M3U8 清单文件
type Manifest struct {
	Segments []*TsSegment
	Key      *EncryptionKey
}

// Parser M3U8 解析器接口
type Parser interface {
	Parse(content string) (*Manifest, error)
}

// M3U8Parser M3U8 解析器实现
type M3U8Parser struct {
	hostURL    string
	httpClient http.Client
	logger     logger.Logger
}

// NewParser 创建新的 M3U8 解析器
func NewParser(hostURL string, httpClient http.Client, logger logger.Logger) Parser {
	return &M3U8Parser{
		hostURL:    hostURL,
		httpClient: httpClient,
		logger:     logger,
	}
}

// Parse 解析 M3U8 清单文件
func (p *M3U8Parser) Parse(content string) (*Manifest, error) {
	if content == "" {
		return nil, errors.New(errors.M3U8Parse, "M3U8 内容为空", nil)
	}

	manifest := &Manifest{
		Segments: make([]*TsSegment, 0),
	}

	lines := strings.Split(content, "\n")
	index := 0

	for _, line := range lines {
		line = strings.TrimSpace(line)

		// 跳过注释和空行
		if line == "" || strings.HasPrefix(line, "#") {
			// 处理加密密钥信息
			if strings.Contains(line, "#EXT-X-KEY") {
				key, err := p.parseKey(line)
				if err == nil {
					manifest.Key = key
				}
			}
			continue
		}

		// 这是一个 TS 段
		index++
		segment, err := p.parseSegment(line, index)
		if err != nil {
			p.logger.Warn("解析 TS 段失败: %v", err)
			continue
		}

		manifest.Segments = append(manifest.Segments, segment)
	}

	if len(manifest.Segments) == 0 {
		return nil, errors.New(errors.M3U8Invalid, "M3U8 中未找到 TS 段", nil)
	}

	p.logger.Info("成功解析 M3U8, 共 %d 个 TS 段", len(manifest.Segments))

	return manifest, nil
}

func (p *M3U8Parser) parseSegment(url string, index int) (*TsSegment, error) {
	url = strings.TrimSpace(url)

	if url == "" {
		return nil, fmt.Errorf("URL 为空")
	}

	// 如果是绝对 URL，直接使用
	if strings.HasPrefix(url, "http") {
		return &TsSegment{
			Name: fmt.Sprintf("%05d.ts", index),
			URL:  url,
		}, nil
	}

	// 否则拼接主机 URL
	url = strings.TrimPrefix(url, "/")
	fullURL := fmt.Sprintf("%s/%s", p.hostURL, url)

	return &TsSegment{
		Name: fmt.Sprintf("%05d.ts", index),
		URL:  fullURL,
	}, nil
}

func (p *M3U8Parser) parseKey(line string) (*EncryptionKey, error) {
	// 检查是否包含 URI
	if !strings.Contains(line, "URI") {
		return nil, fmt.Errorf("密钥信息中未找到 URI")
	}

	// 提取 METHOD
	method := "AES-128"
	if strings.Contains(line, "METHOD=") {
		start := strings.Index(line, "METHOD=")
		end := strings.Index(line[start:], ",")
		if end == -1 {
			end = len(line) - start
		}
		methodPart := line[start+7 : start+end]
		method = strings.Trim(methodPart, "\"")
	}

	// 提取 URI
	uriStart := strings.Index(line, "URI=")
	if uriStart == -1 {
		return nil, fmt.Errorf("未找到 URI 字段")
	}

	// 查找引号对
	firstQuote := strings.Index(line[uriStart:], "\"")
	secondQuote := strings.Index(line[uriStart+firstQuote+1:], "\"")

	if firstQuote == -1 || secondQuote == -1 {
		return nil, fmt.Errorf("URI 字段格式错误")
	}

	keyURL := line[uriStart+firstQuote+1 : uriStart+firstQuote+1+secondQuote]

	// 构造完整 URL
	if !strings.HasPrefix(keyURL, "http") {
		keyURL = fmt.Sprintf("%s/%s", p.hostURL, strings.TrimPrefix(keyURL, "/"))
	}

	key := &EncryptionKey{
		Method: method,
		URL:    keyURL,
	}

	// 尝试下载密钥数据
	data, err := p.httpClient.Get(keyURL)
	if err != nil {
		p.logger.Warn("下载加密密钥失败: %v", err)
		return key, nil
	}

	key.Data = data
	p.logger.Info("成功获取加密密钥: %s", keyURL)

	return key, nil
}
