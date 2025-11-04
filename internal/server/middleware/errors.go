package middleware

import (
	"encoding/json"
	"errors"
	"log/slog"
	"net/http"

	appError "github.com/ActuallyHello/backendstory/internal/core/errors"
	"github.com/ActuallyHello/backendstory/internal/dto"
)

func ErrorHandler(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				slog.Error("Panic recovered", "error", err, "URL", r.URL.Path)
				HandleError(w, r, &appError.TechnicalError{
					ErrorInfo: appError.ErrorInfo{
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
	var logicErr *appError.LogicalError
	var techErr *appError.TechnicalError
	var accessErr *appError.AccessError

	var response dto.ErrorResponse
	switch {
	case errors.As(err, &logicErr):
		response = *dto.NewErrorResponse(
			http.StatusInternalServerError,
			logicErr.Code,
			logicErr.Error(),
			r.URL.Path,
			r.Method,
		)
	case errors.As(err, &techErr):
		response = *dto.NewErrorResponse(
			http.StatusInternalServerError,
			techErr.Code,
			techErr.Error(),
			r.URL.Path,
			r.Method,
		)
	case errors.As(err, &accessErr):
		response = *dto.NewErrorResponse(
			http.StatusUnauthorized,
			accessErr.Code,
			accessErr.Error(),
			r.URL.Path,
			r.Method,
		)
	default:
		response = *dto.NewErrorResponse(
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

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(response.Status)

	json.NewEncoder(w).Encode(response)
}

func HandleValidationError(w http.ResponseWriter, r *http.Request, err error, details map[string]string) {
	response := *dto.NewValidationErrorResponse(
		http.StatusBadRequest,
		err.Error(),
		r.URL.Path,
		details,
	)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusBadRequest)

	json.NewEncoder(w).Encode(response)
}
