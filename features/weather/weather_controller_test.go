package weather

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gerps2/desafio-cloud-run/features/weather/getWeatherByCep"
	getWeatherByCepMocks "github.com/gerps2/desafio-cloud-run/features/weather/getWeatherByCep/mocks"
	httpShared "github.com/gerps2/desafio-cloud-run/shared/http"
	loggerMocks "github.com/gerps2/desafio-cloud-run/shared/logger/mocks"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func setupTestRouter(controller *WeatherController) *gin.Engine {
	gin.SetMode(gin.TestMode)
	router := gin.New()
	
	// Add the controller routes
	controller.RegisterRoutes(router)
	
	return router
}

func TestWeatherControllerGetWeatherByCepSuccess(t *testing.T) {
	// Arrange
	mockUseCase := getWeatherByCepMocks.NewMockGetWeatherByCepUseCaseInterface(t)
	mockLogger := loggerMocks.NewMockLogger(t)

	expectedResult := &getWeatherByCep.GetWeatherByCepOutput{
		TempC: 25.5,
		TempF: 77.9,
		TempK: 298.65,
	}

	// Setup mocks
	mockLogger.EXPECT().Info("GetWeatherByCep endpoint called").Once()
	mockLogger.EXPECT().Info("Weather data retrieved successfully for CEP: %s", "12345-678").Once()

	mockUseCase.EXPECT().Execute(
		mock.Anything, 
		getWeatherByCep.GetWeatherByCepInput{CepString: "12345-678"},
	).Return(expectedResult, nil).Once()

	controller := NewWeatherController(mockUseCase, mockLogger)
	router := setupTestRouter(controller)

	// Act
	req, _ := http.NewRequest("GET", "/api/v1/weather/12345-678", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	// Assert
	assert.Equal(t, http.StatusOK, w.Code)

	var response httpShared.APIResponse
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)

	assert.Equal(t, "Weather data retrieved successfully", response.Message)

	// Check the data structure
	data, ok := response.Data.(map[string]interface{})
	assert.True(t, ok, "Expected data to be a map")

	assert.Equal(t, 25.5, data["temp_C"])
	assert.Equal(t, 77.9, data["temp_F"])
	assert.Equal(t, 298.65, data["temp_K"])

	// Verify mocks
	mockUseCase.AssertExpectations(t)
	mockLogger.AssertExpectations(t)
}

func TestWeatherControllerGetWeatherByCepInvalidCEP(t *testing.T) {
	// Arrange
	mockUseCase := getWeatherByCepMocks.NewMockGetWeatherByCepUseCaseInterface(t)
	mockLogger := loggerMocks.NewMockLogger(t)

	expectedError := getWeatherByCep.NewInvalidZipcodeError()

	// Setup mocks
	mockLogger.EXPECT().Info("GetWeatherByCep endpoint called").Once()
	mockLogger.EXPECT().Error("Error executing GetWeatherByCep use case: %v", expectedError).Once()

	mockUseCase.EXPECT().Execute(
		mock.Anything,
		getWeatherByCep.GetWeatherByCepInput{CepString: "invalid-cep"},
	).Return(nil, expectedError).Once()

	controller := NewWeatherController(mockUseCase, mockLogger)
	router := setupTestRouter(controller)

	// Act
	req, _ := http.NewRequest("GET", "/api/v1/weather/invalid-cep", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	// Assert
	assert.Equal(t, http.StatusUnprocessableEntity, w.Code)

	var response httpShared.APIResponse
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)

	assert.Equal(t, "invalid zipcode", response.Message)

	// Verify mocks
	mockUseCase.AssertExpectations(t)
	mockLogger.AssertExpectations(t)
}

