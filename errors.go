package sqls

import "fmt"

var (
	// ErrInvalidIndex is returned when the reference index is invalid.
	ErrInvalidIndex = fmt.Errorf("invalid index")
)
