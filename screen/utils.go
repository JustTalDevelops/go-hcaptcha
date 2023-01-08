package screen

import (
	"github.com/go-gl/mathgl/mgl64"
	"github.com/justtaldevelops/go-hcaptcha/utils"
)

// merge does a merge on two int slices into a slice of knots.
func merge(a, b []int) []mgl64.Vec2 {
	if len(a) != len(b) {
		panic("arguments must be of same length")
	}

	r := make([]mgl64.Vec2, len(a), len(a))
	for i, e := range a {
		r[i] = mgl64.Vec2{float64(e), float64(b[i])}
	}
	return r
}

// knots generates a random choice of knots based on the size provided.
func knots(firstBoundary, secondBoundary, size int) []int {
	result := make([]int, 0, size)
	for i := 0; i < size; i++ {
		result = append(result, utils.Between(firstBoundary, secondBoundary))
	}
	return result
}
