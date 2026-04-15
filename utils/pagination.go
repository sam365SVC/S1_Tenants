package utils

import "math"

func CalcularPagination(totalCount int64, size int) (pageEnd int) {
	pageEnd = int(math.Ceil(float64(totalCount) / float64(size)))
	return pageEnd
}
