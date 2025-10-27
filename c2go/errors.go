package c2go

import (
	"fmt"
)

// CError 封装C错误
type CError struct {
	Code    int
	Message string
}

func (e *CError) Error() string {
	return fmt.Sprintf("C error %d: %s", e.Code, e.Message)
}

// SimpleError 创建简单的错误
func SimpleError(message string) error {
	return &CError{
		Code:    -1,
		Message: message,
	}
}

// ErrorFromCode 从错误代码创建错误
func ErrorFromCode(code int) error {
	if code == 0 {
		return nil
	}
	return &CError{
		Code:    code,
		Message: fmt.Sprintf("Error code %d", code),
	}
}

// ErrorFromString 从字符串创建错误
func ErrorFromString(message string) error {
	if message == "" {
		return nil
	}
	return &CError{
		Code:    -1,
		Message: message,
	}
}
