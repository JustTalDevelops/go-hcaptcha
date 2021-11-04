package curves

// HumanCurve is used for generating a curve similar to what a human would make when moving a mouse cursor,
// starting and ending at given positions.
type HumanCurve struct {
	fromPoint, toPoint Point
	points             []Point
}

// TODO
