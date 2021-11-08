package hcaptcha

import (
	"github.com/justtaldevelops/go-hcaptcha/utils"
)

// Solver is an interface to solve an hCaptcha task.
type Solver interface {
	// Solve solves the hCaptcha task using the category, question, and the task. If it was successful,
	// it returns true, and in all other cases, it returns false.
	Solve(category, question string, task Task) bool
}

// GuessSolver solves hCaptcha tasks by guessing the solution.
type GuessSolver struct{}

// Solve ...
func (s *GuessSolver) Solve(_, _ string, _ Task) bool {
	return utils.Chance(0.5)
}

// Compile-time check that GuessSolver implements Solver.
var _ Solver = (*GuessSolver)(nil)
