package errors

import (
	stderrors "errors"
	"testing"
)

func TestNew(t *testing.T) {
	err := New(HTTPRequest, "HTTP 请求失败", stderrors.New("connection refused"))

	if err == nil {
		t.Fatal("Expected non-nil error")
	}

	if err.Code != HTTPRequest {
		t.Errorf("Expected code %s, got %s", HTTPRequest, err.Code)
	}

	if err.Message != "HTTP 请求失败" {
		t.Errorf("Expected message 'HTTP 请求失败', got %s", err.Message)
	}
}

func TestIsCode(t *testing.T) {
	err := New(HTTPRequest, "test", nil)

	if !IsCode(err, HTTPRequest) {
		t.Error("Expected IsCode to return true for HTTPRequest")
	}

	if IsCode(err, M3U8Parse) {
		t.Error("Expected IsCode to return false for M3U8Parse")
	}

	// 测试非 Error 类型
	regularErr := stderrors.New("regular error")
	if IsCode(regularErr, HTTPRequest) {
		t.Error("Expected IsCode to return false for non-Error type")
	}
}

func TestErrorString(t *testing.T) {
	tests := []struct {
		name      string
		code      string
		message   string
		err       error
		expectMsg string
	}{
		{
			name:      "with underlying error",
			code:      HTTPRequest,
			message:   "HTTP 请求失败",
			err:       stderrors.New("connection refused"),
			expectMsg: "[HTTP_REQUEST] HTTP 请求失败: connection refused",
		},
		{
			name:      "without underlying error",
			code:      HTTPRequest,
			message:   "HTTP 请求失败",
			err:       nil,
			expectMsg: "[HTTP_REQUEST] HTTP 请求失败",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := New(tt.code, tt.message, tt.err)
			if err.Error() != tt.expectMsg {
				t.Errorf("Expected '%s', got '%s'", tt.expectMsg, err.Error())
			}
		})
	}
}

func TestErrorUnwrap(t *testing.T) {
	underlyingErr := stderrors.New("underlying error")
	err := New(HTTPRequest, "HTTP 请求失败", underlyingErr)

	unwrapped := err.Unwrap()
	if unwrapped.Error() != "underlying error" {
		t.Errorf("Expected 'underlying error', got '%s'", unwrapped.Error())
	}
}
