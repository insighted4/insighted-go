package eventsourcing

import (
	"context"
	"sort"
	"sync"
)

// Record provides the serialized representation of the event
type Record struct {
	// Version contains the version associated with the serialized event
	Version int

	// Data contains the event in serialized form
	Data []byte
}

// History represents
type History []Record

// Len implements sort.Interface
func (h History) Len() int {
	return len(h)
}

// Swap implements sort.Interface
func (h History) Swap(i, j int) {
	h[i], h[j] = h[j], h[i]
}

// Less implements sort.Interface
func (h History) Less(i, j int) bool {
	return h[i].Version < h[j].Version
}

// Store provides an abstraction for the Repository to save data
type Store interface {
	// Save the provided serialized records to the store
	Save(ctx context.Context, aggregateID string, records ...Record) Error

	// Load the history of events up to the version specified.
	// When toVersion is 0, all events will be loaded.
	// To start at the beginning, fromVersion should be set to 0
	Load(ctx context.Context, aggregateID string, fromVersion, toVersion int) (History, Error)
}

// MemStore provides an in-MemStore implementation of Service
type MemStore struct {
	mux        *sync.Mutex
	eventsByID map[string]History
}

// Save implements the Store interface and saves records, serialized events, in-MemStore
func (m *MemStore) Save(ctx context.Context, aggregateID string, records ...Record) Error {
	if _, ok := m.eventsByID[aggregateID]; !ok {
		m.eventsByID[aggregateID] = History{}
	}

	history := append(m.eventsByID[aggregateID], records...)
	sort.Sort(history)
	m.eventsByID[aggregateID] = history

	return nil
}

// Load implements the Store interface and retrieve events from in-MemStore
func (m *MemStore) Load(ctx context.Context, aggregateID string, fromVersion, toVersion int) (History, Error) {
	all, ok := m.eventsByID[aggregateID]
	if !ok {
		return nil, NewError(nil, ErrorAggregateNotFound, "no aggregate found with id, %v", aggregateID)
	}

	history := make(History, 0, len(all))
	if len(all) > 0 {
		for _, record := range all {
			if v := record.Version; v >= fromVersion && (toVersion == 0 || v <= toVersion) {
				history = append(history, record)
			}
		}
	}

	return all, nil
}

// NewMemStore returns a in-MemStore backed store
func NewMemStore() *MemStore {
	return &MemStore{
		mux:        &sync.Mutex{},
		eventsByID: map[string]History{},
	}
}
