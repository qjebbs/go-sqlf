package sqlb

import "github.com/qjebbs/go-sqlf/v2"

func convertFragmentBuilders[T sqlf.Builder](builders []T) []sqlf.Builder {
	r := make([]sqlf.Builder, len(builders))
	for i, b := range builders {
		r[i] = b
	}
	return r
}
