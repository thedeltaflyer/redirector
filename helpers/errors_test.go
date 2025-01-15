package helpers

import (
	"errors"
	"testing"
)

func TestAlreadyExistsError(t *testing.T) {
	t.Run("Error", func(t *testing.T) {
		tests := []struct {
			name     string
			errorObj *AlreadyExistsError
			want     string
		}{
			{"NilErrorObject", nil, "<nil>"},
			{"NilWrappedError", &AlreadyExistsError{}, "<nil>"},
			{"NonNilWrappedError", &AlreadyExistsError{err: errors.New("test error")}, "test error"},
		}

		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				got := tt.errorObj.Error()
				if got != tt.want {
					t.Errorf("Error() = %v, want %v", got, tt.want)
				}
			})
		}
	})

	t.Run("Unwrap", func(t *testing.T) {
		tests := []struct {
			name     string
			errorObj *AlreadyExistsError
			want     error
		}{
			{"NilErrorObject", nil, nil},
			{"NilWrappedError", &AlreadyExistsError{}, nil},
			{"NonNilWrappedError", &AlreadyExistsError{err: errors.New("test error")}, errors.New("test error")},
		}

		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				got := tt.errorObj.Unwrap()
				if got == nil && tt.want == nil {
					return
				}
				if got == nil || tt.want == nil || got.Error() != tt.want.Error() {
					t.Errorf("Unwrap() = %v, want %v", got, tt.want)
				}
			})
		}
	})
}

func TestDoesNotExistError(t *testing.T) {
	t.Run("Error", func(t *testing.T) {
		tests := []struct {
			name     string
			errorObj *DoesNotExistError
			want     string
		}{
			{"NilErrorObject", nil, "<nil>"},
			{"NilWrappedError", &DoesNotExistError{}, "<nil>"},
			{"NonNilWrappedError", &DoesNotExistError{err: errors.New("test error")}, "test error"},
		}

		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				got := tt.errorObj.Error()
				if got != tt.want {
					t.Errorf("Error() = %v, want %v", got, tt.want)
				}
			})
		}
	})

	t.Run("Unwrap", func(t *testing.T) {
		tests := []struct {
			name     string
			errorObj *DoesNotExistError
			want     error
		}{
			{"NilErrorObject", nil, nil},
			{"NilWrappedError", &DoesNotExistError{}, nil},
			{"NonNilWrappedError", &DoesNotExistError{err: errors.New("test error")}, errors.New("test error")},
		}

		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				got := tt.errorObj.Unwrap()
				if got == nil && tt.want == nil {
					return
				}
				if got == nil || tt.want == nil || got.Error() != tt.want.Error() {
					t.Errorf("Unwrap() = %v, want %v", got, tt.want)
				}
			})
		}
	})
}

func TestNewAlreadyExistsError(t *testing.T) {
	tests := []struct {
		name     string
		key      []byte
		expected string
	}{
		{"EmptyKey", []byte(""), `key "" already exists`},
		{"NonEmptyKey", []byte("test"), `key "test" already exists`},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := NewAlreadyExistsError(tt.key)
			if err == nil {
				t.Fatalf("Expected non-nil error")
			}
			if err.Error() != tt.expected {
				t.Errorf("NewAlreadyExistsError() = %v, want %v", err.Error(), tt.expected)
			}
		})
	}
}

func TestNewDoesNotExistError(t *testing.T) {
	tests := []struct {
		name     string
		key      []byte
		expected string
	}{
		{"EmptyKey", []byte(""), `key "" does not exist`},
		{"NonEmptyKey", []byte("test"), `key "test" does not exist`},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := NewDoesNotExistError(tt.key)
			if err == nil {
				t.Fatalf("Expected non-nil error")
			}
			if err.Error() != tt.expected {
				t.Errorf("NewDoesNotExistError() = %v, want %v", err.Error(), tt.expected)
			}
		})
	}
}
