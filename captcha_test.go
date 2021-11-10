package hcaptcha

import (
	"testing"
)

// TestCaptcha ...
func TestCaptcha(t *testing.T) {
	for {
		c, err := NewChallenge("https://minecraftpocket-servers.com/server/41256/vote/", "e6b7bb01-42ff-4114-9245-3d2b7842ed92")
		if err != nil {
			panic(err)
		}
		err = c.Solve(&YOLOSolver{Log: c.Logger()})
		if err != nil {
			c.Logger().Debug(err)
			continue
		}
		c.Logger().Info(c.Token())
	}
}
