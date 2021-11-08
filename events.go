package hcaptcha

import (
	"github.com/iancoleman/orderedmap"
	"github.com/justtaldevelops/go-hcaptcha/agents"
	"github.com/justtaldevelops/go-hcaptcha/screen"
)

// Event is a movement event.
type Event struct {
	screen.Point
	// Type represents the type of event. For example, a "mouse up" event would be "mu".
	Type string
	// Timestamp is the time the event was recorded.
	Timestamp int64
}

// EventRecorder helps record and retrieve movement events.
type EventRecorder struct {
	recording   bool
	agent       agents.Agent
	manifest    *orderedmap.OrderedMap
	timeBuffers map[string]*EventContainer
}

// NewEventRecorder creates a new EventRecorder.
func NewEventRecorder(agent agents.Agent) *EventRecorder {
	return &EventRecorder{
		agent:       agent,
		manifest:    orderedmap.New(),
		timeBuffers: make(map[string]*EventContainer),
	}
}

// Record records a new event.
func (e *EventRecorder) Record() {
	e.manifest.Set("st", e.agent.Unix())
	e.recording = true
}

// Data records the events to the manifest and returns it.
func (e *EventRecorder) Data() *orderedmap.OrderedMap {
	for event, container := range e.timeBuffers {
		e.manifest.Set(event, container.Data())
		e.manifest.Set(event+"-mp", container.MeanPeriod())
	}
	return e.manifest
}

// SetData sets data in the manifest of the EventRecorder.
func (e *EventRecorder) SetData(name string, value interface{}) {
	e.manifest.Set(name, value)
}

// RecordEvent records a new event.
func (e *EventRecorder) RecordEvent(event Event) {
	if !e.recording {
		return
	}

	if _, ok := e.timeBuffers[event.Type]; !ok {
		e.timeBuffers[event.Type] = NewEventContainer(e.agent, 16, 15e3)
	}
	e.timeBuffers[event.Type].Push(event)
}

// EventContainer is a container for event data.
type EventContainer struct {
	agent             agents.Agent
	period, interval  int64
	date              []int64
	data              [][]int64
	previousTimestamp int64
	meanPeriod        int64
	meanCounter       int64
}

// NewEventContainer creates a new EventContainer.
func NewEventContainer(agent agents.Agent, period, interval int64) *EventContainer {
	return &EventContainer{
		agent:    agent,
		period:   period,
		interval: interval,
	}
}

// MeanPeriod returns the mean period of the event container.
func (e *EventContainer) MeanPeriod() int64 {
	return e.meanPeriod
}

// Data returns the data of the event container.
func (e *EventContainer) Data() [][]int64 {
	e.cleanStaleData()
	return e.data
}

// Push adds a new event to the event container.
func (e *EventContainer) Push(event Event) {
	e.cleanStaleData()

	notFirst := len(e.date) > 0

	var timestamp int64
	if notFirst {
		timestamp = e.date[len(e.date)-1]
	}

	if event.Timestamp-timestamp >= e.period {
		e.date = append(e.date, event.Timestamp)
		e.data = append(e.data, []int64{int64(event.Point.X), int64(event.Point.Y), event.Timestamp})

		if notFirst {
			delta := event.Timestamp - e.previousTimestamp
			e.meanPeriod = (e.meanPeriod*e.meanCounter + delta) / (e.meanCounter + 1)
			e.meanCounter++
		}
	}

	e.previousTimestamp = event.Timestamp
}

// cleanStaleData removes stale data from the event container.
func (e *EventContainer) cleanStaleData() {
	date := e.agent.Unix()
	t := len(e.date) - 1

	for t >= 0 {
		if date-e.date[t] >= e.interval {
			e.date = e.date[:t+1]
			e.date = e.date[:t+1]
			break
		}

		t -= 1
	}
}
