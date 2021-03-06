package eventsourcing

import (
	"encoding/json"
	"reflect"
)

// Serializer converts between Events and Records
type Serializer interface {
	// MarshalEvent converts an Event to a Record
	MarshalEvent(event Event) (Record, Error)

	// UnmarshalEvent converts an Event backed into a Record
	UnmarshalEvent(record Record) (Event, Error)
}

type jsonEvent struct {
	Type string          `json:"t"`
	Data json.RawMessage `json:"d"`
}

// JSONSerializer provides a simple serializer implementation
type JSONSerializer struct {
	eventTypes map[string]reflect.Type
}

// Bind registers the specified events with the serializer; may be called more than once
func (j *JSONSerializer) Bind(events ...Event) {
	for _, event := range events {
		eventType, t := EventType(event)
		j.eventTypes[eventType] = t
	}
}

// MarshalEvent converts an event into its persistent type, Record
func (j *JSONSerializer) MarshalEvent(v Event) (Record, Error) {
	eventType, _ := EventType(v)

	data, err := json.Marshal(v)
	if err != nil {
		return Record{}, NewError(err, ErrorInvalidEncoding, "unable to encode data")
	}

	data, err = json.Marshal(jsonEvent{
		Type: eventType,
		Data: json.RawMessage(data),
	})
	if err != nil {
		return Record{}, NewError(err, ErrorInvalidEncoding, "unable to encode event")
	}

	return Record{
		Version: v.EventVersion(),
		Data:    data,
	}, nil
}

// UnmarshalEvent converts the persistent type, Record, into an Event instance
func (j *JSONSerializer) UnmarshalEvent(record Record) (Event, Error) {
	wrapper := jsonEvent{}
	err := json.Unmarshal(record.Data, &wrapper)
	if err != nil {
		return nil, NewError(err, ErrorInvalidEncoding, "unable to unmarshal event")
	}

	t, ok := j.eventTypes[wrapper.Type]
	if !ok {
		return nil, NewError(err, ErrorUnboundEventType, "unbound event type, %v", wrapper.Type)
	}

	v := reflect.New(t).Interface()
	err = json.Unmarshal(wrapper.Data, v)
	if err != nil {
		return nil, NewError(err, ErrorInvalidEncoding, "unable to unmarshal event data into %#v", v)
	}

	return v.(Event), nil
}

// MarshalAll is a utility that marshals all the events provided into a History object
func (j *JSONSerializer) MarshalAll(events ...Event) (History, Error) {
	history := make(History, 0, len(events))

	for _, event := range events {
		record, err := j.MarshalEvent(event)
		if err != nil {
			return nil, err
		}
		history = append(history, record)
	}

	return history, nil
}

// NewJSONSerializer constructs a new JSONSerializer and populates it with the specified events.
// Bind may be subsequently called to add more events.
func NewJSONSerializer(events ...Event) *JSONSerializer {
	serializer := &JSONSerializer{
		eventTypes: map[string]reflect.Type{},
	}
	serializer.Bind(events...)

	return serializer
}
