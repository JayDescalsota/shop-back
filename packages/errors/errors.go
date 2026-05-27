package errors

import (
	"fmt"
	"net/http"
)

type Code string

const (
	CodeInternal         Code = "INTERNAL_ERROR"
	CodeNotFound         Code = "NOT_FOUND"
	CodeValidation       Code = "VALIDATION_ERROR"
	CodeUnauthorized     Code = "UNAUTHORIZED"
	CodeForbidden        Code = "FORBIDDEN"
	CodeConflict         Code = "CONFLICT"
	CodeRateLimited      Code = "RATE_LIMITED"
	CodeTenantMismatch   Code = "TENANT_MISMATCH"
	CodePaymentFailed    Code = "PAYMENT_FAILED"
	CodeDuplicate        Code = "DUPLICATE_ENTRY"
	CodeBadRequest       Code = "BAD_REQUEST"
)

type AppError struct {
	Code    Code   `json:"code"`
	Message string `json:"message"`
	Err     error  `json:"-"`
}

func (e *AppError) Error() string {
	if e.Err != nil {
		return fmt.Sprintf("[%s] %s: %v", e.Code, e.Message, e.Err)
	}
	return fmt.Sprintf("[%s] %s", e.Code, e.Message)
}

func (e *AppError) Unwrap() error {
	return e.Err
}

func HTTPStatus(code Code) int {
	switch code {
	case CodeNotFound:
		return http.StatusNotFound
	case CodeValidation:
		return http.StatusUnprocessableEntity
	case CodeUnauthorized:
		return http.StatusUnauthorized
	case CodeForbidden, CodeTenantMismatch:
		return http.StatusForbidden
	case CodeConflict, CodeDuplicate:
		return http.StatusConflict
	case CodeRateLimited:
		return http.StatusTooManyRequests
	case CodeBadRequest:
		return http.StatusBadRequest
	default:
		return http.StatusInternalServerError
	}
}

func New(code Code, message string) *AppError {
	return &AppError{Code: code, Message: message}
}

func Wrap(code Code, message string, err error) *AppError {
	return &AppError{Code: code, Message: message, Err: err}
}

func NotFound(message string) *AppError {
	return New(CodeNotFound, message)
}

func Unauthorized(message string) *AppError {
	return New(CodeUnauthorized, message)
}

func Forbidden(message string) *AppError {
	return New(CodeForbidden, message)
}

func Validation(message string) *AppError {
	return New(CodeValidation, message)
}

func Internal(message string, err error) *AppError {
	return Wrap(CodeInternal, message, err)
}

func Conflict(message string) *AppError {
	return New(CodeConflict, message)
}

func TenantMismatch() *AppError {
	return New(CodeTenantMismatch, "resource does not belong to your tenant")
}
