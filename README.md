# Desafio Cloud Run - Weather API

API REST desenvolvida em Go seguindo **Clean Architecture** + **DDD** com **Vertical Slice Architecture**. A aplicação consulta o clima de uma cidade através do CEP, integrando-se com as APIs ViaCep e WeatherAPI.

## 🏗️ Arquitetura da Aplicação

### Visão Geral

A aplicação segue uma arquitetura em camadas bem definidas, onde cada componente tem uma responsabilidade específica e as dependências fluem sempre para dentro (Dependency Inversion Principle).

```
┌─────────────────┐    ┌─────────────────┐    ┌─────────────────┐
│   HTTP Client   │    │   HTTP Client   │    │   HTTP Client   │
│   (Browser)     │    │   (Postman)     │    │   (Mobile)      │
└─────────┬───────┘    └─────────┬───────┘    └─────────┬───────┘
          │                      │                      │
          └──────────────────────┼──────────────────────┘
                                 │
                    ┌─────────────▼─────────────┐
                    │      HTTP Server         │
                    │    (Gin Framework)       │
                    └─────────────┬─────────────┘
                                 │
                    ┌─────────────▼─────────────┐
                    │     Controllers          │
                    │  (HTTP Request/Response) │
                    └─────────────┬─────────────┘
                                 │
                    ┌─────────────▼─────────────┐
                    │      Use Cases           │
                    │   (Business Logic)       │
                    └─────────────┬─────────────┘
                                 │
          ┌──────────────────────┼──────────────────────┐
          │                      │                      │
┌─────────▼─────────┐  ┌─────────▼─────────┐  ┌─────────▼─────────┐
│   ViaCep Repo     │  │  Weather Repo     │  │   Value Objects   │
│  (CEP → Address)  │  │ (City → Weather)  │  │   (CEP, etc.)     │
└─────────┬─────────┘  └─────────┬─────────┘  └───────────────────┘
          │                      │
┌─────────▼─────────┐  ┌─────────▼─────────┐
│   ViaCep API      │  │  WeatherAPI       │
│  (External)       │  │  (External)       │
└───────────────────┘  └───────────────────┘
```

### Estrutura do Projeto

```
desafio-cloud-run/
├── cmd/api/                           # 🚀 Entry Point
│   ├── main.go                       # Inicialização da aplicação
│   ├── wire.go                       # Configuração de DI (Wire)
│   └── wire_gen.go                   # Código gerado pelo Wire
│
├── features/                          # 🎯 Features (Vertical Slices)
│   └── weather/                      # Feature de consulta de clima
│       ├── weather_controller.go     # HTTP Controllers
│       ├── weather_routes.go         # Definição de rotas
│       └── getWeatherByCep/          # Use Case específico
│           ├── get_weather_by_cep_usecase.go
│           ├── types.go              # DTOs de entrada/saída
│           └── errors.go             # Erros específicos do domínio
│
├── shared/                            # 🔧 Componentes Compartilhados
│   ├── config/                       # Configurações (Viper)
│   ├── logger/                       # Sistema de logs estruturado
│   ├── http/                         # Servidor HTTP e middlewares
│   ├── errors/                       # Tratamento global de erros
│   ├── domain/valueObjects/          # Value Objects (CEP, etc.)
│   └── repositories/external_apis/   # Integrações externas
│       ├── viapcep/                  # Cliente ViaCep
│       └── weather/                  # Cliente WeatherAPI
│
├── test/                             # 🧪 Testes
│   ├── e2e/                         # Testes End-to-End
│   └── mocks/                       # Mocks gerados (Mockery)
│
├── .env                              # Variáveis de ambiente
├── Makefile                          # Comandos de build/test
└── go.mod                           # Dependências Go
```

## 🔄 Fluxo da Aplicação

### Como Funciona uma Requisição

1. **📥 Entrada HTTP**: Cliente faz requisição `GET /api/v1/weather/getWeatherByCep?cep=01310-100`

2. **🎛️ Controller**: Recebe a requisição, extrai o parâmetro CEP e chama o Use Case

3. **💼 Use Case (Business Logic)**:
   - Valida o formato do CEP usando Value Object
   - Consulta o endereço via **ViaCep Repository**
   - Consulta o clima via **Weather Repository** usando a cidade encontrada
   - Converte temperaturas (Celsius, Fahrenheit, Kelvin)
   - Retorna resposta padronizada

