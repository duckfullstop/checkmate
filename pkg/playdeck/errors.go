package playdeck

import "errors"

// Errors throwable by this module.
var (
	ErrDeckEmpty         = errors.New("deck is empty")
	ErrDeckUninitialized = errors.New("deck is uninitialized")
)
