package utils

import (
	"math/rand"
	"time"
)

// sRand is a random number generator for hCaptchas.
var sRand = rand.New(rand.NewSource(time.Now().UnixNano()))

// Chance returns true if the given chance is greater than the random number.
func Chance(chance float64) bool {
	return sRand.Float64() < chance
}

// Between returns a number between two numbers.
func Between(min, max int) int {
	return sRand.Intn(max-min) + min
}
