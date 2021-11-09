package hcaptcha

import (
	"testing"
)

// TestCaptcha ...
func TestCaptcha(t *testing.T) {
	c, err := NewChallenge("https://minecraftpocket-servers.com/server/41256/vote/", "e6b7bb01-42ff-4114-9245-3d2b7842ed92")
	if err != nil {
		panic(err)
	}
	err = c.Solve(&YOLOSolver{Log: c.log})
	if err != nil {
		c.log.Panic(err)
	}
	c.log.Info(c.Token())
}
