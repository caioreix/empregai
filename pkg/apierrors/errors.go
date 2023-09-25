package apierrors

import (
	"database/sql"
	"errors"
	"fmt"
	"net/http"
	"strings"

	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
)

const separator = ": "

// APIError represents a common API error structure.
type APIError struct {
	HTTPStatus int    `json:"-"`
	ErrCode    string `json:"code"`
	Message    string `json:"message"`
}

// NewAPIError creates a new instance of APIError.
func NewAPIError(statusCode int, code, message string) *APIError {
	if code == "" {
		code = fmt.Sprint(statusCode)
	}
	if message == "" {
		message = http.StatusText(statusCode)
	}

	return &APIError{
		HTTPStatus: statusCode,
		ErrCode:    code,
		Message:    message,
	}
}

func (e *APIError) Code(code string) {
	e.ErrCode = code
}

// BadRequest creates a 400 Bad Request error.
func BadRequest(messages ...string) *APIError {
	return NewAPIError(http.StatusBadRequest, "", strings.Join(messages, separator))
}

// Unauthorized creates a 401 Unauthorized error.
func Unauthorized(messages ...string) *APIError {
	return NewAPIError(http.StatusUnauthorized, "", strings.Join(messages, separator))
}

// Forbidden creates a 403 Forbidden error.
func Forbidden(messages ...string) *APIError {
	return NewAPIError(http.StatusForbidden, "", strings.Join(messages, separator))
}

// NotFound creates a 404 Not Found error.
func NotFound(messages ...string) *APIError {
	return NewAPIError(http.StatusNotFound, "", strings.Join(messages, separator))
}

func Conflict(messages ...string) *APIError {
	return NewAPIError(http.StatusConflict, "", strings.Join(messages, separator))
}

// InternalServerError creates a 500 Internal Server Error.
func InternalServerError(messages ...string) *APIError {
	return NewAPIError(http.StatusInternalServerError, "", strings.Join(messages, separator))
}

// NotImplemented creates a 501 Not Implemented error.
func NotImplemented(messages ...string) *APIError {
	return NewAPIError(http.StatusNotImplemented, "", strings.Join(messages, separator))
}

func Parse(err error) *APIError {
	var apiErr *APIError

	switch {
	case errors.As(err, &apiErr): // First check!
		return apiErr
	case errors.Is(err, http.ErrNoCookie):
		return Unauthorized()
	case errors.Is(err, sql.ErrNoRows):
		return NotFound()
	case uuid.IsInvalidLengthError(err):
		return BadRequest()
	case errors.Is(err, redis.Nil):
		return NotFound()
	case strings.Contains(err.Error(), "duplicate key value"):
		return Conflict()
	default:
		return InternalServerError()
	}
}

// JSON represents the error in JSON format.
func (e *APIError) JSON() (int, map[string]interface{}) {
	return e.HTTPStatus, map[string]interface{}{
		"code":    e.ErrCode,
		"message": e.Message,
	}
}

// JSON represents the error in JSON format.
func JSON(err error) (int, map[string]interface{}) {
	var apiErr *APIError

	if errors.As(err, &apiErr) {
		return apiErr.JSON()
	}

	return InternalServerError().JSON()
}

// StatusCode returns the associated HTTP status code for the error.
func (e *APIError) StatusCode() int {
	return e.HTTPStatus
}

// Error implements the error interface.
func (e *APIError) Error() string {
	return fmt.Sprintf("API Error: %s - %s", e.ErrCode, e.Message)
}
