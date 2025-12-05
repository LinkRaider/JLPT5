package errors

import (
	"fmt"
	"net/http"
)

// ErrorCode represents a specific error type
type ErrorCode string

const (
	// General errors
	ErrCodeInternal       ErrorCode = "INTERNAL_ERROR"
	ErrCodeBadRequest     ErrorCode = "BAD_REQUEST"
	ErrCodeNotFound       ErrorCode = "NOT_FOUND"
	ErrCodeUnauthorized   ErrorCode = "UNAUTHORIZED"
	ErrCodeForbidden      ErrorCode = "FORBIDDEN"
	ErrCodeConflict       ErrorCode = "CONFLICT"
	ErrCodeValidation     ErrorCode = "VALIDATION_ERROR"

	// Database errors
	ErrCodeDatabase       ErrorCode = "DATABASE_ERROR"
	ErrCodeDuplicate      ErrorCode = "DUPLICATE_ENTRY"

	// Authentication errors
	ErrCodeInvalidCredentials ErrorCode = "INVALID_CREDENTIALS"
	ErrCodeTokenExpired       ErrorCode = "TOKEN_EXPIRED"
	ErrCodeTokenInvalid       ErrorCode = "TOKEN_INVALID"
)

// AppError represents an application error with code and HTTP status
type AppError struct {
	Code       ErrorCode `json:"code"`
	Message    string    `json:"message"`
	StatusCode int       `json:"-"`
	Err        error     `json:"-"`
}

// Error implements the error interface
func (e *AppError) Error() string {
	if e.Err != nil {
		return fmt.Sprintf("%s: %s (%v)", e.Code, e.Message, e.Err)
	}
	return fmt.Sprintf("%s: %s", e.Code, e.Message)
}

// Unwrap returns the wrapped error
func (e *AppError) Unwrap() error {
	return e.Err
}

// NewAppError creates a new application error
func NewAppError(code ErrorCode, message string, statusCode int, err error) *AppError {
	return &AppError{
		Code:       code,
		Message:    message,
		StatusCode: statusCode,
		Err:        err,
	}
}

// Common error constructors

// Internal creates an internal server error
func Internal(message string, err error) *AppError {
	return NewAppError(ErrCodeInternal, message, http.StatusInternalServerError, err)
}

// BadRequest creates a bad request error
func BadRequest(message string) *AppError {
	return NewAppError(ErrCodeBadRequest, message, http.StatusBadRequest, nil)
}

// NotFound creates a not found error
func NotFound(message string) *AppError {
	return NewAppError(ErrCodeNotFound, message, http.StatusNotFound, nil)
}

// Unauthorized creates an unauthorized error
func Unauthorized(message string) *AppError {
	return NewAppError(ErrCodeUnauthorized, message, http.StatusUnauthorized, nil)
}

// Forbidden creates a forbidden error
func Forbidden(message string) *AppError {
	return NewAppError(ErrCodeForbidden, message, http.StatusForbidden, nil)
}

// Conflict creates a conflict error
func Conflict(message string) *AppError {
	return NewAppError(ErrCodeConflict, message, http.StatusConflict, nil)
}

// Validation creates a validation error
func Validation(message string) *AppError {
	return NewAppError(ErrCodeValidation, message, http.StatusBadRequest, nil)
}

// InvalidCredentials creates an invalid credentials error
func InvalidCredentials() *AppError {
	return NewAppError(ErrCodeInvalidCredentials, "Invalid email or password", http.StatusUnauthorized, nil)
}

// TokenExpired creates a token expired error
func TokenExpired() *AppError {
	return NewAppError(ErrCodeTokenExpired, "Token has expired", http.StatusUnauthorized, nil)
}

// TokenInvalid creates an invalid token error
func TokenInvalid() *AppError {
	return NewAppError(ErrCodeTokenInvalid, "Invalid token", http.StatusUnauthorized, nil)
}
