package hcaptcha

import (
	"testing"
	"time"
)

// TestCaptcha ...
func TestCaptcha(t *testing.T) {
	for {
		c, err := NewChallenge(
			"https://accounts.hcaptcha.com/demo",
			"a5f74b19-9e45-40e0-b45d-47ff91b7a6c2",
			ChallengeOptions{
				Timeout: 10 * time.Second,
			},
		)
		if err != nil {
			panic(err)
		}
		err = c.Solve(&GuessSolver{})
		if err != nil {
			c.Logger().Debugf("Error from hCaptcha API: %s", err)
			continue
		}
		c.Logger().Info(c.Token())
		break
	}
}
