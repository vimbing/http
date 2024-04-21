package http

import "errors"

var (
	ErrRequestTimedOut = errors.New("request timed out, context cancelled")
)
