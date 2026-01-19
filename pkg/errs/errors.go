package errs

import (
	"errors"
	"fmt"
)

// NotFoundError is a wrapper for an error when something is not found
type NotFoundError struct {
	message string
}

func (err *NotFoundError) Error() string {
	return err.message
}

// NotFound returns a new NotFoundError
func NotFound(message string) *NotFoundError {
	return &NotFoundError{message: message}
}

func IsNotFound(err error) bool {
	if err == nil {
		return false
	}

	var e *NotFoundError
	return errors.As(err, &e)
}

// InvalidArgumentError is a wrapper for an error when an user provides an invalid input
type InvalidArgumentError struct {
	message string
}

func (err *InvalidArgumentError) Error() string {
	return err.message
}

// InvalidArgument returns a new InvalidArgumentError
func InvalidArgument(message string) *InvalidArgumentError {
	return &InvalidArgumentError{
		message: message,
	}
}

// InternalError is a wrapper for an error when something bad happened, but we don't want to inform the user
// of the details
type InternalError struct {
	err error
}

func (err *InternalError) Error() string {
	return fmt.Sprintf("Internal error: %s", err.err)
}

func (err *InternalError) Message() string {
	return err.err.Error()
}

// Internal returns a new InternalError
func Internal(message string, err error) *InternalError {
	return &InternalError{
		err: fmt.Errorf("%s: %w", message, err),
	}
}

func IsInternal(err error) bool {
	if err == nil {
		return false
	}

	switch err.(type) {
	case *NotFoundError, *InvalidArgumentError,
		*PermissionDeniedError, *AuthenticationRequiredError, *AlreadyExistsError, *TooManyRequestsError:
		return false
	default:
		return true
	}
}

// PermissionDeniedError is a wrapper for an error when an authenticated user try to perform an
// action he does not have the rights for
type PermissionDeniedError struct {
	message string
}

func (err *PermissionDeniedError) Error() string {
	return err.message
}

// PermissionDenied returns a new PermissionDeniedError
func PermissionDenied(message string) *PermissionDeniedError {
	return &PermissionDeniedError{message: message}
}

func IsPermissionDenied(err error) bool {
	if err == nil {
		return false
	}

	var e *PermissionDeniedError
	return errors.As(err, &e)
}

// AlreadyExistsError is a wrapper for an error when something already exists
type AlreadyExistsError struct {
	message string
}

func (err *AlreadyExistsError) Error() string {
	return err.message
}

// AlreadyExists returns a new AlreadyExistsError
func AlreadyExists(message string) *AlreadyExistsError {
	return &AlreadyExistsError{
		message: message,
	}
}

// AuthenticationRequiredError is a wrapper for an error when something an unauthenticated person try
// to perform an action which requires authentication
type AuthenticationRequiredError struct {
	message string
}

func (err *AuthenticationRequiredError) Error() string {
	if err.message == "" {
		return "Authentication required."
	}
	return err.message
}

// AuthenticationRequired returns a new AuthenticationRequiredError
func AuthenticationRequired(message string) *AuthenticationRequiredError {
	return &AuthenticationRequiredError{message: message}
}

// TooManyRequestsError is a wrapper for an error when a client is sending too many requests
type TooManyRequestsError struct {
	// message string
}

func (err *TooManyRequestsError) Error() string {
	return "Too many requests. Please try again later."
}

// TooManyRequests returns a new TooManyRequestsError
func TooManyRequests() *TooManyRequestsError {
	return &TooManyRequestsError{}
}
