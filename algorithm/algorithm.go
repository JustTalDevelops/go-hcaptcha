package algorithm

// Algorithm is the algorithm used to provide proof.
type Algorithm interface {
	// Encode encodes the algorithm name as a string.
	Encode() string
	// Initialize is called on the algorithms' initialization.
	Initialize()
	// Prove returns proof (N) of the request (C).
	Prove(request string) string
}

// Compile time checks to make sure HSL and HSW implement Algorithm.
var _, _ Algorithm = (*HSL)(nil), (*HSW)(nil)

// Solve solves for N given a specific algorithm and request.
func Solve(algorithm, request string) string {
	return findAlgorithm(algorithm).Prove(request)
}
