package sqlb

import "github.com/qjebbs/go-sqlf/v2"

func convertFragmentBuilders[T sqlf.FragmentBuilder](builders []T) []sqlf.FragmentBuilder {
	r := make([]sqlf.FragmentBuilder, len(builders))
	for i, b := range builders {
		r[i] = b
	}
	return r
}
