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
	defer s.Close()
	// We provide a deadline that the solver must have the solution done by.
	// If the deadline is not reached, an error is sent instead of the solution.
	solution, err := s.Solve(time.Now().Add(1 * time.Minute))
	if err != nil {
		panic(err)
	}
	// F0_eyJ0eXAiOiJKV1Q...
	fmt.Println(solution)
}
