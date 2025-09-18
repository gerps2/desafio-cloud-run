package getWeatherByCep

import (
	"net/http"

	sharedErrors "github.com/gerps2/desafio-cloud-run/shared/errors"
)

const (
	CodeInvalidZipcode      = "INVALID_ZIPCODE"
	CodeZipcodeNotFound     = "ZIPCODE_NOT_FOUND"
	CodeWeatherServiceError = "WEATHER_SERVICE_ERROR"
)

func NewInvalidZipcodeError() *sharedErrors.APIError {
	return sharedErrors.NewAPIError(
		CodeInvalidZipcode,
		"invalid zipcode",
		http.StatusUnprocessableEntity,
		[]string{"The provided zipcode format is invalid"},
	)
}

func NewZipcodeNotFoundError() *sharedErrors.APIError {
	return sharedErrors.NewAPIError(
		CodeZipcodeNotFound,
		"can not find zipcode",
		http.StatusNotFound,
		[]string{"The provided zipcode was not found"},
	)
}

func NewWeatherServiceError() *sharedErrors.APIError {
	return sharedErrors.NewExternalServiceError(
		"Weather service temporarily unavailable",
		[]string{"Unable to fetch weather data from external service"},
	)
}

func NewWeatherValidationError(message string, causes []string) *sharedErrors.APIError {
	return sharedErrors.NewValidationError(message, causes)
}

func NewWeatherBusinessError(code, message string, causes []string) *sharedErrors.APIError {
	return sharedErrors.NewBusinessError(code, message, causes)
}

func NewWeatherExternalError(message string, causes []string) *sharedErrors.APIError {
	return sharedErrors.NewExternalServiceError(message, causes)
}
