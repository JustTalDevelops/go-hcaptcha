package algorithm

import (
	"github.com/justtaldevelops/hcaptcha-solver-go/utils"
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

// Compile time checks to make sure HSL and HSW implement Algorithm.
var _, _ Algorithm = (*HSL)(nil), (*HSW)(nil)

// Solve solves for N given a specific algorithm and request.
func Solve(algorithm, request string) (string, error) {
	return findAlgorithm(algorithm).Prove(request)
}

// script gets the script of the algorithm from hCaptcha.
func script(script string) string {
	resp, err := http.Get("https://newassets.hcaptcha.com/c/" + utils.AssetVersion() + "/" + script)
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
