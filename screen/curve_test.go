package screen

import (
	"github.com/go-gl/mathgl/mgl64"
	"github.com/go-vgo/robotgo"
	"github.com/justtaldevelops/go-hcaptcha/utils"
	"testing"
	"time"
)

func TestCurve(t *testing.T) {
	mousePosX, mousePosY := robotgo.GetMousePos()
	start := mgl64.Vec2{float64(mousePosX), float64(mousePosY)}
	end := mgl64.Vec2{1187, 719}

	humanCurve := NewCurve(start, end, &CurveOpts{})
	pause := utils.Between(5, 15)
	for _, point := range humanCurve.points {
		robotgo.MoveMouse(int(point.X()), int(point.Y()))
		time.Sleep(time.Duration(pause) * time.Millisecond)
	}
}
