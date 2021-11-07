package agents

// Agent is used to generate user-agent specific data.
type Agent interface {
	// ScreenProperties returns the screen properties of the agent.
	ScreenProperties() map[string]interface{}
	// NavigatorProperties returns the navigator properties of the agent.
	NavigatorProperties() map[string]interface{}

	// Unix returns the current timestamp with any added offsets.
	Unix(asMilliseconds bool) int64
	// OffsetUnix offsets the Unix timestamp with the given offset.
	OffsetUnix(offset int64)
	// ResetUnix resets the Unix timestamp with offsets to the current time.
	ResetUnix()
}

// screenSize is the size of an agents screen.
type screenSize [2]int

// Default returns the default agent.
func Default() Agent {
	// Chrome only for now, sorry!
	return &Chrome{}
}

// Compile time check to make sure that the Agent interface is implemented by Chrome.
var _ Agent = (*Chrome)(nil)
