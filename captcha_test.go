package hcaptcha

import (
	"fmt"
	"testing"
)

// TestCaptcha ...
func TestCaptcha(t *testing.T) {
	for {
		c, err := NewChallenge("https://minecraftpocket-servers.com/server/41256/vote/", "e6b7bb01-42ff-4114-9245-3d2b7842ed92")
		if err != nil {
			panic(err)
		}
		err = c.Solve(&GuessSolver{})
		if err != nil {
			fmt.Println(err)
			continue
		}
		fmt.Println(c.Token())
		break
	}
}
