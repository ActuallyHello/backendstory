package dto

import "time"

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
