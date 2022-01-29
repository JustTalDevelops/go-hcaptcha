package screen

import (
	"github.com/go-gl/mathgl/mgl64"
	"github.com/justtaldevelops/go-hcaptcha/utils"
	"math"
	"math/rand"
)

// Curve is used for generating a curve similar to what a human would make when moving a mouse cursor,
// starting and ending at given positions.
type Curve struct {
	fromPoint, toPoint mgl64.Vec2
	points             []mgl64.Vec2
}

// NewCurve creates a new Curve.
func NewCurve(fromPoint, toPoint mgl64.Vec2, opts *CurveOpts) *Curve {
	h := &Curve{
		fromPoint: fromPoint,
		toPoint:   toPoint,
	}
	h.points = h.generateCurve(opts)
	return h
}

// FromPoint returns the starting point of the curve.
func (h *Curve) FromPoint() mgl64.Vec2 {
	return h.fromPoint
}

// ToPoint returns the ending point of the curve.
func (h *Curve) ToPoint() mgl64.Vec2 {
	return h.toPoint
}

// Points returns the points of the curve.
func (h *Curve) Points() []mgl64.Vec2 {
	return h.points
}

// generateCurve generates a curve according to the parameters in CurveOpts.
func (h *Curve) generateCurve(opts *CurveOpts) []mgl64.Vec2 {
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
func (*Curve) generateInternalKnots(leftBoundary, rightBoundary, downBoundary, upBoundary, knotsCount int) []mgl64.Vec2 {
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
func (h *Curve) generatePoints(knots []mgl64.Vec2) []mgl64.Vec2 {
	midPointsCount := int(math.Max(
		math.Max(
			math.Abs(h.fromPoint.X()-h.toPoint.X()),
			math.Abs(h.fromPoint.Y()-h.toPoint.Y()),
		), 2,
	))

	knots = append([]mgl64.Vec2{h.fromPoint}, append(knots, h.toPoint)...)
	return mgl64.MakeBezierCurve2D(midPointsCount, knots)
}

// distortPoints distorts the curve described by the points, so that the curve is not ideally smooth. Distortion
// happens by randomly, according to normal distribution, adding an offset to some points.
func (h *Curve) distortPoints(points []mgl64.Vec2, distortionMean, distortionStdDev, distortionFrequency float64) []mgl64.Vec2 {
	if distortionFrequency < 0 || distortionFrequency > 1 {
		panic("distortionFrequency must be between 0 and 1")
	}

	distortedPoints := make([]mgl64.Vec2, len(points))
	for i := 1; i < len(points)-1; i++ {
		point := points[i]
		if utils.Chance(distortionFrequency) {
			delta := rand.NormFloat64()*distortionStdDev + distortionMean
			distortedPoints[i] = mgl64.Vec2{point.X(), point.Y() + delta}
		} else {
			distortedPoints[i] = point
		}
	}

	distortedPoints = append([]mgl64.Vec2{points[0]}, append(distortedPoints, points[len(points)-1])...)
	return distortedPoints
}

// tweenPoints chooses a number of target points from the points according to tween function. This controls the
// velocity of mouse movement.
func (h *Curve) tweenPoints(points []mgl64.Vec2, tween func(float64) float64, targetPoints int) []mgl64.Vec2 {
	if targetPoints < 2 {
		panic("targetPoints must be at least 2")
	}

	var tweenedPoints []mgl64.Vec2
	for i := 0; i < targetPoints; i++ {
		index := int(tween(float64(i)/(float64(targetPoints)-1)) * (float64(len(points)) - 1))
		tweenedPoints = append(tweenedPoints, points[index])
	}
	return tweenedPoints
}

// defaultTween is the default tween function. It is a quadratic tween function that begins fast and then decelerates.
func defaultTween(n float64) float64 {
	if n < 0 || n > 1 {
		panic("parameter must be between 0.0 and 1.0")
	}
	return -n * (n - 2)
}
