package eventsourcing

import (
	"errors"
	"fmt"
	"io"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewError(t *testing.T) {
	err := NewError(io.EOF, "code", "hello %v", "world")
	assert.NotNil(t, err)

	v, ok := err.(Error)
	assert.True(t, ok)
	assert.Equal(t, io.EOF, v.Cause())
	assert.Equal(t, "code", v.Code())
	assert.Equal(t, "hello world", v.Message())
	assert.Equal(t, "hello world: EOF", v.Error())

	s, ok := err.(fmt.Stringer)
	assert.True(t, ok)
	assert.Equal(t, v.Error(), s.String())
}

func TestIsNotFound(t *testing.T) {
	testCases := map[string]struct {
		Err        error
		IsNotFound bool
	}{
		"nil": {
			Err:        nil,
			IsNotFound: false,
		},
		"Error": {
			Err:        NewError(nil, ErrorAggregateNotFound, "not found"),
			IsNotFound: true,
		},
		"nested Error": {
			Err: NewError(
				NewError(nil, ErrorAggregateNotFound, "not found"),
				ErrorUnboundEventType,
				"not found",
			),
			IsNotFound: true,
		},
	}

	for label, tc := range testCases {
		t.Run(label, func(t *testing.T) {
			assert.Equal(t, tc.IsNotFound, IsNotFound(tc.Err))
		})
	}
}

func TestErrHasCode(t *testing.T) {
	code := "code"

	testCases := map[string]struct {
		Err        error
		ErrHasCode bool
	}{
		"simple": {
			Err:        NewError(nil, code, "blah"),
			ErrHasCode: true,
		},
		"nope": {
			Err:        errors.New("blah"),
			ErrHasCode: false,
		},
		"nested": {
			Err:        NewError(NewError(nil, code, "blah"), "blah", "blah"),
			ErrHasCode: true,
		},
	}

	for label, tc := range testCases {
		t.Run(label, func(t *testing.T) {
			assert.Equal(t, tc.ErrHasCode, ErrHasCode(tc.Err, code))
		})
	}
}
