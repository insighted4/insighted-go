package eventsourcing

import (
	"context"
	"reflect"
	"time"

	"github.com/sirupsen/logrus"
)

// Aggregate represents the aggregate root in the domain driven design sense.
// It represents the current state of the domain object and can be thought of
// as a left fold over events.
type Aggregate interface {
	// On will be called for each event; returns err if the event could not be
	// applied
	On(event Event) Error
}

// Repository provides the primary abstraction to saving and loading events
type Repository struct {
	prototype  reflect.Type
	store      Store
	serializer Serializer
	observers  []func(Event)
	logger     logrus.FieldLogger
}

// New returns a new instance of the aggregate
func (r *Repository) New() Aggregate {
	return reflect.New(r.prototype).Interface().(Aggregate)
}

// Save persists the events into the underlying Store
func (r *Repository) Save(ctx context.Context, events ...Event) Error {
	if len(events) == 0 {
		return nil
	}
	aggregateID := events[0].AggregateID()

	history := make(History, 0, len(events))
	for _, event := range events {
		record, err := r.serializer.MarshalEvent(event)
		if err != nil {
			return err
		}

		history = append(history, record)
	}

	return r.store.Save(ctx, aggregateID, history...)
}

// Load retrieves the specified aggregate from the underlying store
func (r *Repository) Load(ctx context.Context, aggregateID string) (Aggregate, Error) {
	v, _, err := r.loadVersion(ctx, aggregateID, 0)
	return v, err
}

// LoadVersion retrieves the specified aggregate from the underlying store at a particular version.
func (r *Repository) LoadVersion(ctx context.Context, aggregateID string, version int) (Aggregate, error) {
	v, _, err := r.loadVersion(ctx, aggregateID, version)
	return v, err
}

// LoadTime retrieves the specified aggregate from the underlying store as of some point in time.
func (r *Repository) LoadTime(ctx context.Context, aggregateID string, endTime time.Time) (Aggregate, error) {
	v, _, err := r.loadTime(ctx, aggregateID, endTime)
	return v, err
}

// LoadVersion loads the specified aggregate from the store and returns both the Aggregate and the
// current version number of the aggregate
func (r *Repository) loadVersion(ctx context.Context, aggregateID string, version int) (Aggregate, int, Error) {
	history, err := r.store.Load(ctx, aggregateID, 0, version)
	if err != nil {
		return nil, 0, err
	}

	entryCount := len(history)
	if entryCount == 0 {
		return nil, 0, NewError(nil, ErrorAggregateNotFound, "unable to load %v, %v", r.New(), aggregateID)
	}

	r.logger.Infof("Loaded %v event(s) for aggregate id, %v", entryCount, aggregateID)
	aggregate := r.New()

	version = 0
	for _, record := range history {
		event, err := r.serializer.UnmarshalEvent(record)
		if err != nil {
			return nil, 0, err
		}

		err = aggregate.On(event)
		if err != nil {
			eventType, _ := EventType(event)
			return nil, 0, NewError(err, ErrorUnhandledEvent, "aggregate was unable to handle event, %v", eventType)
		}

		version = event.EventVersion()
	}

	return aggregate, version, nil
}

// loadTime loads the specified aggregate from the store at some point in time and returns
// both the Aggregate and the current version number of the aggregate.
func (r *Repository) loadTime(ctx context.Context, aggregateID string, endTime time.Time) (Aggregate, int, error) {
	history, err := r.store.Load(ctx, aggregateID, 0, 0)
	if err != nil {
		return nil, 0, err
	}

	entryCount := len(history)
	if entryCount == 0 {
		return nil, 0, NewError(nil, ErrorAggregateNotFound, "unable to load %v, %v", r.New(), aggregateID)
	}

	r.logger.Infof("Loaded %v event(s) for aggregate id, %v", entryCount, aggregateID)
	aggregate := r.New()

	version := 0
	for _, record := range history {
		event, err := r.serializer.UnmarshalEvent(record)
		if err != nil {
			return nil, 0, err
		}

		// Stop if after the end time.
		if event.EventAt().After(endTime) {
			break
		}

		err = aggregate.On(event)
		if err != nil {
			eventType, _ := EventType(event)
			return nil, 0, NewError(err, ErrorUnhandledEvent, "aggregate was unable to handle event, %v", eventType)
		}

		version = event.EventVersion()
	}

	return aggregate, version, nil
}

// Apply executes the command specified and returns the current version of the aggregate
func (r *Repository) Apply(ctx context.Context, command Command) (int, Error) {
	if command == nil {
		return 0, NewError(nil, ErrorInvalidArgument, "command provided to Repository.Dispatch may not be nil")
	}
	aggregateID := command.AggregateID()
	if aggregateID == "" {
		return 0, NewError(nil, ErrorInvalidArgument, "command provided to Repository.Dispatch may not contain a blank aggregate ID")
	}

	aggregate, version, err := r.loadVersion(ctx, aggregateID, 0)
	if err != nil {
		aggregate = r.New()
	}

	h, ok := aggregate.(CommandHandler)
	if !ok {
		return 0, NewError(nil, ErrorInvalidArgument, "aggregate %v, does not implement CommandHandler", aggregate)
	}
	events, err := h.Apply(ctx, command)
	if err != nil {
		return 0, err
	}

	err = r.Save(ctx, events...)
	if err != nil {
		return 0, err
	}

	if v := len(events); v > 0 {
		version = events[v-1].EventVersion()
	}

	// publish events to observers
	if r.observers != nil {
		for _, event := range events {
			for _, observer := range r.observers {
				observer(event)
			}
		}
	}

	return version, nil
}

// Store returns the underlying Store
func (r *Repository) Store() Store {
	return r.store
}

// Serializer returns the underlying serializer
func (r *Repository) Serializer() Serializer {
	return r.serializer
}

// NewRepository creates a new Repository using the JSON serializer and In-Memory store.
// Observers should invoke very short lived operations as calls will block until the observer is finished.
func NewRepository(prototype Aggregate, store Store, serializer Serializer, logger logrus.FieldLogger, observers ...func(event Event)) *Repository {
	t := reflect.TypeOf(prototype)
	if t.Kind() == reflect.Ptr {
		t = t.Elem()
	}

	r := &Repository{
		prototype:  t,
		store:      store,
		serializer: serializer,
		logger:     logger,
		observers:  observers,
	}

	return r
}
