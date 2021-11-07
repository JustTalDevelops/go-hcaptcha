package algorithm

// HSW is one of a few proof algorithms for hCaptcha services.
type HSW struct{}

// Encode ...
func (h *HSW) Encode() string {
	return "hsw"
}

// Initialize ...
func (h *HSW) Initialize() {

}

// Prove ...
func (h HSW) Prove(request string) string {
	panic("not implemented")
}