func TestWeatherControllerGetWeatherByCepCEPNotFound(t *testing.T) {
	// Arrange
	mockUseCase := getWeatherByCepMocks.NewMockGetWeatherByCepUseCaseInterface(t)
	mockLogger := loggerMocks.NewMockLogger(t)

	expectedError := getWeatherByCep.NewZipcodeNotFoundError()

	// Setup mocks
	mockLogger.EXPECT().Info("GetWeatherByCep endpoint called").Once()
	mockLogger.EXPECT().Error("Error executing GetWeatherByCep use case: %v", expectedError).Once()

	mockUseCase.EXPECT().Execute(
		mock.Anything,
		getWeatherByCep.GetWeatherByCepInput{CepString: "99999-999"},
	).Return(nil, expectedError).Once()

	controller := NewWeatherController(mockUseCase, mockLogger)
	router := setupTestRouter(controller)

	// Act
	req, _ := http.NewRequest("GET", "/api/v1/weather/99999-999", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	// Assert
	assert.Equal(t, http.StatusNotFound, w.Code)

	var response httpShared.APIResponse
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)

	assert.Equal(t, "can not find zipcode", response.Message)

	// Verify mocks
	mockUseCase.AssertExpectations(t)
	mockLogger.AssertExpectations(t)
}

func TestWeatherControllerGetWeatherByCepWeatherServiceError(t *testing.T) {
	// Arrange
	mockUseCase := getWeatherByCepMocks.NewMockGetWeatherByCepUseCaseInterface(t)
	mockLogger := loggerMocks.NewMockLogger(t)

	expectedError := getWeatherByCep.NewWeatherServiceError()

	// Setup mocks
	mockLogger.EXPECT().Info("GetWeatherByCep endpoint called").Once()
	mockLogger.EXPECT().Error("Error executing GetWeatherByCep use case: %v", expectedError).Once()

	mockUseCase.EXPECT().Execute(
		mock.Anything,
		getWeatherByCep.GetWeatherByCepInput{CepString: "12345-678"},
	).Return(nil, expectedError).Once()

	controller := NewWeatherController(mockUseCase, mockLogger)
	router := setupTestRouter(controller)

	// Act
	req, _ := http.NewRequest("GET", "/api/v1/weather/12345-678", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	// Assert
	assert.Equal(t, http.StatusBadGateway, w.Code)

	var response httpShared.APIResponse
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)

	// Verify mocks
	mockUseCase.AssertExpectations(t)
	mockLogger.AssertExpectations(t)
}

func TestWeatherControllerGetWeatherByCepUnknownError(t *testing.T) {
	// Arrange
	mockUseCase := getWeatherByCepMocks.NewMockGetWeatherByCepUseCaseInterface(t)
	mockLogger := loggerMocks.NewMockLogger(t)

	unknownError := errors.New("unknown error")

	// Setup mocks
	mockLogger.EXPECT().Info("GetWeatherByCep endpoint called").Once()
	mockLogger.EXPECT().Error("Error executing GetWeatherByCep use case: %v", unknownError).Once()

	mockUseCase.EXPECT().Execute(
		mock.Anything,
		getWeatherByCep.GetWeatherByCepInput{CepString: "12345-678"},
	).Return(nil, unknownError).Once()

	controller := NewWeatherController(mockUseCase, mockLogger)
	router := setupTestRouter(controller)

	// Act
	req, _ := http.NewRequest("GET", "/api/v1/weather/12345-678", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	// Assert
	assert.Equal(t, http.StatusInternalServerError, w.Code)

	var response httpShared.APIResponse
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)

	assert.Equal(t, "Failed to get weather data", response.Message)

	// Verify mocks
	mockUseCase.AssertExpectations(t)
	mockLogger.AssertExpectations(t)
}

func TestWeatherControllerGetWeatherByCepMissingCEPParameter(t *testing.T) {
	// Arrange
	mockUseCase := getWeatherByCepMocks.NewMockGetWeatherByCepUseCaseInterface(t)
	mockLogger := loggerMocks.NewMockLogger(t)

	controller := NewWeatherController(mockUseCase, mockLogger)
	router := setupTestRouter(controller)

	// Act
	req, _ := http.NewRequest("GET", "/api/v1/weather/", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	// Assert
	// This should return 404 because the route doesn't match
	assert.Equal(t, http.StatusNotFound, w.Code)

	// Use case should not be called for missing parameter
	mockUseCase.AssertNotCalled(t, "Execute")
}
