package config

import (
	"testing"
	"time"
)

// TestDefaultConfig 测试默认配置
func TestDefaultConfig(t *testing.T) {
	cfg := DefaultConfig()

	if cfg == nil {
		t.Fatal("DefaultConfig 返回 nil")
	}

	if cfg.Download.MaxGoroutines != 24 {
		t.Errorf("期望 MaxGoroutines=24, 得到 %d", cfg.Download.MaxGoroutines)
	}

	if cfg.HTTP.Timeout != 5*time.Second {
		t.Errorf("期望 Timeout=5s, 得到 %v", cfg.HTTP.Timeout)
	}

	if cfg.Download.AutoClear != true {
		t.Error("期望 AutoClear=true")
	}
}

// TestConfigValidate 测试配置验证
func TestConfigValidate(t *testing.T) {
	tests := []struct {
		name    string
		cfg     *Config
		wantErr bool
	}{
		{
			name:    "有效配置",
			cfg:     DefaultConfig(),
			wantErr: false,
		},
		{
			name: "无效的线程数 - 太小",
			cfg: &Config{
				Download: DownloadConfig{
					MaxGoroutines: 0,
				},
			},
			wantErr: true,
		},
		{
			name: "无效的线程数 - 太大",
			cfg: &Config{
				Download: DownloadConfig{
					MaxGoroutines: 300,
				},
				HTTP: HTTPConfig{
					Timeout: 5 * time.Second,
				},
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.cfg.Validate()
			if (err != nil) != tt.wantErr {
				t.Errorf("Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
