package algorithm

import (
	"github.com/justtaldevelops/go-hcaptcha/utils"
	"io/ioutil"
	"net/http"
)

// Algorithm is the algorithm used to provide proof.
type Algorithm interface {
	// Encode encodes the algorithm name as a string.
	Encode() string
	// Prove returns proof (N) of the request (C).
	Prove(request string) (string, error)
}

// Proof is the full proof of a request.
type Proof struct {
	// Algorithm is the algorithm used to provide proof.
	Algorithm Algorithm
	// Request is the original request (C) which was used to calculate the proof (N).
	Request string
	// Proof is the calculated proof (N).
	Proof string
}

// Compile time checks to make sure HSL and HSW implement Algorithm.
var _, _ Algorithm = (*HSL)(nil), (*HSW)(nil)

// Solve solves for N given a specific algorithm and request. It returns the full proof.
func Solve(algorithm, request string) (Proof, error) {
	algo := findAlgorithm(algorithm)
	proof, err := algo.Prove(request)
	if err != nil {
		return Proof{}, err
	}

	return Proof{Algorithm: algo, Request: `{"type":"` + algo.Encode() + `","req":"` + request + `"}`, Proof: proof}, nil
}

// script gets the script of the algorithm from hCaptcha.
func script(script string) string {
	resp, err := http.Get("https://newassets.hcaptcha.com/c/" + utils.AssetVersion + "/" + script)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	b, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}

	return string(b)
}
