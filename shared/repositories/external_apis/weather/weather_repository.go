package weather

import (
	"context"
)

//go:generate mockery --name=WeatherRepositoryInterface
type WeatherRepositoryInterface interface {
	GetWeather(ctx context.Context, city string) (*WeatherResponse, error)
}

type WeatherRepository struct {
	client *WeatherClient
}

func NewWeatherRepository(client *WeatherClient) WeatherRepositoryInterface {
	return &WeatherRepository{
		client: client,
	}
}

func (r *WeatherRepository) GetWeather(ctx context.Context, city string) (*WeatherResponse, error) {
	return r.client.GetWeather(ctx, city)
}
