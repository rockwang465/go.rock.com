package utils

import "fmt"

type HTTPError struct {
	Error     string `json:"error" example:"error message here"`
	ErrorCode int    `json:"error_code" example:"50000001"`
}

type RockError struct {
	HttpCode int
	ErrCode  int
	Message  string
}

// NewRockError这个error中必须要有Error() string这样的方法+返回值，这是基础的error样例
func (r *RockError) Error() string {
	return fmt.Sprintf("%s", r.Message)
}

func NewRockError(httpCode, errCode int, msg string) *RockError {
	return &RockError{
		HttpCode: httpCode,
		ErrCode:  errCode,
		Message:  msg,
	}
}
