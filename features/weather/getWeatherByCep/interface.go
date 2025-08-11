package getWeatherByCep

import (
	"context"
)

//go:generate mockery --name=GetWeatherByCepUseCaseInterface
type GetWeatherByCepUseCaseInterface interface {
	Execute(ctx context.Context, input GetWeatherByCepInput) (*GetWeatherByCepOutput, error)
}