4. **🔌 Repositories**: Fazem chamadas HTTP para APIs externas
   - **ViaCep**: `https://viacep.com.br/ws/{cep}/json/`
   - **WeatherAPI**: `http://api.weatherapi.com/v1/current.json?key={key}&q={city}`

5. **📤 Resposta**: Controller retorna JSON padronizado ao cliente

### Exemplo de Fluxo Completo

```
Cliente → Controller → Use Case → ViaCep Repo → ViaCep API
                         ↓
                    Weather Repo → WeatherAPI
                         ↓
                    Resposta Final ← Controller ← Use Case
```

## 🛠️ Tecnologias e Padrões

### Stack Tecnológica

- **Go 1.21+**: Linguagem principal
- **Gin**: Framework HTTP
- **Google Wire**: Injeção de dependência
- **Viper**: Gerenciamento de configuração
- **Testify**: Framework de testes
- **Mockery**: Geração automática de mocks

### Padrões Arquiteturais

- **Clean Architecture**: Separação em camadas com dependências invertidas
- **DDD**: Value Objects, Entities, Repositories
- **Vertical Slice**: Features organizadas por funcionalidade completa
- **Repository Pattern**: Abstração para acesso a dados externos
- **Dependency Injection**: Inversão de controle via Wire
- **Error Handling**: Tratamento padronizado de erros

### Middleware Global

- **Timeout**: Configurável via `REQUEST_TIMEOUT_SEC` (padrão: 300s)
- **Recovery**: Captura panics e retorna erro 500
- **CORS**: Configurado para desenvolvimento
- **Logging**: Log estruturado de todas as requisições

## 📋 Configuração

### Variáveis de Ambiente (.env)

```env
# Servidor HTTP
PORT=8080
HOST=localhost

# APIs Externas
VIACEP_BASE_URL=https://viacep.com.br/ws/
WEATHER_BASE_URL=http://api.weatherapi.com/v1/current.json?key=
WEATHER_API_KEY=your-weather-api-key

# Configurações
REQUEST_TIMEOUT_SEC=300
```

## 🚀 Como Executar

### Pré-requisitos

- Go 1.21+
- Google Wire CLI: `go install github.com/google/wire/cmd/wire@latest`
- Mockery: `go install github.com/vektra/mockery/v2@latest`
- Conta no WeatherAPI para obter API key

### Comandos Disponíveis

```bash
# Executar aplicação
make run

# Executar todos os testes
make test

# Executar testes E2E
make test-e2e

# Gerar mocks
make mocks

# Gerar código Wire
make wire
```

### Instalação e Execução

1. **Clone o repositório**
```bash
git clone <repository-url>
cd desafio-cloud-run
```

2. **Instale as dependências**
```bash
go mod tidy
```

3. **Configure as variáveis de ambiente**
```bash
cp .env.example .env
# Edite o .env com suas configurações
```

4. **Gere o código Wire e mocks**
```bash
make wire
make mocks
```

5. **Execute a aplicação**
```bash
make run
```

A aplicação estará disponível em `http://localhost:8080`

## 📡 API Endpoints

### Weather API

#### Consultar Clima por CEP
```http
GET /api/v1/weather/getWeatherByCep?cep={cep}
```

**Parâmetros:**
- `cep` (string): CEP no formato `00000-000` ou `00000000`

**Exemplo de Requisição:**
```bash
curl "http://localhost:8080/api/v1/weather/getWeatherByCep?cep=01310-100"
```

**Resposta de Sucesso (200):**
```json
{
  "message": "Weather data retrieved successfully",
  "data": {
    "temp_C": 23.5,
    "temp_F": 74.3,
    "temp_K": 296.65
  }
}
```

**Respostas de Erro:**

**CEP Inválido (422):**
```json
{
  "message": "invalid zipcode"
}
```

**CEP Não Encontrado (404):**
```json
{
  "message": "can not find zipcode"
}
```

**Erro do Serviço (502):**
```json
{
  "message": "Weather service temporarily unavailable"
}
```

### Health Check

#### Verificar Status da API
```http
GET /health
```

**Resposta:**
```json
{
  "status": "OK",
  "timestamp": "2024-01-01T12:00:00Z"
}
```

## 🧪 Testes

### Estrutura de Testes

A aplicação possui uma cobertura de testes robusta com **80%+ de cobertura**:

