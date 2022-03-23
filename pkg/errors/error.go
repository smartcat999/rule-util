package errors

import (
	"fmt"
)

//type Error *error

type Error struct {
	Code    string
	Message string
	Args    []interface{}
}

//type Error device.ErrorMsg

func NewCodeError(code string, msg string, args ...interface{}) *Error {
	return &Error{
		Code:    code,
		Message: msg,
		Args:    args,
	}
}

type errorMsg struct {
	En   string
	ZhCN string
}

func (e Error) Error() string {
	return e.Code
}

//var DefaultErrorMessage = map[string]*Error{
//	constant.INTERNAL_ERROR:   InternalErrorMsg,
//	constant.GRPC_ERROR:       InternalGrpcErrorMsg,
//	constant.IOT_PARAMS_ERROR: ParamsErrorMsg,
//}

//内部错误（未知错误）

type InternalErrorCode struct {
	msg  string
	args []interface{}
}

func NewInternalError(msg string, args ...interface{}) error {
	return &InternalErrorCode{msg: msg, args: args}
}
func (i *InternalErrorCode) Error() string {
	return fmt.Sprintf(i.msg, i.args...)
}

//数据库错误
type DBErrorCode struct {
	msg  string
	args []interface{}
}

func NewDBError(msg string, args ...interface{}) error {
	return &DBErrorCode{msg: msg, args: args}
}
func (i *DBErrorCode) Error() string {
	return fmt.Sprintf(i.msg, i.args...)
}
