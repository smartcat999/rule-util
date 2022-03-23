package errors

import (
	"errors"
	"fmt"
)

var (
	ErrEncodeError  = New("encode error")
	ErrDecodeError  = New("decode error")
	ErrNeedRestart  = New("connection need restart")
	ErrNeedSleep    = New("connection need sleep")
	ErrMaySkip      = New("connection may skip")
	ErrOpenSource   = New("open source error")
	ErrOpenSink     = New("open sink error")
	ErrReopenSource = New("reopen source error")
	ErrReopenSink   = New("reopen sink error")
	ErrRerunSource  = New("rerun source error")
	ErrRerunSink    = New("rerun sink error")
	ErrSendError    = New("send error")
)

func New(text string) error {
	return errors.New(text)
}

func Errorf(format string, a ...interface{}) error {
	return fmt.Errorf(format, a...)
}
