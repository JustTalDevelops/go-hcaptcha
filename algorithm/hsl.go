package algorithm

import (
	"fmt"
	"gopkg.in/square/go-jose.v2/jwt"
	"strings"
	"time"
)

// HSL is one of a few proof algorithms for hCaptcha services.
type HSL struct{}

// Encode ...
func (h *HSL) Encode() string {
	return "hsl"
}

// Prove ...
func (h *HSL) Prove(request string) (string, error) {
	tok, err := jwt.ParseSigned(request)
	if err != nil {
		panic(err)
	}

	var claims map[string]interface{}
	if err = tok.UnsafeClaimsWithoutVerification(&claims); err != nil {
		panic(err)
	}

	now := time.Now().UTC().Format("2006-01-02T15:04:05.999Z")
	now = now[:len(now)-5]
	now = strings.ReplaceAll(now, "-", "")
	now = strings.ReplaceAll(now, ":", "")
	now = strings.ReplaceAll(now, "T", "")

	return strings.Join([]string{
		"1",
		fmt.Sprint(int(claims["s"].(float64))),
		now,
		claims["d"].(string),
		"",
		"1", // TODO: Figure out if this is right?
	}, ":"), nil
}
