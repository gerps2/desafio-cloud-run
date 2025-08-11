package viacep

import (
	"context"
	"crypto/tls"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"os"

	"github.com/gerps2/desafio-cloud-run/shared/domain/valueObjects"
)

type ViaCepClient struct {
	BaseURL string
}

func NewClient(baseURL string) *ViaCepClient {
	return &ViaCepClient{
		BaseURL: baseURL,
	}
}

func (c *ViaCepClient) GetAddress(ctx context.Context, cep valueObjects.Cep) (*ViaCepResponse, error) {
	req, err := http.NewRequestWithContext(ctx, "GET", fmt.Sprintf("%s%s/json/", c.BaseURL, cep.String()), nil)
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
		return nil, errors.New("failed to fetch address")
	}

	var address ViaCepResponse
	err = json.NewDecoder(resp.Body).Decode(&address)
	if err != nil {
		return nil, err
	}

	return &address, nil
}
