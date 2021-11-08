package agents

import "github.com/iancoleman/orderedmap"

// Agent is used to generate user-agent specific data.
type Agent interface {
	// UserAgent returns the user-agent string for the agent.
	UserAgent() string
	// ScreenProperties returns the screen properties of the agent.
	ScreenProperties() *orderedmap.OrderedMap
	// NavigatorProperties returns the navigator properties of the agent.
	NavigatorProperties() *orderedmap.OrderedMap

	// Unix returns the current timestamp with any added offsets.
	Unix() int64
	// OffsetUnix offsets the Unix timestamp with the given offset.
	OffsetUnix(offset int64)
	// ResetUnix resets the Unix timestamp with offsets to the current time.
	ResetUnix()
}

// Compile time check to make sure that the Agent interface is implemented by Chrome.
var _ Agent = (*Chrome)(nil)
