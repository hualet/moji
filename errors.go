package moji

import "errors"

var (
	// ErrAppCodeEnv means the app code environment variable is not set.
	ErrAppCodeEnv = errors.New(appCodeEnv + " environment variable is not set")
)
