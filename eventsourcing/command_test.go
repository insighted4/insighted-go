package eventsourcing

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCommandModel_AggregateID(t *testing.T) {
	m := CommandModel{ID: "abc"}
	assert.Equal(t, m.ID, m.AggregateID())
}
