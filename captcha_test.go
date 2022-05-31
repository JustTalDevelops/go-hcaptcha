package hcaptcha

import (
	"testing"
	"time"
)

// TestCaptcha ...
func TestCaptcha(t *testing.T) {
	for {
		c, err := NewChallenge(
			"https://democaptcha.com/demo-form-eng/hcaptcha.html",
			"51829642-2cda-4b09-896c-594f89d700cc",
			ChallengeOptions{
				Timeout: 10 * time.Second,
			},
		)
		if err != nil {
			panic(err)
		}
		err = c.Solve(&GuessSolver{})
		if err != nil {
			c.Logger().Debug(err)
			continue
		}
		c.Logger().Info(c.Token())
		break
	}
}
