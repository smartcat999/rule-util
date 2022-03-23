package prometheus

import "errors"

const (
	ApiStatusFailure = "failure"
	ApiStatusSuccess = "success"
)

var (
	ErrResolvePort = errors.New("resolve port failed")
)
