package screen

import (
	"github.com/justtaldevelops/hcaptcha-solver-go/utils"
	"math"
)

// Point represents a point in 2D space.
type Point struct {
	// X, Y represent the coordinates of the point.
	X, Y float64
}

// merge does a merge on two int slices into a slice of knots.
func merge(a, b []int) []Point {
	if len(a) != len(b) {
		panic("arguments must be of same length")
	}

	r := make([]Point, len(a), len(a))
	for i, e := range a {
		r[i] = Point{float64(e), float64(b[i])}
	}

	return r
}

// knots generates a random choice of knots based on the size provided.
func knots(firstBoundary, secondBoundary, size int) []int {
	result := make([]int, size)
	for i := 0; i < size; i++ {
		result = append(result, utils.Between(firstBoundary, secondBoundary))
	}
	return result
}

// binomial returns the binomial coefficient "n choose k".
func binomial(n, k int) float64 {
	return float64(factorial(n)) / float64(factorial(k)*factorial(n-k))
}

// bernsteinPolynomialPoint calculates the i-th component of a bernstein polynomial of degree n.
func bernsteinPolynomialPoint(x float64, i, n int) float64 {
	return binomial(n, i) * math.Pow(x, float64(i)) * (math.Pow(1-x, float64(n-i)))
}

// factorial returns the factorial of n.
func factorial(n int) int {
	if n >= 1 {
		return n * factorial(n-1)
	}
	return 1
}
