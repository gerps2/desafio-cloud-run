//go:build wireinject
// +build wireinject

package main

import (
	"github.com/gerps2/desafio-cloud-run/features/weather"
	"github.com/gerps2/desafio-cloud-run/shared/config"
	"github.com/gerps2/desafio-cloud-run/shared/http"
	"github.com/gerps2/desafio-cloud-run/shared/logger"
	"github.com/gerps2/desafio-cloud-run/shared/providers"

	"github.com/google/wire"
)

func InitializeApp() (*App, error) {
	wire.Build(
		// Shared dependencies
		config.Load,
		logger.New,
		http.NewServer,

		// External APIs providers
		providers.ProvideViaCepClient,
		providers.ProvideViaCepRepository,
		providers.ProvideWeatherClient,
		providers.ProvideWeatherRepository,

		// Weather feature dependencies
		weather.ProvideGetWeatherByCepUseCase,
		weather.NewWeatherController,

		// App
		NewApp,
	)
	return &App{}, nil
}
