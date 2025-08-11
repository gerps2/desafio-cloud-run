package getWeatherByCep

import (
	"context"

	"github.com/gerps2/desafio-cloud-run/shared/domain/valueObjects"
	"github.com/gerps2/desafio-cloud-run/shared/logger"
	viacep "github.com/gerps2/desafio-cloud-run/shared/repositories/external_apis/viapcep"
	weather "github.com/gerps2/desafio-cloud-run/shared/repositories/external_apis/weather"
)

type GetWeatherByCepInput struct {
	CepString string
}

type GetWeatherByCepOutput struct {
	TempC float64 `json:"temp_C"`
	TempF float64 `json:"temp_F"`
	TempK float64 `json:"temp_K"`
}

type GetWeatherByCepUseCase interface {
	Execute(ctx context.Context, input GetWeatherByCepInput) (*GetWeatherByCepOutput, error)
}

type getWeatherByCepUseCase struct {
	viaCepRepo    viacep.ViaCepRepositoryInterface
	weatherRepo   weather.WeatherRepositoryInterface
	logger        logger.Logger
}

func NewGetWeatherByCepUseCase(viaCepRepo viacep.ViaCepRepositoryInterface, weatherRepo weather.WeatherRepositoryInterface, logger logger.Logger) GetWeatherByCepUseCaseInterface {
	return &getWeatherByCepUseCase{
		viaCepRepo:  viaCepRepo,
		weatherRepo: weatherRepo,
		logger:      logger,
	}
}

func (gwbc *getWeatherByCepUseCase) Execute(ctx context.Context, input GetWeatherByCepInput) (*GetWeatherByCepOutput, error) {
	gwbc.logger.Debug("Executing get weather by cep use case for CEP: %s", input.CepString)

	cep, err := valueObjects.NewCep(input.CepString)
	if err != nil {
		gwbc.logger.Error("Invalid CEP format: %s", input.CepString)
		return nil, NewInvalidZipcodeError()
	}

	address, err := gwbc.viaCepRepo.GetAddress(ctx, cep)
	if err != nil {
		gwbc.logger.Error("Error fetching address for CEP %s: %v", input.CepString, err)
		return nil, NewZipcodeNotFoundError()
	}

	gwbc.logger.Info("Address found for CEP %s: %s, %s", input.CepString, address.City, address.State)

	weatherData, err := gwbc.weatherRepo.GetWeather(ctx, address.City)
	if err != nil {
		gwbc.logger.Error("Error fetching weather for city %s: %v", address.City, err)
		return nil, NewWeatherServiceError()
	}

	gwbc.logger.Info("Weather data found for city %s: %.1f°C", address.City, weatherData.Current.TempC)

	tempKelvin := weatherData.Current.TempC + 273.15

	return &GetWeatherByCepOutput{
		TempC: weatherData.Current.TempC,
		TempF: weatherData.Current.TempF,
		TempK: tempKelvin,
	}, nil
}
