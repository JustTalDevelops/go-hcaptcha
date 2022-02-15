package utils

import (
	"gopkg.in/square/go-jose.v2/jwt"
)

// ParseJWT parses a JWT using go-jose and returns the payload as a map.
func ParseJWT(raw string) map[string]interface{} {
	tok, err := jwt.ParseSigned(raw)
	if err != nil {
		panic(err)
	}

	var out map[string]interface{}
	if err = tok.UnsafeClaimsWithoutVerification(&out); err != nil {
		panic(err)
	}

	return out
}
