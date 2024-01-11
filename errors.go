package sqlf

import "errors"

var (
	// ErrInvalidIndex is returned when the reference index is invalid.
	ErrInvalidIndex = errors.New("invalid index")
)