- **Testes Unitários**: Use Cases, Value Objects, Repositories
- **Testes de Integração**: Controllers, HTTP endpoints
- **Mocks**: Gerados automaticamente com Mockery

### Executar Testes

```bash
# Todos os testes
make test

# Testes com cobertura
go test -coverprofile=coverage.out ./...
go tool cover -html=coverage.out -o coverage.html

# Regenerar mocks
make mocks
```

### Cobertura Atual

- **Use Cases**: 87.5% de cobertura
- **Controllers**: 75.0% de cobertura  
- **Value Objects**: 100% de cobertura

## 🏗️ Desenvolvimento

### Adicionando Nova Feature

Para adicionar uma nova feature seguindo a arquitetura:

1. **Criar estrutura da feature**
```bash
mkdir -p features/nova-feature
mkdir -p features/nova-feature/usecase
```

2. **Implementar componentes**
- Controller (HTTP layer)
- Use Case (Business logic)
- Routes (Endpoint definitions)
- Types (DTOs)
- Errors (Domain errors)

3. **Atualizar Wire**
```bash
# Adicionar providers no wire.go
make wire
```

4. **Criar testes**
- Testes unitários para Use Cases
- Testes de integração para Controllers
- Mocks para dependências externas

### Padrões de Código

- **Interfaces**: Definidas no mesmo arquivo da implementação
- **Errors**: Específicos por domínio com códigos HTTP apropriados
- **Logging**: Estruturado com níveis (Debug, Info, Error)
- **Context**: Propagado em todas as operações para timeout/cancelamento
- **DTOs**: Input/Output tipados para Use Cases

## 🚀 Deploy

### Docker (Futuro)

```dockerfile
FROM golang:1.21-alpine AS builder
WORKDIR /app
COPY . .
RUN go mod download
RUN go build -o main ./cmd/api

FROM alpine:latest
RUN apk --no-cache add ca-certificates
WORKDIR /root/
COPY --from=builder /app/main .
CMD ["./main"]
```

### Cloud Run (Google Cloud)

```bash
# Build e deploy
gcloud builds submit --tag gcr.io/PROJECT-ID/weather-api
gcloud run deploy --image gcr.io/PROJECT-ID/weather-api --platform managed
```

## 📚 Referências

- [Clean Architecture - Uncle Bob](https://blog.cleancoder.com/uncle-bob/2012/08/13/the-clean-architecture.html)
- [Domain Driven Design](https://martinfowler.com/bliki/DomainDrivenDesign.html)
- [Vertical Slice Architecture](https://jimmybogard.com/vertical-slice-architecture/)
- [Google Wire](https://github.com/google/wire)
- [Gin Framework](https://gin-gonic.com/)

---

**Desenvolvido com ❤️ usando Go e Clean Architecture**
```

O servidor estará disponível em `http://localhost:8080`

## 📡 Endpoints

### Health Check
```
GET /health
```
Resposta:
```json
{"status":"healthy"}
```

### Ping
```
GET /api/v1/ping
```
Resposta:
```json
{"message":"pong","status":"success"}
```

## ⚙️ Configuração

As configurações são gerenciadas através do arquivo `.env`:

```env
PORT=8080
HOST=localhost
ENV=development
```

## 🛠️ Tecnologias Utilizadas

- **[Gin](https://github.com/gin-gonic/gin)**: Framework web HTTP
- **[Viper](https://github.com/spf13/viper)**: Gerenciamento de configuração
- **[Wire](https://github.com/google/wire)**: Injeção de dependência
- **Go 1.21**: Linguagem de programação

## 📁 Adicionando Novas Features

Para adicionar uma nova feature seguindo o padrão vertical slice:

1. Crie um diretório em `features/nome-da-feature/`
2. Implemente as camadas:
   - `domain/`: Entidades e regras de negócio
   - `usecase/`: Casos de uso
   - `controller/`: Controllers HTTP
   - `routes/`: Definição de rotas
3. Registre as dependências no `wire.go`
4. Registre as rotas no `main.go`

## 🧪 Testes

Para executar os testes:
```bash
go test ./...
```

## 📝 Logs

O sistema de logs está configurado para diferentes níveis:
- INFO: Informações gerais
- ERROR: Erros
- DEBUG: Informações de debug
- WARN: Avisos

## 🔧 Desenvolvimento

Para desenvolvimento, o servidor roda em modo debug. Para produção, configure:
```bash
export GIN_MODE=release
```
