package errors

import (
	"net/http"
)

type APIError struct {
	Code       string   `json:"code"`
	Message    string   `json:"message"`
	StatusCode int      `json:"-"`
	Causes     []string `json:"causes,omitempty"`
	Context    string   `json:"context,omitempty"`
}

func (e APIError) Error() string {
	return e.Message
}

type ErrorType string

const (
	ValidationError ErrorType = "validation"
	BusinessError   ErrorType = "business"
	SystemError     ErrorType = "system"
	ExternalError   ErrorType = "external"
)

const (
	// Erros de validação (400)
	CodeInvalidInput     = "INVALID_INPUT"
	CodeMissingParameter = "MISSING_PARAMETER"
	CodeInvalidFormat    = "INVALID_FORMAT"

	// Erros de negócio (400-404)
	CodeResourceNotFound = "RESOURCE_NOT_FOUND"
	CodeBusinessRule     = "BUSINESS_RULE_VIOLATION"

	// Erros de sistema (500)
	CodeInternalError = "INTERNAL_SERVER_ERROR"
	CodeDatabaseError = "DATABASE_ERROR"

	// Erros externos (502-504)
	CodeExternalService    = "EXTERNAL_SERVICE_ERROR"
	CodeServiceTimeout     = "SERVICE_TIMEOUT"
	CodeServiceUnavailable = "SERVICE_UNAVAILABLE"

	// Erros de autenticação/autorização (401-403)
	CodeUnauthorized = "UNAUTHORIZED"
	CodeForbidden    = "FORBIDDEN"
)

func NewAPIError(code, message string, statusCode int, causes []string) *APIError {
	return &APIError{
		Code:       code,
		Message:    message,
		StatusCode: statusCode,
		Causes:     causes,
	}
}

func NewValidationError(message string, causes []string) *APIError {
	return &APIError{
		Code:       CodeInvalidInput,
		Message:    message,
		StatusCode: http.StatusBadRequest,
		Causes:     causes,
		Context:    string(ValidationError),
	}
}

func NewBusinessError(code, message string, causes []string) *APIError {
	return &APIError{
		Code:       code,
		Message:    message,
		StatusCode: http.StatusBadRequest,
		Causes:     causes,
		Context:    string(BusinessError),
	}
}

func NewNotFoundError(message string, causes []string) *APIError {
	return &APIError{
		Code:       CodeResourceNotFound,
		Message:    message,
		StatusCode: http.StatusNotFound,
		Causes:     causes,
		Context:    string(BusinessError),
	}
}

func NewInternalError(message string, causes []string) *APIError {
	return &APIError{
		Code:       CodeInternalError,
		Message:    message,
		StatusCode: http.StatusInternalServerError,
		Causes:     causes,
		Context:    string(SystemError),
	}
}

func NewExternalServiceError(message string, causes []string) *APIError {
	return &APIError{
		Code:       CodeExternalService,
		Message:    message,
		StatusCode: http.StatusBadGateway,
		Causes:     causes,
		Context:    string(ExternalError),
	}
}

func NewTimeoutError(message string, causes []string) *APIError {
	return &APIError{
		Code:       CodeServiceTimeout,
		Message:    message,
		StatusCode: http.StatusGatewayTimeout,
		Causes:     causes,
		Context:    string(ExternalError),
	}
}
