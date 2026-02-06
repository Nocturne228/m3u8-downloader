package errors

import "fmt"

// Error 表示一个应用程序错误
type Error struct {
	Code    string
	Message string
	Err     error
}

func (e *Error) Error() string {
	if e.Err != nil {
		return fmt.Sprintf("[%s] %s: %v", e.Code, e.Message, e.Err)
	}
	return fmt.Sprintf("[%s] %s", e.Code, e.Message)
}

func (e *Error) Unwrap() error {
	return e.Err
}

// New 创建新错误
func New(code, message string, err error) *Error {
	return &Error{
		Code:    code,
		Message: message,
		Err:     err,
	}
}

// 错误码定义
const (
	HTTPRequest    = "HTTP_REQUEST"
	HTTPTimeout    = "HTTP_TIMEOUT"
	HTTPStatus     = "HTTP_STATUS"
	M3U8Parse      = "M3U8_PARSE"
	M3U8Invalid    = "M3U8_INVALID"
	DownloadFailed = "DOWNLOAD_FAILED"
	FFmpegNotFound = "FFMPEG_NOT_FOUND"
	FFmpegFailed   = "FFMPEG_FAILED"
	MergeFailed    = "MERGE_FAILED"
	FileRead       = "FILE_READ"
	FileWrite      = "FILE_WRITE"
	DirCreate      = "DIR_CREATE"
	InvalidURL     = "INVALID_URL"
	InvalidConfig  = "INVALID_CONFIG"
)

// IsCode 检查错误是否为特定错误码
func IsCode(err error, code string) bool {
	if e, ok := err.(*Error); ok {
		return e.Code == code
	}
	return false
}
