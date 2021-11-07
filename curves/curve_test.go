package curves

import (
	"github.com/go-vgo/robotgo"
	"github.com/justtaldevelops/hcaptcha-solver-go/utils"
	"testing"
	"time"
)

func TestCurve(t *testing.T) {
	mousePosX, mousePosY := robotgo.GetMousePos()
	start := Point{float64(mousePosX), float64(mousePosY)}
	end := Point{1187, 719}

	humanCurve := NewHumanCurve(start, end, &CurveOpts{})
	pause := utils.Between(5, 15)
	for _, point := range humanCurve.points {
		robotgo.MoveMouse(int(point.x), int(point.y))
		time.Sleep(time.Duration(pause) * time.Millisecond)
	}
}
