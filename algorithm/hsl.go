package algorithm

// HSL is one of a few proof algorithms for hCaptcha services.
type HSL struct{}

// Encode ...
func (h *HSL) Encode() string {
	return "hsl"
}

// Prove ...
func (h *HSL) Prove(request string) string {
	panic("not implemented")
}
