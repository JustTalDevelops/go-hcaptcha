package algorithm

import (
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"strings"
)

// HSL is one of a few proof algorithms for hCaptcha services.
type HSL struct {
	// script is the cached script file.
	script string
}

// Encode ...
func (h *HSL) Encode() string {
	return "hsl"
}

// Initialize ...
func (h *HSL) Initialize() {
	h.script = script("hsl.js")
}

// Prove ...
func (h *HSL) Prove(request string) (string, error) {
	f, err := ioutil.TempFile(os.TempDir(), "hsl.*.js")
	if err != nil {
		panic(err)
	}
	defer f.Close()
	defer os.Remove(f.Name())

	_, err = f.WriteString(h.script)
	if err != nil {
		panic(err)
	}

	cmd := exec.Command("node", f.Name(), request)
	cmd.Dir = os.TempDir()
	out, err := cmd.Output()
	if err != nil {
		return "", err
	}
	output := string(out)

	if strings.Contains(output, "Token is invalid.") {
		return "", fmt.Errorf(output)
	}
	return strings.TrimSpace(output), nil
}
