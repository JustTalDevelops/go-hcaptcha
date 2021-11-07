package algorithm

// algorithms is a map of all available algorithms.
var algorithms = make(map[string]Algorithm)

// init registers all available algorithms.
func init() {
	registerAlgorithm(&HSL{})
	registerAlgorithm(&HSW{})
}

// findAlgorithm returns the algorithm with the given name.
func findAlgorithm(name string) Algorithm {
	return algorithms[name]
}

// registerAlgorithm registers an algorithm.
func registerAlgorithm(algorithm Algorithm) {
	algorithms[algorithm.Encode()] = algorithm
}
