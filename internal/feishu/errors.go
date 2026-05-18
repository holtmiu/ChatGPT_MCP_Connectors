package feishu

import "fmt"

type Code string

const (
	ErrAuthRequired            Code = "AUTH_REQUIRED"
	ErrAuthExpired             Code = "AUTH_EXPIRED"
	ErrPermissionDenied        Code = "PERMISSION_DENIED"
	ErrDocumentNotFound        Code = "DOCUMENT_NOT_FOUND"
	ErrUnsupportedDocumentType Code = "UNSUPPORTED_DOCUMENT_TYPE"
	ErrRateLimited             Code = "RATE_LIMITED"
	ErrPartialContent          Code = "PARTIAL_CONTENT"
	ErrWriteConflict           Code = "WRITE_CONFLICT"
	ErrInvalidInput            Code = "INVALID_INPUT"
	ErrUpstream                Code = "UPSTREAM_ERROR"
)

type ConnectorError struct {
	Code    Code   `json:"code"`
	Message string `json:"message"`
	Cause   error  `json:"-"`
}

func (e *ConnectorError) Error() string {
	if e == nil {
		return ""
	}
	if e.Cause != nil {
		return fmt.Sprintf("%s: %s: %v", e.Code, e.Message, e.Cause)
	}
	return fmt.Sprintf("%s: %s", e.Code, e.Message)
}

func (e *ConnectorError) Unwrap() error {
	if e == nil {
		return nil
	}
	return e.Cause
}

func newError(code Code, message string, cause error) *ConnectorError {
	return &ConnectorError{Code: code, Message: message, Cause: cause}
}
