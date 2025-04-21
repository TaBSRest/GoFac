package TaBSError

import (
	"errors"
	"fmt"
	"runtime"
	"strings"
	"slices"
)

type TaBSError struct {
	location   string
	message    string
	children []error
}

// New creates a TaBSError, capturing a caller location (skipping 2 frames).
func New(message string) *TaBSError {
	var innerErrors []error

	return &TaBSError{
		location: getCallerLocation(2),
		message:  message,
		children: innerErrors,
	}
}

// NewWithGeneric creates a TaBSError, adds generic type info in brackets,
// and captures the caller location.
func NewWithGeneric(genericName string, message string) *TaBSError {
	callerLocation := fmt.Sprintf("%s[%s]", getCallerLocation(2), genericName)

	var innerErrors []error

	return &TaBSError{
		location: callerLocation,
		message:  message,
		children: innerErrors,
	}
}

// NewWithLocation creates a TaBSError with a location, and
// it is used for some functions which are not able to track locations
func NewWithLocation(location string, message string) *TaBSError {
	var innerErrors []error

	return &TaBSError{
		location: location,
		message:  message,
		children: innerErrors,
	}
}

// Join attaches an inner error to the TaBSError (for chaining).
func (e *TaBSError) Join(child error) *TaBSError {
	copy := *e
	copy.children = append(slices.Clone(e.children), child)
	return &copy
}

func (e *TaBSError) JoinMultiple(children []error) *TaBSError {
	copy := *e
	copy.children = append(slices.Clone(e.children), children...)
	return &copy
}

func (e *TaBSError) JoinString(innerError string) *TaBSError {
	copy := *e
	copy.children = append(slices.Clone(e.children), errors.New(innerError))
	return &copy
}

// Error formats the TaBSError, optionally including the inner error message.
func (e *TaBSError) Error() string {
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%s: %s\n\t", e.location, e.message))

	if len(e.children) > 0 {
		for _, child := range e.children {
			sb.WriteString("with inner error: ")
			if child != nil {
				lines := strings.Split(child.Error(), "\n")
				for _, line := range lines {
					sb.WriteString(line)
					sb.WriteString("\n\t")
				}
			}
		}
	}

	return strings.TrimRight(sb.String(), "\n\t")
}

// getCallerLocation finds the caller location using runtime.Caller.
func getCallerLocation(skipFrames int) string {
	pc, _, _, ok := runtime.Caller(skipFrames)
	if !ok {
		return "<unknown location>"
	}

	fullFuncName := runtime.FuncForPC(pc).Name()
	fullFuncName = strings.SplitN(fullFuncName, "[", 2)[0]

	parts := strings.SplitN(fullFuncName, ".", -1)

	location := strings.TrimPrefix(parts[1], "com/TaBSRest/")
	location = strings.SplitN(location, "(", 2)[0]

	return fmt.Sprintf("%s.%s", location, parts[len(parts)-1])
}


func (e *TaBSError) Unwrap() error {
	if len(e.children) == 0 {
		return e.children[0] // Pick the first child as the primary cause
	}
	return New("The error contains multilpe errors. It cannot unwrap to one error.")
}

