.PHONY: run test mocks wire build docker-build clean

# Executar aplicação
run:
	go run ./cmd/api

# Executar todos os testes
test:
	go test -v ./...

# Gerar mocks com Mockery
mocks:
	mockery --all

# Gerar código Wire
wire:
	wire ./cmd/api

# Build da aplicação
build: wire
	go build -o bin/api ./cmd/api

# Build Docker image para Cloud Run
docker-build:
	docker build -t weather-api:latest .

# Limpar arquivos gerados
clean:
	rm -rf bin/
	rm -f coverage.out coverage.html
	rm -f cmd/api/wire_gen.go
