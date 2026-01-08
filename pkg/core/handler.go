package core

import (
	"encoding/json"
	"errors"
	"log/slog"
	"net/http"
	"runtime/debug"

	"github.com/go-playground/validator/v10"
)

func CollectValidationDetails(err error) map[string]string {
	validationErrors := err.(validator.ValidationErrors)

	details := make(map[string]string)
	if len(validationErrors) > 0 {
		for _, validationError := range validationErrors {
			details[validationError.Field()] = validationError.Tag()
		}
	}

	return details
}

func ErrorHandler(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				slog.Error("Panic recovered", "error", err, "URL", r.URL.Path)
				HandleError(w, r, &TechnicalError{
					ErrorInfo: ErrorInfo{
						Code:    "INTERNAL_ERROR",
						Message: "Internal server error",
					},
				})
			}
		}()

		rw := &responseWriter{ResponseWriter: w, status: http.StatusOK}
		next.ServeHTTP(rw, r)
	})
}

type responseWriter struct {
	http.ResponseWriter
	status int
}

func (rw *responseWriter) WriteHeader(code int) {
	rw.status = code
	rw.ResponseWriter.WriteHeader(code)
}

func HandleError(w http.ResponseWriter, r *http.Request, err error) {
	var logicErr *LogicalError
	var techErr *TechnicalError
	var accessErr *AccessError

	var response ErrorResponse
	switch {
	case errors.As(err, &logicErr):
		response = *NewErrorResponse(
			http.StatusInternalServerError,
			logicErr.Code,
			logicErr.Error(),
			r.URL.Path,
			r.Method,
		)
	case errors.As(err, &techErr):
		response = *NewErrorResponse(
			http.StatusInternalServerError,
			techErr.Code,
			techErr.Error(),
			r.URL.Path,
			r.Method,
		)
	case errors.As(err, &accessErr):
		response = *NewErrorResponse(
			http.StatusUnauthorized,
			accessErr.Code,
			accessErr.Error(),
			r.URL.Path,
			r.Method,
		)
	default:
		response = *NewErrorResponse(
			http.StatusInternalServerError,
			"ERROR",
			err.Error(),
			r.URL.Path,
			r.Method,
		)
	}

	slog.Error(
		"Handle error",
		"Code", response.Code,
		"Message", response.Message,
	)
	debug.PrintStack()

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(response.Status)

	json.NewEncoder(w).Encode(response)
}

func HandleValidationError(w http.ResponseWriter, r *http.Request, err error, details map[string]string) {
	response := *NewValidationErrorResponse(
		http.StatusBadRequest,
		err.Error(),
		r.URL.Path,
		details,
	)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusBadRequest)

	json.NewEncoder(w).Encode(response)
}
