package core

import (
	"fmt"
	"runtime"
	"time"
)

type StackTracer interface {
	error
	StackTrace() []uintptr
}

type StackError struct {
	msg   string
	cause error
	stack []uintptr
}

func NewStackError(msg string) *StackError {
	pcs := make([]uintptr, 32)
	n := runtime.Callers(2, pcs)

	return &StackError{
		msg:   msg,
		stack: pcs[:n],
	}
}

func (e *StackError) Error() string {
	if e.cause != nil {
		// TODO: refactor errors
		// logical/techinal errors duplicate log info
		return e.cause.Error()
	}
	return e.msg
}

func (e *StackError) Unwrap() error {
	return e.cause
}

func (e *StackError) StackTrace() []uintptr {
	return e.stack
}

func WrapStack(err error, msg string) error {
	if err == nil {
		return nil
	}

	if se, ok := err.(*StackError); ok {
		return &StackError{
			msg:   msg,
			cause: se,
			stack: se.stack,
		}
	}

	pcs := make([]uintptr, 32)
	n := runtime.Callers(2, pcs)
	return &StackError{
		msg:   msg,
		cause: err,
		stack: pcs[:n],
	}
}

type NotFoundError struct {
	Message string
}

func (e *NotFoundError) Error() string {
	return e.Message
}

func (e *NotFoundError) Is(target error) bool {
	_, ok := target.(*NotFoundError)
	return ok
}

func NewNotFoundError(msg string) *NotFoundError {
	return &NotFoundError{Message: msg}
}

type ErrorInfo struct {
	Code    string
	Message string
	Err     error
}

func (e *ErrorInfo) Unwrap() error {
	return e.Err
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
			Err:     WrapStack(err, message),
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
			Err:     WrapStack(err, message),
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
	msg := fmt.Sprintf("[%s] %s", code, err.Error())
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
