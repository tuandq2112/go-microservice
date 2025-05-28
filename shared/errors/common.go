package errors

import (
	"encoding/json"
	"net/http"
	"strings"
)

// ErrorResponse represents a standardized error response
type ErrorResponse struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Error   string `json:"error,omitempty"`
}

// Common error message IDs for localization
const (
	ErrUnauthorized     = "error.unauthorized"
	ErrForbidden        = "error.forbidden"
	ErrNotFound         = "error.not_found"
	ErrInternalServer   = "error.internal_server"
	ErrBadRequest       = "error.bad_request"
	ErrValidationFailed = "error.validation_failed"
	ErrConnectionFailed = "error.connection_failed"
)

// Common HTTP errors
var (
	UnauthorizedError = ErrorResponse{
		Code:    http.StatusUnauthorized,
		Message: ErrUnauthorized,
	}

	ForbiddenError = ErrorResponse{
		Code:    http.StatusForbidden,
		Message: ErrForbidden,
	}

	NotFoundError = ErrorResponse{
		Code:    http.StatusNotFound,
		Message: ErrNotFound,
	}

	InternalServerError = ErrorResponse{
		Code:    http.StatusInternalServerError,
		Message: ErrInternalServer,
	}

	BadRequestError = ErrorResponse{
		Code:    http.StatusBadRequest,
		Message: ErrBadRequest,
	}

	ValidationError = ErrorResponse{
		Code:    http.StatusBadRequest,
		Message: ErrValidationFailed,
	}

	ConnectionError = ErrorResponse{
		Code:    http.StatusServiceUnavailable,
		Message: ErrConnectionFailed,
	}
)

// WriteError writes an error response to the HTTP response writer
func WriteError(w http.ResponseWriter, err ErrorResponse, localizedMessage string, originalError error) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(err.Code)

	response := ErrorResponse{
		Code:    err.Code,
		Message: localizedMessage,
	}

	if originalError != nil {
		response.Error = originalError.Error()
	}

	json.NewEncoder(w).Encode(response)
}

// NewError creates a new error response with a custom message
func NewError(code int, messageID string) ErrorResponse {
	return ErrorResponse{
		Code:    code,
		Message: messageID,
	}
}

// IsConnectionError checks if the error is a connection error
func IsConnectionError(err error) bool {
	if err == nil {
		return false
	}
	errStr := err.Error()
	return strings.Contains(errStr, "connection error") ||
		strings.Contains(errStr, "connection refused") ||
		strings.Contains(errStr, "dial tcp") ||
		strings.Contains(errStr, "transport: Error while dialing")
}
