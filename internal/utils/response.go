package utils

import "fmt"

// WVPResult is the unified API response wrapper, matching Java WVPResult<T>
type WVPResult[T any] struct {
	Code int    `json:"code"`
	Msg  string `json:"msg"`
	Data T      `json:"data"`
}

// Error codes matching the Java ErrorCode.java
const (
	CodeSuccess      = 0
	CodeFailure      = 100
	CodeParamError   = 400
	CodeUnauthorized = 401
	CodeForbidden    = 403
	CodeNotFound     = 404
	CodeTimeout      = 408
	CodeNoResponse   = 486
	CodeServerError  = 500
)

func Success[T any](data T) *WVPResult[T] {
	return &WVPResult[T]{
		Code: CodeSuccess,
		Msg:  "成功",
		Data: data,
	}
}

func SuccessWithMsg[T any](data T, msg string) *WVPResult[T] {
	return &WVPResult[T]{
		Code: CodeSuccess,
		Msg:  msg,
		Data: data,
	}
}

func SuccessNoData() *WVPResult[any] {
	return &WVPResult[any]{
		Code: CodeSuccess,
		Msg:  "成功",
	}
}

func Fail(code int, msg string) *WVPResult[any] {
	return &WVPResult[any]{
		Code: code,
		Msg:  msg,
	}
}

func Failf(code int, format string, args ...interface{}) *WVPResult[any] {
	return &WVPResult[any]{
		Code: code,
		Msg:  fmt.Sprintf(format, args...),
	}
}
