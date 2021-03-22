package main

import (
	"fmt"
	"github.com/justtaldevelops/hcaptcha-solver-go"
	"os"
	"time"
)

func main() {
	// In order to use Vision API, you need to set this environment variable.
	err := os.Setenv("GOOGLE_APPLICATION_CREDENTIALS", "C:\\Users\\JustTal\\VisionAPI\\Project.json")
	if err != nil {
		panic(err)
	}

	// Create the solver with the default options.
	s, err := hcaptcha.NewSolver("minecraftpocket-servers.com", hcaptcha.SolverOptions{
		SiteKey: "e6b7bb01-42ff-4114-9245-3d2b7842ed92",
	})
	if err != nil {
		panic(err)
	}
	defer s.Close()

	// We provide a deadline that the solver must have the solution done by.
	// If the deadline is not reached, an error is sent instead of the solution.
	solution, err := s.Solve(time.Now().Add(3 * time.Minute))
	if err != nil {
		panic(err)
	}
	// P0_eyJ0eXAiOiJKV1Q...
	fmt.Println(solution)
}
