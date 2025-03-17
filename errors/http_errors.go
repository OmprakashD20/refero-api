package errors

import (
	"fmt"
	"net/http"
)

// HTTPError represents an HTTP error with a status code, a message, and an optional cause.
type HTTPError struct {
	ErrorMsg   string // Error message
	StatusCode int    // HTTP status code
	Cause      error  // Underlying error, if any
}

// Option is a functional option for configuring an HTTPError.
type Option func(*HTTPError)

// Error returns the formatted error message including the status code
// and message. If there is a cause, it includes that as well.
func (e *HTTPError) Error() string {
	if e.Cause != nil {
		return fmt.Sprintf("[%d %s]: %s | Cause: %v", e.StatusCode, http.StatusText(e.StatusCode), e.ErrorMsg, e.Cause)
	}
	return fmt.Sprintf("[%d %s]: %s", e.StatusCode, http.StatusText(e.StatusCode), e.ErrorMsg)
}

// WithCause returns an Option that sets the cause of an HTTPError.
func WithCause(err error) Option {
	return func(e *HTTPError) {
		e.Cause = err
	}
}

// WithError returns an Option that sets the error message for an HTTPError.
// The provided value can be a string or an error.
func WithError(err any) Option {
	return func(e *HTTPError) {
		switch v := err.(type) {
		case error:
			e.ErrorMsg = v.Error()
		case string:
			e.ErrorMsg = v
		default:
			e.ErrorMsg = "Unknown Error"
		}
	}
}

// WithStatus returns an Option that sets the HTTP status code for an HTTPError.
func WithStatus(statusCode int) Option {
	return func(e *HTTPError) {
		e.StatusCode = statusCode
	}
}

// NewHTTPError creates a new HTTPError with the given options.
// If no options are provided, it defaults to an internal server error.
func NewHTTPError(opts ...Option) *HTTPError {
	err := &HTTPError{
		ErrorMsg:   "Internal Server Error",
		StatusCode: http.StatusInternalServerError,
	}
	for _, opt := range opts {
		opt(err)
	}
	return err
}

// WithOptions returns an Option that applies multiple functional options to an HTTPError.
func WithOptions(opts ...Option) Option {
	return func(e *HTTPError) {
		for _, opt := range opts {
			opt(e)
		}
	}
}

func NotFound(err any, opts ...Option) *HTTPError {
	return NewHTTPError(WithStatus(http.StatusNotFound), WithError(err), WithOptions(opts...))
}

func Unauthorized(err any, opts ...Option) *HTTPError {
	return NewHTTPError(WithStatus(http.StatusUnauthorized), WithError(err), WithOptions(opts...))
}

func Forbidden(err any, opts ...Option) *HTTPError {
	return NewHTTPError(WithStatus(http.StatusForbidden), WithError(err), WithOptions(opts...))
}

func BadRequest(err any, opts ...Option) *HTTPError {
	return NewHTTPError(WithStatus(http.StatusBadRequest), WithError(err), WithOptions(opts...))
}

func Conflict(err any, opts ...Option) *HTTPError {
	return NewHTTPError(WithStatus(http.StatusConflict), WithError(err), WithOptions(opts...))
}

func Validation(err any, opts ...Option) *HTTPError {
	return NewHTTPError(WithStatus(http.StatusUnprocessableEntity), WithError(err), WithOptions(opts...))
}

func InternalServerError(opts ...Option) *HTTPError {
	return NewHTTPError(WithOptions(opts...))
}
