package utils

import (
	"fmt"
)

const skip = 3

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

	//pc, file, line, ok := runtime.Caller(skip)
	//if !ok {
	//	return fmt.Sprintf("%s", r.Message)
	//}
	//
	//fileName := path.Base(file)
	//lenFuncName := len(strings.Split(runtime.FuncForPC(pc).Name(), "."))        // 先取长度
	//funcName := strings.Split(runtime.FuncForPC(pc).Name(), ".")[lenFuncName-1] // 从结果main.main中，取后面一个函数名就好，不然太多有点难看
	//if funcName != "doPrintf" {
	//	fmt.Printf("Error: func_name: %v, file_name:[%v], line:[%v] :%s\n", funcName, fileName, line, r.Message)
	//}
	//errMsg := fmt.Sprintf("func_name: [%v], file_name:[%v], line:[%v] :%s", funcName, fileName, line, r.Message)
	//return fmt.Sprintf("%s", errMsg)
}

func NewRockError(httpCode, errCode int, msg string) *RockError {
	return &RockError{
		HttpCode: httpCode,
		ErrCode:  errCode,
		Message:  msg,
	}
}
