package errors

import (
	"bytes"
	"encoding/json"
	"errors"
	"strings"
)

// Error defines a standard application error.
type Error struct {
	// Error classification for the application.
	Kind Kind `json:"kind"`

	// Human-readable message.
	Message string `json:"message"`

	// Wrapped underlying error.
	WrappedErr error `json:"wrapped_err,omitempty"`
}

// Error returns the string representation of the error message.
func (e *Error) Error() string {
	var buf bytes.Buffer
	json.NewEncoder(&buf).Encode(e)
	return buf.String()
}

// Kind defines the kind or class of an error.
type Kind uint8

// Transport agnostic error "kinds"
const (
	Other    Kind = iota // Unclassified error
	Internal             // Internal error
	Conflict             // Conflict when an entity already exists
	Invalid              // Invalid input, validation error etc
	NotFound             // Entity does not exist
)

func (k Kind) String() string {
	switch k {
	case Other:
		return "unclassified error"
	case Internal:
		return "internal error"
	case Invalid:
		return "invalid input"
	case NotFound:
		return "entity not found"
	default:
		return "unknown error kind"
	}
}

func (k Kind) MarshalJSON() ([]byte, error) {
	return json.Marshal(k.String())
}

// E is a helper function which constructs an `*Error`
// You can pass it Kind, error (Err) or string (Message) in any order and it'll construct it.
func E(args ...interface{}) error {
	e := &Error{}
	for _, arg := range args {
		switch arg := arg.(type) {
		case Kind:
			e.Kind = arg
		case error:
			e.WrappedErr = arg
		case string:
			e.Message = arg
		}
	}
	return e
}

type MultiErr []error

func (m MultiErr) Error() string {
	msg := &strings.Builder{}
	msg.WriteRune('[')
	for i, err := range m {
		if i > 0 {
			msg.WriteRune(',')
		}
		msg.WriteRune('"')
		msg.WriteString(err.Error())
		msg.WriteRune('"')
	}
	msg.WriteRune(']')
	return msg.String()
}

var (
	As = errors.As
	Is = errors.Is
)
