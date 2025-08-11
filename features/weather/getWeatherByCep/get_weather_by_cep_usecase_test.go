package getWeatherByCep

import (
	"context"
	"errors"
	"testing"

	sharedErrors "github.com/gerps2/desafio-cloud-run/shared/errors"
	loggerMocks "github.com/gerps2/desafio-cloud-run/shared/logger/mocks"
	viacep "github.com/gerps2/desafio-cloud-run/shared/repositories/external_apis/viapcep"
	viacepMocks "github.com/gerps2/desafio-cloud-run/shared/repositories/external_apis/viapcep/mocks"
	weather "github.com/gerps2/desafio-cloud-run/shared/repositories/external_apis/weather"
	weatherMocks "github.com/gerps2/desafio-cloud-run/shared/repositories/external_apis/weather/mocks"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestGetWeatherByCepUseCaseExecuteSuccess(t *testing.T) {
	// Arrange
	mockViaCepRepo := viacepMocks.NewMockViaCepRepositoryInterface(t)
	mockWeatherRepo := weatherMocks.NewMockWeatherRepositoryInterface(t)
	mockLogger := loggerMocks.NewMockLogger(t)

	expectedAddress := &viacep.ViaCepResponse{
		Cep:        "12345-678",
		Street:     "Rua Teste",
		Complement: "",
		District:   "Centro",
		City:       "São Paulo",
		State:      "SP",
		IbgeCode:   "3550308",
		GiaCode:    "1004",
		SiafiCode:  "7107",
	}

	expectedWeather := &weather.WeatherResponse{
		Location: struct {
			Name    string `json:"name"`
			Region  string `json:"region"`
			Country string `json:"country"`
		}{
			Name:    "São Paulo",
			Region:  "Sao Paulo",
			Country: "Brazil",
		},
		Current: struct {
			TempC     float64 `json:"temp_c"`
			TempF     float64 `json:"temp_f"`
			Condition struct {
				Text string `json:"text"`
			} `json:"condition"`
		}{
			TempC: 25.5,
			TempF: 77.9,
		},
	}

	mockLogger.EXPECT().Debug("Executing get weather by cep use case for CEP: %s", "12345-678").Once()
	mockLogger.EXPECT().Info("Address found for CEP %s: %s, %s", "12345-678", "São Paulo", "SP").Once()
	mockLogger.EXPECT().Info("Weather data found for city %s: %.1f°C", "São Paulo", 25.5).Once()

	mockViaCepRepo.EXPECT().GetAddress(mock.Anything, mock.AnythingOfType("valueObjects.Cep")).Return(expectedAddress, nil).Once()
	mockWeatherRepo.EXPECT().GetWeather(mock.Anything, "São Paulo").Return(expectedWeather, nil).Once()

	useCase := NewGetWeatherByCepUseCase(mockViaCepRepo, mockWeatherRepo, mockLogger)

	input := GetWeatherByCepInput{
		CepString: "12345-678",
	}

	// Act
	result, err := useCase.Execute(context.Background(), input)

	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, 25.5, result.TempC)
	assert.Equal(t, 77.9, result.TempF)
	assert.Equal(t, 298.65, result.TempK) // 25.5 + 273.15

	mockViaCepRepo.AssertExpectations(t)
	mockWeatherRepo.AssertExpectations(t)
	mockLogger.AssertExpectations(t)
}

func TestGetWeatherByCepUseCaseExecuteInvalidCEP(t *testing.T) {
	// Arrange
	mockViaCepRepo := viacepMocks.NewMockViaCepRepositoryInterface(t)
	mockWeatherRepo := weatherMocks.NewMockWeatherRepositoryInterface(t)
	mockLogger := loggerMocks.NewMockLogger(t)

	mockLogger.EXPECT().Debug("Executing get weather by cep use case for CEP: %s", "invalid-cep").Once()
	mockLogger.EXPECT().Error("Invalid CEP format: %s", "invalid-cep").Once()

	useCase := NewGetWeatherByCepUseCase(mockViaCepRepo, mockWeatherRepo, mockLogger)

	input := GetWeatherByCepInput{
		CepString: "invalid-cep",
	}

	// Act
	result, err := useCase.Execute(context.Background(), input)

	// Assert
	assert.Error(t, err)
	assert.Nil(t, result)

	apiErr, ok := err.(*sharedErrors.APIError)
	assert.True(t, ok, "Expected APIError")
	assert.Equal(t, CodeInvalidZipcode, apiErr.Code)
	assert.Equal(t, "invalid zipcode", apiErr.Message)

	mockLogger.AssertExpectations(t)
	mockViaCepRepo.AssertNotCalled(t, "GetAddress")
	mockWeatherRepo.AssertNotCalled(t, "GetWeather")
}

