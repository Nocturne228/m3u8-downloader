package config

import "time"

// Config 应用程序配置
type Config struct {
	HTTP     HTTPConfig
	Download DownloadConfig
	FFmpeg   FFmpegConfig
	Log      LogConfig
}

// HTTPConfig HTTP 相关配置
type HTTPConfig struct {
	Timeout    time.Duration
	MaxRetries int
	UserAgent  string
}

// DownloadConfig 下载相关配置
type DownloadConfig struct {
	MaxGoroutines      int
	TsNameTemplate     string
	LossTolerance      float64
	HostType           string
	SavePath           string
	AutoClear          bool
	InsecureSkipVerify bool
	Cookie             string
}

// FFmpegConfig FFmpeg 相关配置
type FFmpegConfig struct {
	Enabled bool
	Path    string
	Options []string
}

// LogConfig 日志配置
type LogConfig struct {
	Level  string
	Format string
}

// DefaultConfig 返回默认配置
func DefaultConfig() *Config {
	return &Config{
		HTTP: HTTPConfig{
			Timeout:    5 * time.Second,
			MaxRetries: 5,
			UserAgent:  "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_13_6) AppleWebKit/537.36",
		},
		Download: DownloadConfig{
			MaxGoroutines:      24,
			TsNameTemplate:     "%05d.ts",
			LossTolerance:      0.1,
			HostType:           "v1",
			AutoClear:          true,
			InsecureSkipVerify: false,
		},
		FFmpeg: FFmpegConfig{
			Enabled: true,
			Path:    "ffmpeg",
			Options: []string{"-c", "copy", "-y"},
		},
		Log: LogConfig{
			Level:  "info",
			Format: "text",
		},
	}
}

// Validate 验证配置
func (c *Config) Validate() error {
	if c.Download.MaxGoroutines <= 0 || c.Download.MaxGoroutines > 256 {
		return NewConfigError("下载线程数必须在 1-256 之间")
	}

	if c.HTTP.Timeout <= 0 {
		return NewConfigError("HTTP 超时时间必须大于 0")
	}

	if c.HTTP.MaxRetries < 0 {
		return NewConfigError("最大重试次数不能为负数")
	}

	if c.Download.LossTolerance < 0 || c.Download.LossTolerance >= 1 {
		return NewConfigError("损失率必须在 0-1 之间")
	}

	return nil
}

// ConfigError 配置错误
type ConfigError struct {
	message string
}

func (e *ConfigError) Error() string {
	return "配置错误: " + e.message
}

func NewConfigError(message string) *ConfigError {
	return &ConfigError{message: message}
}
