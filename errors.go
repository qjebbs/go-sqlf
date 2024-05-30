package sqlf

import "errors"

var (
	// ErrInvalidIndex is returned when the reference index is invalid.
	// It's a required behaviour for a custom #func to be compatible with #join.
	ErrInvalidIndex = errors.New("invalid index")
)