func TestGetWeatherByCepUseCaseExecuteAddressNotFound(t *testing.T) {
	// Arrange
	mockViaCepRepo := viacepMocks.NewMockViaCepRepositoryInterface(t)
	mockWeatherRepo := weatherMocks.NewMockWeatherRepositoryInterface(t)
	mockLogger := loggerMocks.NewMockLogger(t)

	mockLogger.EXPECT().Debug("Executing get weather by cep use case for CEP: %s", "99999-999").Once()
	mockLogger.EXPECT().Error("Error fetching address for CEP %s: %v", "99999-999", mock.AnythingOfType("*errors.errorString")).Once()

	mockViaCepRepo.EXPECT().GetAddress(mock.Anything, mock.AnythingOfType("valueObjects.Cep")).Return(nil, errors.New("address not found")).Once()

	useCase := NewGetWeatherByCepUseCase(mockViaCepRepo, mockWeatherRepo, mockLogger)

	input := GetWeatherByCepInput{
		CepString: "99999-999",
	}

	// Act
	result, err := useCase.Execute(context.Background(), input)

	// Assert
	assert.Error(t, err)
	assert.Nil(t, result)

	apiErr, ok := err.(*sharedErrors.APIError)
	assert.True(t, ok, "Expected APIError")
	assert.Equal(t, CodeZipcodeNotFound, apiErr.Code)
	assert.Equal(t, "can not find zipcode", apiErr.Message)

	mockViaCepRepo.AssertExpectations(t)
	mockLogger.AssertExpectations(t)
	mockWeatherRepo.AssertNotCalled(t, "GetWeather")
}

func TestGetWeatherByCepUseCaseExecuteWeatherServiceError(t *testing.T) {
	// Arrange
	mockViaCepRepo := viacepMocks.NewMockViaCepRepositoryInterface(t)
	mockWeatherRepo := weatherMocks.NewMockWeatherRepositoryInterface(t)
	mockLogger := loggerMocks.NewMockLogger(t)

	expectedAddress := &viacep.ViaCepResponse{
		Cep:   "12345-678",
		City:  "São Paulo",
		State: "SP",
	}

	mockLogger.EXPECT().Debug("Executing get weather by cep use case for CEP: %s", "12345-678").Once()
	mockLogger.EXPECT().Info("Address found for CEP %s: %s, %s", "12345-678", "São Paulo", "SP").Once()
	mockLogger.EXPECT().Error("Error fetching weather for city %s: %v", "São Paulo", mock.AnythingOfType("*errors.errorString")).Once()

	mockViaCepRepo.EXPECT().GetAddress(mock.Anything, mock.AnythingOfType("valueObjects.Cep")).Return(expectedAddress, nil).Once()
	mockWeatherRepo.EXPECT().GetWeather(mock.Anything, "São Paulo").Return(nil, errors.New("weather service error")).Once()

	useCase := NewGetWeatherByCepUseCase(mockViaCepRepo, mockWeatherRepo, mockLogger)

	input := GetWeatherByCepInput{
		CepString: "12345-678",
	}

	// Act
	result, err := useCase.Execute(context.Background(), input)

	// Assert
	assert.Error(t, err)
	assert.Nil(t, result)

	apiErr, ok := err.(*sharedErrors.APIError)
	assert.True(t, ok, "Expected APIError")
	assert.Equal(t, "EXTERNAL_SERVICE_ERROR", apiErr.Code)

	mockViaCepRepo.AssertExpectations(t)
	mockWeatherRepo.AssertExpectations(t)
	mockLogger.AssertExpectations(t)
}

func TestGetWeatherByCepUseCaseExecuteContextCancellation(t *testing.T) {
	// Arrange
	mockViaCepRepo := viacepMocks.NewMockViaCepRepositoryInterface(t)
	mockWeatherRepo := weatherMocks.NewMockWeatherRepositoryInterface(t)
	mockLogger := loggerMocks.NewMockLogger(t)

	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	mockLogger.EXPECT().Debug("Executing get weather by cep use case for CEP: %s", "12345-678").Once()
	mockLogger.EXPECT().Error("Error fetching address for CEP %s: %v", "12345-678", mock.AnythingOfType("*errors.errorString")).Once()

	mockViaCepRepo.EXPECT().GetAddress(mock.Anything, mock.AnythingOfType("valueObjects.Cep")).Return(nil, context.Canceled).Once()

	useCase := NewGetWeatherByCepUseCase(mockViaCepRepo, mockWeatherRepo, mockLogger)

	input := GetWeatherByCepInput{
		CepString: "12345-678",
	}

	// Act
	result, err := useCase.Execute(ctx, input)

	// Assert
	assert.Error(t, err)
	assert.Nil(t, result)

	mockViaCepRepo.AssertExpectations(t)
	mockLogger.AssertExpectations(t)
	mockWeatherRepo.AssertNotCalled(t, "GetWeather")
}
