package weather

import (
	"context"
	"crypto/tls"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"os"
)

type WeatherClient struct {
	BaseURL string
	APIKey  string
}

func NewClient(baseURL string, apiKey string) *WeatherClient {
	return &WeatherClient{
		BaseURL: baseURL,
		APIKey:  apiKey,
	}
}

func (c *WeatherClient) GetWeather(ctx context.Context, city string) (*WeatherResponse, error) {
	safeCity := url.QueryEscape(city)
	fullURL := fmt.Sprintf("%s%s&q=%s", c.BaseURL, c.APIKey, safeCity)

	req, err := http.NewRequestWithContext(ctx, "GET", fullURL, nil)
	if err != nil {
		return nil, err
	}

	client := &http.Client{
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{
				InsecureSkipVerify: os.Getenv("ENV") == "production",
			},
		},
	}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return nil, errors.New("failed to fetch weather data")
	}

	var weather WeatherResponse
	err = json.NewDecoder(resp.Body).Decode(&weather)
	if err != nil {
		return nil, err
	}

	return &weather, nil
}
