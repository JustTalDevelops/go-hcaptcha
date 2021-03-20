package main

import (
	"fmt"
	"github.com/justtaldevelops/hcaptcha-solver-go"
	"time"
)

func main() {
	s, err := hcaptcha.NewSolver("minecraftpocket-servers.com")
	if err != nil {
		panic(err)
	}
	solution, err := s.Solve(time.Now().Add(1 * time.Minute))
	if err != nil {
		panic(err)
	}
	// F0_eyJ0eXAiOiJKV1Q...
	fmt.Println(solution)
}
