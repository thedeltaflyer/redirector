package helpers

import "fmt"

// AlreadyExistsError represents an error indicating that a requested key or entry already exists in the datastore.
type AlreadyExistsError struct {
	err error
}

// Error returns the error message. If the receiver or wrapped error is nil, it returns "<nil>".
func (e *AlreadyExistsError) Error() string {
	if e == nil || e.err == nil {
		return "<nil>"
	}
	return e.err.Error()
}

// Unwrap returns the wrapped error if it exists; otherwise, it returns nil.
func (e *AlreadyExistsError) Unwrap() error {
	if e == nil {
		return nil
	}
	return e.err
}

// NewAlreadyExistsError creates a new AlreadyExistsError with a message indicating the given key already exists.
func NewAlreadyExistsError(key []byte) *AlreadyExistsError {
	return &AlreadyExistsError{
		err: fmt.Errorf("key %q already exists", key),
	}
}

// DoesNotExistError represents an error indicating that a requested key or item does not exist in the storage.
type DoesNotExistError struct {
	err error
}

// Error returns the error message. If the receiver or wrapped error is nil, it returns "<nil>".
func (e *DoesNotExistError) Error() string {
	if e == nil || e.err == nil {
		return "<nil>"
	}
	return e.err.Error()
}

// Unwrap returns the wrapped error if it exists; otherwise, it returns nil.
func (e *DoesNotExistError) Unwrap() error {
	if e == nil {
		return nil
	}
	return e.err
}

// NewDoesNotExistError creates a new DoesNotExistError with a message indicating the given key does not exist.
func NewDoesNotExistError(key []byte) *DoesNotExistError {
	return &DoesNotExistError{
		err: fmt.Errorf("key %q does not exist", key),
	}
}
