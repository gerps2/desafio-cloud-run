package providers

import (
	"github.com/gerps2/desafio-cloud-run/shared/config"
	viacep "github.com/gerps2/desafio-cloud-run/shared/repositories/external_apis/viapcep"
	"github.com/gerps2/desafio-cloud-run/shared/repositories/external_apis/weather"
)

func ProvideViaCepClient(cfg *config.Config) *viacep.ViaCepClient {
	return viacep.NewClient(cfg.ExternalAPIs.ViaCep.BaseURL)
}

func ProvideViaCepRepository(client *viacep.ViaCepClient) viacep.ViaCepRepositoryInterface {
	return viacep.NewViaCepRepository(client)
}

func ProvideWeatherClient(cfg *config.Config) *weather.WeatherClient {
	return weather.NewClient(cfg.ExternalAPIs.Weather.BaseURL, cfg.ExternalAPIs.Weather.APIKey)
}

func ProvideWeatherRepository(client *weather.WeatherClient) weather.WeatherRepositoryInterface {
	return weather.NewWeatherRepository(client)
}
