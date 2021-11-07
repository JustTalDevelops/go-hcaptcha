package algorithm

import (
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"strings"
)

const (
	// atob is a part of the full HSL generation script which converts base64 to binary.
	atob = "function atob(r){return Buffer.from(r,\"base64\").toString(\"binary\")}"
	// run is a part of the full HSL generation script which runs the script.
	run = "hsl(process.argv[2]).then(function(r){console.log(r)})"
)

// HSL is one of a few proof algorithms for hCaptcha services.
type HSL struct{}

// Encode ...
func (h *HSL) Encode() string {
	return "hsl"
}

// Prove ...
func (h *HSL) Prove(request string) (string, error) {
	f, err := ioutil.TempFile(os.TempDir(), "hsl.*.js")
	if err != nil {
		panic(err)
	}
	defer func() {
		_ = f.Close()
		_ = os.Remove(f.Name())
	}()

	_, err = f.WriteString(atob + script("hsl.js") + run)
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
