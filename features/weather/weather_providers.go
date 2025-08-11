package weather

import (
	"github.com/gerps2/desafio-cloud-run/features/weather/getWeatherByCep"
	"github.com/gerps2/desafio-cloud-run/shared/logger"
	viacep "github.com/gerps2/desafio-cloud-run/shared/repositories/external_apis/viapcep"
	weather "github.com/gerps2/desafio-cloud-run/shared/repositories/external_apis/weather"
)

func ProvideGetWeatherByCepUseCase(
	viaCepRepo viacep.ViaCepRepositoryInterface,
	weatherRepo weather.WeatherRepositoryInterface,
	logger logger.Logger,
) getWeatherByCep.GetWeatherByCepUseCaseInterface {
	return getWeatherByCep.NewGetWeatherByCepUseCase(viaCepRepo, weatherRepo, logger)
}
