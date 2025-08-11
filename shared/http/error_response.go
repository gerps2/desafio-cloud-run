package http

import (
	"net/http"

	"github.com/gerps2/desafio-cloud-run/shared/errors"
	"github.com/gin-gonic/gin"
)

type APIResponse struct {
	Data    interface{} `json:"data"`
	Message string      `json:"message"`
	Causes  []string    `json:"causes,omitempty"`
}

func RespondWithSuccess(c *gin.Context, data interface{}, message string) {
	response := APIResponse{
		Data:    data,
		Message: message,
		Causes:  nil,
	}
	c.JSON(http.StatusOK, response)
}

func RespondWithAPIError(c *gin.Context, apiError *errors.APIError) {
	response := APIResponse{
		Data:    nil,
		Message: apiError.Message,
		Causes:  apiError.Causes,
	}
	c.JSON(apiError.StatusCode, response)
}

func RespondWithError(c *gin.Context, statusCode int, message string, causes []string) {
	response := APIResponse{
		Data:    nil,
		Message: message,
		Causes:  causes,
	}
	c.JSON(statusCode, response)
}

func RespondWithValidationError(c *gin.Context, message string, causes []string) {
	apiError := errors.NewValidationError(message, causes)
	RespondWithAPIError(c, apiError)
}

func RespondWithBusinessError(c *gin.Context, code, message string, causes []string) {
	apiError := errors.NewBusinessError(code, message, causes)
	RespondWithAPIError(c, apiError)
}

func RespondWithNotFound(c *gin.Context, message string, causes []string) {
	apiError := errors.NewNotFoundError(message, causes)
	RespondWithAPIError(c, apiError)
}

func RespondWithInternalError(c *gin.Context, message string, causes []string) {
	if message == "" {
		message = "Internal server error occurred"
	}
	apiError := errors.NewInternalError(message, causes)
	RespondWithAPIError(c, apiError)
}

func RespondWithExternalServiceError(c *gin.Context, message string, causes []string) {
	apiError := errors.NewExternalServiceError(message, causes)
	RespondWithAPIError(c, apiError)
}

func RespondWithTimeout(c *gin.Context, message string, causes []string) {
	if message == "" {
		message = "Request timeout exceeded"
	}
	apiError := errors.NewTimeoutError(message, causes)
	RespondWithAPIError(c, apiError)
}
