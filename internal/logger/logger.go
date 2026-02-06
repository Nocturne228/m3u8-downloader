package logger

import (
	"fmt"
	"os"
	"time"
)

// Level 日志级别
type Level int

const (
	DebugLevel Level = iota
	InfoLevel
	WarnLevel
	ErrorLevel
	FatalLevel
)

func (l Level) String() string {
	switch l {
	case DebugLevel:
		return "DEBUG"
	case InfoLevel:
		return "INFO"
	case WarnLevel:
		return "WARN"
	case ErrorLevel:
		return "ERROR"
	case FatalLevel:
		return "FATAL"
	default:
		return "UNKNOWN"
	}
}

// Logger 日志接口
type Logger interface {
	Debug(msg string, args ...interface{})
	Info(msg string, args ...interface{})
	Warn(msg string, args ...interface{})
	Error(msg string, args ...interface{})
	Fatal(msg string, args ...interface{})
	DebugWithFields(msg string, fields map[string]interface{})
	InfoWithFields(msg string, fields map[string]interface{})
	ErrorWithFields(msg string, fields map[string]interface{})
}

// ConsoleLogger 控制台日志实现
type ConsoleLogger struct {
	level Level
}

// New 创建新的文本日志记录器
func New(levelStr string) Logger {
	level := parseLevel(levelStr)
	return &ConsoleLogger{
		level: level,
	}
}

func parseLevel(levelStr string) Level {
	switch levelStr {
	case "debug":
		return DebugLevel
	case "warn":
		return WarnLevel
	case "error":
		return ErrorLevel
	case "fatal":
		return FatalLevel
	default:
		return InfoLevel
	}
}

func (l *ConsoleLogger) Debug(msg string, args ...interface{}) {
	if l.level <= DebugLevel {
		l.log(DebugLevel, msg, args...)
	}
}

func (l *ConsoleLogger) Info(msg string, args ...interface{}) {
	if l.level <= InfoLevel {
		l.log(InfoLevel, msg, args...)
	}
}

func (l *ConsoleLogger) Warn(msg string, args ...interface{}) {
	if l.level <= WarnLevel {
		l.log(WarnLevel, msg, args...)
	}
}

func (l *ConsoleLogger) Error(msg string, args ...interface{}) {
	if l.level <= ErrorLevel {
		l.log(ErrorLevel, msg, args...)
	}
}

func (l *ConsoleLogger) Fatal(msg string, args ...interface{}) {
	l.log(FatalLevel, msg, args...)
	os.Exit(1)
}

func (l *ConsoleLogger) DebugWithFields(msg string, fields map[string]interface{}) {
	if l.level <= DebugLevel {
		l.logWithFields(DebugLevel, msg, fields)
	}
}

func (l *ConsoleLogger) InfoWithFields(msg string, fields map[string]interface{}) {
	if l.level <= InfoLevel {
		l.logWithFields(InfoLevel, msg, fields)
	}
}

func (l *ConsoleLogger) ErrorWithFields(msg string, fields map[string]interface{}) {
	if l.level <= ErrorLevel {
		l.logWithFields(ErrorLevel, msg, fields)
	}
}

func (l *ConsoleLogger) log(level Level, msg string, args ...interface{}) {
	timestamp := time.Now().Format("2006-01-02 15:04:05")
	prefix := fmt.Sprintf("[%s] [%s]", timestamp, level)

	if len(args) > 0 {
		msg = fmt.Sprintf(msg, args...)
	}

	fmt.Printf("%s %s\n", prefix, msg)
}

func (l *ConsoleLogger) logWithFields(level Level, msg string, fields map[string]interface{}) {
	timestamp := time.Now().Format("2006-01-02 15:04:05")
	prefix := fmt.Sprintf("[%s] [%s]", timestamp, level)

	fieldsStr := ""
	for k, v := range fields {
		fieldsStr += fmt.Sprintf(" %s=%v", k, v)
	}

	fmt.Printf("%s %s%s\n", prefix, msg, fieldsStr)
}
