package eventsourcing

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestEvent(t *testing.T) {
	m := Model{
		ID:      "abc",
		Version: 123,
		At:      time.Now(),
	}

	assert.Equal(t, m.ID, m.AggregateID())
	assert.Equal(t, m.Version, m.EventVersion())
	assert.Equal(t, m.At, m.EventAt())
}

type Custom struct {
	Model
}

func (c Custom) EventType() string {
	return "blah"
}

func TestEventType(t *testing.T) {
	m := Custom{}
	eventType, _ := EventType(m)
	assert.Equal(t, "blah", eventType)
}
