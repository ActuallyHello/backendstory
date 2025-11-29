package core

import (
	"fmt"
	"time"
)

type NotFoundError struct {
	Message string
}

func (e *NotFoundError) Error() string {
	return e.Message
}

func NewNotFoundError(msg string) *NotFoundError {
	return &NotFoundError{Message: msg}
}

type ErrorInfo struct {
	Code    string
	Message string
	Err     error
}

type LogicalError struct {
	ErrorInfo
}

func (e *LogicalError) Error() string {
	return errorString(e.Err, e.Code, e.Message)
}

func NewLogicalError(err error, code, message string) *LogicalError {
	return &LogicalError{
		ErrorInfo: ErrorInfo{
			Code:    code,
			Message: message,
			Err:     err,
		},
	}
}

type TechnicalError struct {
	ErrorInfo
}

func (e *TechnicalError) Error() string {
	return errorString(e.Err, e.Code, e.Message)
}

func NewTechnicalError(err error, code, message string) *TechnicalError {
	return &TechnicalError{
		ErrorInfo: ErrorInfo{
			Code:    code,
			Message: message,
			Err:     err,
		},
	}
}

type ValidationError struct {
	ErrorInfo
}

func (e *ValidationError) Error() string {
	return errorString(e.Err, e.Code, e.Message)
}

func NewValidationError(err error, code, message string) *ValidationError {
	return &ValidationError{
		ErrorInfo: ErrorInfo{
			Code:    code,
			Message: message,
			Err:     err,
		},
	}
}

type AccessError struct {
	ErrorInfo
}

func (e *AccessError) Error() string {
	return errorString(e.Err, e.Code, e.Message)
}

func NewAccessError(err error, code, message string) *AccessError {
	return &AccessError{
		ErrorInfo: ErrorInfo{
			Code:    code,
			Message: message,
			Err:     err,
		},
	}
}

func errorString(err error, code, message string) string {
	msg := fmt.Sprintf("[%s] %s", code, message)
	// TODO: refactor len(message)!=len(err) - chose error message by application start mod (local, dev, prod)
	if err != nil && len(message) != len(err.Error()) {
		return fmt.Sprintf("%s. ERROR: %v", msg, err)
	}
	return msg
}

// ErrorResponse represents standard error response
// @Name ErrorResponse
type ErrorResponse struct {
	Status    int       `json:"status"`
	Code      string    `json:"code"`
	Message   string    `json:"message"`
	Path      string    `json:"path"`
	Method    string    `json:"method"`
	Timestamp time.Time `json:"timestamp"`
}

// ValidationErrorResponse represents validation error response
// @Name ValidationErrorResponse
type ValidationErrorResponse struct {
	ErrorResponse
	Details map[string]string `json:"details"`
}

func NewErrorResponse(status int, code, message, path, method string) *ErrorResponse {
	return &ErrorResponse{
		Status:    status,
		Code:      code,
		Message:   message,
		Path:      path,
		Method:    method,
		Timestamp: time.Now(),
	}
}

func NewValidationErrorResponse(status int, message, path string, details map[string]string) *ValidationErrorResponse {
	return &ValidationErrorResponse{
		ErrorResponse: ErrorResponse{
			Status:    status,
			Message:   message,
			Path:      path,
			Timestamp: time.Now(),
		},
		Details: details,
	}
}
