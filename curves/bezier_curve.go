package curves

// bernsteinPolynomial returns a function, which, given by a list of control points, and a point (0, 1), returns
// a point in the Bézier curve described by these points.
func bernsteinPolynomial(points []Point) func(t float64) Point {
	n := len(points) - 1

	return func(t float64) (res Point) {
		for i, point := range points {
			bern := bernsteinPolynomialPoint(t, i, n)
			res.x += point.x * bern
			res.y += point.y * bern
		}
		return res
	}
}

// curvePoints returns n points on the Bézier curve described by these points, given a list of control points.
func curvePoints(n int, points []Point) (curvePoints []Point) {
	f := float64(n)
	polynomial := bernsteinPolynomial(points)
	for i := float64(0); i < f; i++ {
		t := i / (f - 1)
		curvePoints = append(curvePoints, polynomial(t))
	}
	return curvePoints
}

// defaultTween is the default tween function. It is a quadratic tween function that begins fast and then decelerates.
func defaultTween(n float64) float64 {
	if n < 0 || n > 1 {
		panic("parameter must be between 0.0 and 1.0")
	}
	return -n * (n - 2)
}
