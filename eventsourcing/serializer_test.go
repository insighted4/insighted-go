package eventsourcing

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

type EntitySetName struct {
	Model
	Name string
}

func TestNewJSONSerializer(t *testing.T) {
	event := EntitySetName{
		Model: Model{
			ID:      "123",
			Version: 456,
		},
		Name: "blah",
	}

	serializer := NewJSONSerializer(event)
	record, err := serializer.MarshalEvent(event)
	assert.Nil(t, err)
	assert.NotNil(t, record)

	v, err := serializer.UnmarshalEvent(record)
	assert.Nil(t, err)

	found, ok := v.(*EntitySetName)
	assert.True(t, ok)
	assert.Equal(t, &event, found)
}

func TestJSONSerializer_MarshalAll(t *testing.T) {
	event := EntitySetName{
		Model: Model{
			ID:      "123",
			Version: 456,
		},
		Name: "blah",
	}

	serializer := NewJSONSerializer(event)
	history, err := serializer.MarshalAll(event)
	assert.Nil(t, err)
	assert.NotNil(t, history)

	v, err := serializer.UnmarshalEvent(history[0])
	assert.Nil(t, err)

	found, ok := v.(*EntitySetName)
	assert.True(t, ok)
	assert.Equal(t, &event, found)
}
