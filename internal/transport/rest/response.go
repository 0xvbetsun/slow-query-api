package rest

import (
	"fmt"
	"net/http"
	"strings"
)

// Pagination used for splitting large data into the small pieces
type Pagination struct {
	Page     int `json:"page"`
	PageSize int `json:"page_size"`
}

// ResponseWithPayload common response with any data in Payload field.
type ResponseWithPayload struct {
	Message    string      `json:"message,omitempty"`
	Pagination Pagination  `json:"pagination,omitempty"`
	Payload    interface{} `json:"payload,omitempty"`
}

// ErrorResponse defines the body of an error response.
type ErrorResponse struct {
	Message string `json:"message"`
	Field   string `json:"field,omitempty"`
}

// ErrWithHint custom error for ErrorHandler.
type ErrWithHint struct {
	Code    int
	Message string
	Field   string
	Err     error
}

func (e ErrWithHint) Error() string {
	return fmt.Sprintf("%d: %s", e.Code, e.Message)
}

func (e *ErrWithHint) WithErr(err error) *ErrWithHint {
	e.Err = err

	return e
}
func ErrBadRequestInvalidParameter(name string, description ...string) *ErrWithHint {
	return &ErrWithHint{
		Code:    http.StatusBadRequest,
		Message: fmt.Sprintf("Invalid parameter '%s'. %s", name, strings.Join(description, ".")),
		Field:   name,
	}
}
