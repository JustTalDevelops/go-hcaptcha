package hcaptcha

import (
	"github.com/justtaldevelops/go-hcaptcha/utils"
)

// Solver is an interface to solve hCaptcha tasks.
type Solver interface {
	// Solve solves the hCaptcha tasks using the category, question, and the task. If it was successful,
	// it returns true, and in all other cases, it returns false.
	Solve(category, question string, tasks []Task) []Task
}

// GuessSolver solves hCaptcha tasks by guessing the solution.
type GuessSolver struct{}

// Solve ...
func (s *GuessSolver) Solve(_, _ string, tasks []Task) (answers []Task) {
	for _, task := range tasks {
		if utils.Chance(0.5) {
			answers = append(answers, task)
		}
	}
	return answers
}
