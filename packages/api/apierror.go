package api

import "fmt"

// APIError captures an HTTP-oriented error to keep handler logic framework-agnostic.
type APIError struct {
	Status  int
	Message string
}

func (e *APIError) Error() string {
	if e == nil {
		return ""
	}
	return fmt.Sprintf("[%d] %s", e.Status, e.Message)
}

// newAPIError creates a new APIError helper for internal package use.
func newAPIError(status int, message string) *APIError {
	return &APIError{
		Status:  status,
		Message: message,
	}
}

// NewAPIError exposes APIError creation for external consumers.
func NewAPIError(status int, message string) *APIError {
	return &APIError{
		Status:  status,
		Message: message,
	}
}
