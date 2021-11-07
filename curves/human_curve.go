package curves

import (
	"github.com/justtaldevelops/hcaptcha-solver-go/utils"
	"math"
	"math/rand"
)

// HumanCurve is used for generating a curve similar to what a human would make when moving a mouse cursor,
// starting and ending at given positions.
type HumanCurve struct {
	fromPoint, toPoint Point
	points             []Point
}

// NewHumanCurve creates a new HumanCurve.
func NewHumanCurve(fromPoint, toPoint Point, opts *CurveOpts) *HumanCurve {
	h := &HumanCurve{
		fromPoint: fromPoint,
		toPoint:   toPoint,
	}
	h.points = h.generateCurve(opts)
	return h
}

// FromPoint returns the starting point of the curve.
func (h *HumanCurve) FromPoint() Point {
	return h.fromPoint
}

// ToPoint returns the ending point of the curve.
func (h *HumanCurve) ToPoint() Point {
	return h.toPoint
}

// Points returns the points of the curve.
func (h *HumanCurve) Points() []Point {
	return h.points
}

// generateCurve generates a curve according to the parameters in CurveOpts.
func (h *HumanCurve) generateCurve(opts *CurveOpts) []Point {
	h.defaultCurveOpts(opts)

	offsetBoundaryX := *opts.OffsetBoundaryX
	offsetBoundaryY := *opts.OffsetBoundaryY

	leftBoundary := *opts.LeftBoundary - offsetBoundaryX
	rightBoundary := *opts.RightBoundary + offsetBoundaryX
	downBoundary := *opts.DownBoundary - offsetBoundaryY
	upBoundary := *opts.UpBoundary + offsetBoundaryY
	count := *opts.KnotsCount
	distortionMean := *opts.DistortionMean
	distortionStdDev := *opts.DistortionStdDev
	distortionFrequency := *opts.DistortionFrequency
	targetPoints := *opts.TargetPoints

	internalKnots := h.generateInternalKnots(leftBoundary, rightBoundary, downBoundary, upBoundary, count)
	points := h.generatePoints(internalKnots)
	points = h.distortPoints(points, distortionMean, distortionStdDev, distortionFrequency)
	points = h.tweenPoints(points, opts.Tween, targetPoints)
	return points
}

// generateInternalKnots generates the internal knots for the curve.
func (*HumanCurve) generateInternalKnots(leftBoundary, rightBoundary, downBoundary, upBoundary, knotsCount int) []Point {
	if knotsCount < 0 {
		panic("knotsCount can't be negative")
	}
	if leftBoundary > rightBoundary {
		panic("leftBoundary must be less than or equal to rightBoundary")
	}
	if downBoundary > upBoundary {
		panic("downBoundary must be less than or equal to upBoundary")
	}

	knotsX := knots(leftBoundary, rightBoundary, knotsCount)
	knotsY := knots(downBoundary, upBoundary, knotsCount)
	return merge(knotsX, knotsY)
}

// generatePoints generates Bézier curve points on a curve, according to the internal knots passed as parameter.
func (h *HumanCurve) generatePoints(knots []Point) []Point {
	midPointsCount := int(math.Max(
		math.Max(
			math.Abs(h.fromPoint.x-h.toPoint.x),
			math.Abs(h.fromPoint.y-h.toPoint.y),
		), 2,
	))

	knots = append([]Point{h.fromPoint}, append(knots, h.toPoint)...)
	return curvePoints(midPointsCount, knots)
}

// distortPoints distorts the curve described by the points, so that the curve is not ideally smooth. Distortion
// happens by randomly, according to normal distribution, adding an offset to some points.
func (h *HumanCurve) distortPoints(points []Point, distortionMean, distortionStdDev, distortionFrequency float64) []Point {
	if distortionFrequency < 0 || distortionFrequency > 1 {
		panic("distortionFrequency must be between 0 and 1")
	}

	distortedPoints := make([]Point, len(points))
	for i := 1; i < len(points)-1; i++ {
		point := points[i]
		if utils.Chance(distortionFrequency) {
			delta := rand.NormFloat64()*distortionStdDev + distortionMean
			distortedPoints[i] = Point{x: point.x, y: point.y + delta}
		} else {
			distortedPoints[i] = point
		}
	}

	distortedPoints = append([]Point{points[0]}, append(distortedPoints, points[len(points)-1])...)
	return distortedPoints
}

// tweenPoints chooses a number of target points from the points according to tween function. This controls the
// velocity of mouse movement.
func (h *HumanCurve) tweenPoints(points []Point, tween func(float64) float64, targetPoints int) []Point {
	if targetPoints < 2 {
		panic("targetPoints must be at least 2")
	}

	var tweenedPoints []Point
	for i := 0; i < targetPoints; i++ {
		index := int(tween(float64(i)/(float64(targetPoints)-1)) * (float64(len(points)) - 1))
		tweenedPoints = append(tweenedPoints, points[index])
	}

	return tweenedPoints
}
