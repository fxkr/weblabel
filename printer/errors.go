package printer

import (
	"errors"
)

var (
	// Client errors
	ErrBadRequest = errors.New("Bad request.")
	ErrNotFound   = errors.New("Not found.")
)
