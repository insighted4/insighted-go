package eventsourcing

import "fmt"

const (
	//AggregateNil      = "AggregateNil"
	//DuplicateID       = "DuplicateID"
	//DuplicateVersion  = "DuplicateVersion"
	//DuplicateAt       = "DuplicateAt"
	//DuplicateType     = "DuplicateType"
	//InvalidID         = "InvalidID"
	//InvalidAt         = "InvalidAt"
	//InvalidVersion    = "InvalidVersion"

	// ErrorInvalidArgument is returned when the caller passed an incorrect value
	ErrorInvalidArgument = "invalid argument"

	// ErrorInvalidEncoding is returned when the Serializer cannot marshal the event
	ErrorInvalidEncoding = "invalid encoding"

	// ErrorUnboundEventType when the Serializer cannot unmarshal the serialized event
	ErrorUnboundEventType = "unbound event type"

	// ErrorAggregateNotFound will be returned when attempting to load an aggregateID
	// that does not exist in the Store
	ErrorAggregateNotFound = "aggregate not found"

	// ErrorAggregateNotSaved will be returned when attempting to save an aggregate in the Store
	ErrorAggregateNotSaved = "aggregate not saved"

	// ErrorUnhandledCommand occurs when the command handler is unable to handle a command
	ErrorUnhandledCommand = "unhandled command"

	// ErrorUnhandledEvent occurs when the Aggregate is unable to handle an event and returns
	// a non-nill err
	ErrorUnhandledEvent = "unhandled event"
)

// Error provides a standardized error interface for eventsource
type Error interface {
	error

	// Returns the original error if one was set.  Nil is returned if not set.
	Cause() error

	// Returns the short phrase depicting the classification of the error.
	Code() string

	// Returns the error details message.
	Message() string
}

type baseErr struct {
	cause   error
	code    string
	message string
}

func (b *baseErr) Cause() error    { return b.cause }
func (b *baseErr) Code() string    { return b.code }
func (b *baseErr) Message() string { return b.message }
func (b *baseErr) Error() string   { return fmt.Sprintf("%s: %s", b.message, b.cause.Error()) }
func (b *baseErr) String() string  { return b.Error() }

// NewError generates the common error structure
func NewError(err error, code, format string, args ...interface{}) Error {
	return &baseErr{
		code:    code,
		message: fmt.Sprintf(format, args...),
		cause:   err,
	}
}

// ErrHasCode returns true if any error in the cause chain has the specified code
func ErrHasCode(err error, code string) bool {
	if err == nil {
		return false
	}

	v, ok := err.(Error)
	if !ok {
		return false
	}

	if v.Code() == code {
		return true
	}

	if cause := v.Cause(); cause != nil {
		return ErrHasCode(cause, code)
	}

	return false
}

// IsNotFound returns true if the issue as the aggregate was not found
func IsNotFound(err error) bool {
	for err != nil {
		if err == nil {
			return false
		}

		v, ok := err.(Error)
		if !ok {
			return false
		}

		if v.Code() == ErrorAggregateNotFound {
			return true
		}

		err = v.Cause()
	}

	return false
}
