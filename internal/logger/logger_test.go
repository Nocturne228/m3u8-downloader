package logger

import (
	"testing"
)

// TestParseLevel 测试日志级别解析函数
func TestParseLevel(t *testing.T) {
	tests := []struct {
		name     string
		levelStr string
		expected Level
	}{
		{"debug level", "debug", DebugLevel},
		{"info level", "info", InfoLevel},
		{"warn level", "warn", WarnLevel},
		{"error level", "error", ErrorLevel},
		{"fatal level", "fatal", FatalLevel},
		{"default to info", "unknown", InfoLevel},
		{"empty string", "", InfoLevel},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			level := parseLevel(tt.levelStr)
			if level != tt.expected {
				t.Errorf("Expected level %v, got %v", tt.expected, level)
			}
		})
	}
}

// TestNewLogger 测试日志实例创建
func TestNewLogger(t *testing.T) {
	logger := New("info")

	if logger == nil {
		t.Error("Expected non-nil logger")
	}

	consoleLogger, ok := logger.(*ConsoleLogger)
	if !ok {
		t.Error("Expected ConsoleLogger type")
	}

	if consoleLogger.level != InfoLevel {
		t.Errorf("Expected InfoLevel, got %v", consoleLogger.level)
	}
}

// TestLevelString 测试日志级别字符串转换
func TestLevelString(t *testing.T) {
	tests := []struct {
		level    Level
		expected string
	}{
		{DebugLevel, "DEBUG"},
		{InfoLevel, "INFO"},
		{WarnLevel, "WARN"},
		{ErrorLevel, "ERROR"},
		{FatalLevel, "FATAL"},
		{Level(999), "UNKNOWN"},
	}

	for _, tt := range tests {
		if tt.level.String() != tt.expected {
			t.Errorf("Expected '%s', got '%s'", tt.expected, tt.level.String())
		}
	}
}

// TestLoggerMethods 测试日志输出方法不会panic
func TestLoggerMethods(t *testing.T) {
	logger := New("info")

	// 这些应该不会panic
	logger.Debug("debug message")
	logger.Info("info message")
	logger.Warn("warn message")
	logger.Error("error message")

	logger.DebugWithFields("debug with fields", map[string]interface{}{"key": "value"})
	logger.InfoWithFields("info with fields", map[string]interface{}{"key": "value"})
	logger.ErrorWithFields("error with fields", map[string]interface{}{"key": "value"})
}

// TestDebugLevelFiltering 测试日志级别过滤
func TestDebugLevelFiltering(t *testing.T) {
	logger := New("error")

	// 这些都应该被过滤（不输出）
	logger.Debug("should be filtered")
	logger.Info("should be filtered")
	logger.Warn("should be filtered")

	// 这个应该被输出
	logger.Error("should be output")
}
