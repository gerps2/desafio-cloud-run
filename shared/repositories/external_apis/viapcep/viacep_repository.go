package viacep

import (
	"context"

	"github.com/gerps2/desafio-cloud-run/shared/domain/valueObjects"
)

//go:generate mockery --name=ViaCepRepositoryInterface
type ViaCepRepositoryInterface interface {
	GetAddress(ctx context.Context, cep valueObjects.Cep) (*ViaCepResponse, error)
}

type ViaCepRepository struct {
	client *ViaCepClient
}

func NewViaCepRepository(client *ViaCepClient) ViaCepRepositoryInterface {
	return &ViaCepRepository{
		client: client,
	}
}

func (r *ViaCepRepository) GetAddress(ctx context.Context, cep valueObjects.Cep) (*ViaCepResponse, error) {
	return r.client.GetAddress(ctx, cep)
}
