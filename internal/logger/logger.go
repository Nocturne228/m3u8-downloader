package logger

import (
	"fmt"
	"os"
	"sync"
	"time"
)

/*
Catppuccin Mocha ANSI colors
https://catppuccin.com/palette
*/
const (
	reset = "\033[0m"

	text      = "\033[38;2;205;214;244m"
	subtext1  = "\033[38;2;186;194;222m"
	surface1  = "\033[38;2;69;71;90m"

	sky       = "\033[38;2;137;220;235m"
	green     = "\033[38;2;166;227;161m"
	yellow    = "\033[38;2;249;226;175m"
	red       = "\033[38;2;243;139;168m"
	maroon    = "\033[38;2;235;160;172m"
	lavender  = "\033[38;2;180;190;254m"
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

// levelColor returns Catppuccin color for level
func levelColor(l Level) string {
	switch l {
	case DebugLevel:
		return sky
	case InfoLevel:
		return green
	case WarnLevel:
		return yellow
	case ErrorLevel:
		return red
	case FatalLevel:
		return maroon
	default:
		return text
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
	return &ConsoleLogger{
		level: parseLevel(levelStr),
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
	ts := time.Now().Format("2006-01-02 15:04:05")

	if len(args) > 0 {
		msg = fmt.Sprintf(msg, args...)
	}

	prefix := fmt.Sprintf(
		"%s[%s]%s %s[%s]%s",
		surface1, ts, reset,
		levelColor(level), level, reset,
	)

	printWithProgressRedraw(func() {
		fmt.Printf("%s %s%s%s\n", prefix, text, msg, reset)
	})
}

func (l *ConsoleLogger) logWithFields(level Level, msg string, fields map[string]interface{}) {
	ts := time.Now().Format("2006-01-02 15:04:05")

	fieldsStr := ""
	for k, v := range fields {
		fieldsStr += fmt.Sprintf(
			" %s%s%s=%s%v%s",
			lavender, k, reset,
			subtext1, v, reset,
		)
	}

	prefix := fmt.Sprintf(
		"%s[%s]%s %s[%s]%s",
		surface1, ts, reset,
		levelColor(level), level, reset,
	)

	printWithProgressRedraw(func() {
		fmt.Printf("%s %s%s%s\n", prefix, text, msg, fieldsStr)
	})
}

/* ---------- progress redraw support ---------- */

var (
	progressRedraw   func()
	progressRedrawMu sync.Mutex
)

func printWithProgressRedraw(printFn func()) {
	progressRedrawMu.Lock()
	hasRedraw := progressRedraw != nil
	progressRedrawMu.Unlock()

	if hasRedraw {
		fmt.Print("\r\033[K")
	}

	printFn()

	if hasRedraw {
		progressRedrawMu.Lock()
		cb := progressRedraw
		progressRedrawMu.Unlock()
		if cb != nil {
			cb()
		}
	}
}

// RegisterProgressRedraw registers a callback that will be invoked after each
// log line so a long-running progress line can be redrawn without being
// interrupted by log output.
func RegisterProgressRedraw(cb func()) {
	progressRedrawMu.Lock()
	defer progressRedrawMu.Unlock()
	progressRedraw = cb
}

