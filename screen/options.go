package screen

import "math"

// CurveOpts contains options for a curve.
type CurveOpts struct {
	// OffsetBoundaryX is the boundary for X of the offset curve.
	OffsetBoundaryX *int
	// OffsetBoundaryY is the boundary for Y of the offset curve.
	OffsetBoundaryY *int
	// LeftBoundary is the boundary for the left of the curve.
	LeftBoundary *int
	// RightBoundary is the boundary for the right of the curve.
	RightBoundary *int
	// DownBoundary is the boundary for the bottom of the curve.
	DownBoundary *int
	// UpBoundary is the boundary for the top of the curve.
	UpBoundary *int
	// KnotsCount is the number of knots in the curve.
	KnotsCount *int
	// DistortionMean is the mean of the distortion.
	DistortionMean *float64
	// DistortionStdDev is the standard deviation of the distortion.
	DistortionStdDev *float64
	// DistortionFrequency is the frequency of the distortion.
	DistortionFrequency *float64
	// Tween is the function that tweens values.
	Tween func(float64) float64
	// TargetPoints is the target points of the curve.
	TargetPoints *int
}

// defaultCurveOpts returns the default curve options.
func (h *Curve) defaultCurveOpts(opts *CurveOpts) {
	defaultOffsetBoundaryX := 100
	if opts.OffsetBoundaryX == nil {
		opts.OffsetBoundaryX = &defaultOffsetBoundaryX
	}

	defaultOffsetBoundaryY := 100
	if opts.OffsetBoundaryY == nil {
		opts.OffsetBoundaryY = &defaultOffsetBoundaryY
	}

	defaultLeftBoundary := int(math.Min(h.fromPoint.X(), h.toPoint.X()))
	if opts.LeftBoundary == nil {
		opts.LeftBoundary = &defaultLeftBoundary
	}

	defaultRightBoundary := int(math.Max(h.fromPoint.X(), h.toPoint.X()))
	if opts.RightBoundary == nil {
		opts.RightBoundary = &defaultRightBoundary
	}

	defaultDownBoundary := int(math.Min(h.fromPoint.Y(), h.toPoint.Y()))
	if opts.DownBoundary == nil {
		opts.DownBoundary = &defaultDownBoundary
	}

	defaultUpBoundary := int(math.Max(h.fromPoint.Y(), h.toPoint.Y()))
	if opts.UpBoundary == nil {
		opts.UpBoundary = &defaultUpBoundary
	}

	knotsCount := 2
	if opts.KnotsCount == nil {
		opts.KnotsCount = &knotsCount
	}

	distortionMean := 1.0
	if opts.DistortionMean == nil {
		opts.DistortionMean = &distortionMean
	}

	distortionStdDev := 0.6
	if opts.DistortionStdDev == nil {
		opts.DistortionStdDev = &distortionStdDev
	}

	distortionFrequency := 0.5
	if opts.DistortionFrequency == nil {
		opts.DistortionFrequency = &distortionFrequency
	}

	if opts.Tween == nil {
		opts.Tween = defaultTween
	}

	targetPoints := 300
	if opts.TargetPoints == nil {
		opts.TargetPoints = &targetPoints
	}
}
