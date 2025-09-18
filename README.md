# Desafio Cloud Run - Weather API

API REST desenvolvida em Go seguindo **Clean Architecture** + **DDD** com **Vertical Slice Architecture**. A aplicação consulta o clima de uma cidade através do CEP, integrando-se com as APIs ViaCep e WeatherAPI.

## 🚀 Como Executar

> **⚡ Quer testar rapidamente?** Use Docker Compose (mais fácil):
> ```bash
> git clone <repository-url>
> cd desafio-cloud-run
> echo "WEATHER_API_KEY=sua-chave-aqui" > .env
> docker-compose up --build
> ```

> **🌐 API já está rodando no Cloud Run:** https://desafio-cloud-run-308065563700.us-central1.run.app

### Pré-requisitos

#### Obrigatórios
- **Go 1.22+**: [Download aqui](https://golang.org/dl/)
- **WeatherAPI Key**: Conta gratuita em [WeatherAPI](https://www.weatherapi.com/) (1000 chamadas/dia)

#### Opcionais (para desenvolvimento)
- **Google Wire CLI**: `go install github.com/google/wire/cmd/wire@latest`
- **Mockery**: `go install github.com/vektra/mockery/v2@latest`
- **Docker**: Para execução containerizada
- **Docker Compose**: Para orquestração de containers

### 🚀 Execução Rápida (Docker Compose)

1. **Clone o repositório**
```bash
git clone <repository-url>
cd desafio-cloud-run
```

2. **Configure a Weather API Key**
```bash
# Crie um arquivo .env com sua chave da WeatherAPI
echo "WEATHER_API_KEY=sua-chave-weather-api-aqui" > .env
```

3. **Execute com Docker Compose**
```bash
# Build e execução
docker-compose up --build

# Ou em background
docker-compose up -d --build
```

4. **Acesse a aplicação**
```bash
# API disponível em
curl "http://localhost:8080/api/v1/weather/18077346"

# Health check
curl "http://localhost:8080/health"
```

### 🖥️ Execução Local (Desenvolvimento)

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
# Crie o arquivo .env baseado no exemplo
cp .env.example .env
# Edite o .env com sua Weather API Key
```

4. **Gere o código Wire e mocks** (opcional)
```bash
make wire
make mocks
```

5. **Execute a aplicação**
```bash
make run
# ou
go run ./cmd/api
```

A aplicação estará disponível em `http://localhost:8080`

### 🐳 Outras Formas de Execução

#### Docker (Standalone)
```bash
# Build da imagem
docker build -t weather-api:latest .

# Execute o container
docker run -p 8080:8080 \
  -e WEATHER_API_KEY=sua-chave-aqui \
  -e PORT=8080 \
  -e HOST=0.0.0.0 \
  weather-api:latest
```

#### Comandos do Makefile
```bash
# Executar aplicação
make run

# Executar todos os testes
make test

# Gerar mocks
make mocks

# Gerar código Wire
make wire

# Build da aplicação
make build

# Limpar arquivos gerados
make clean
```

## 📋 Configuração

### Arquivo de Configuração (.env)

Crie um arquivo `.env` na raiz do projeto com as seguintes configurações:

```env
# ===========================================
# CONFIGURAÇÃO DO SERVIDOR
# ===========================================
PORT=8080
HOST=localhost
ENV=development

# ===========================================
# APIs EXTERNAS
# ===========================================
# ViaCep API (endereço por CEP)
VIACEP_BASE_URL=https://viacep.com.br/ws/

# WeatherAPI (clima por cidade)
WEATHER_BASE_URL=http://api.weatherapi.com/v1/current.json?key=
WEATHER_API_KEY=sua-chave-weather-api-aqui

# ===========================================
# CONFIGURAÇÕES DA APLICAÇÃO
# ===========================================
# Timeout das requisições em segundos (padrão: 300 = 5 minutos)
REQUEST_TIMEOUT_SEC=300
```

#### Como obter a Weather API Key

1. Acesse [WeatherAPI](https://www.weatherapi.com/)
2. Crie uma conta gratuita
3. Vá para "My Account" → "API Keys"
4. Copie sua API Key
5. Adicione ao arquivo `.env`:
```bash
WEATHER_API_KEY=sua-chave-real-aqui
```

> **Nota**: A conta gratuita permite 1000 chamadas por dia, o que é suficiente para desenvolvimento e testes.

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

1. **📥 Entrada HTTP**: Cliente faz requisição `GET /api/v1/weather/01310-100`

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

- **Go 1.22+**: Linguagem principal
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

## 📡 API Endpoints

### 🧪 Arquivos de Teste da API

O projeto inclui arquivos prontos para testar a API:

- **`api-dev.http`** - Para testar a API em produção (Cloud Run)
- **`api.http`** - Para testar a API localmente (desenvolvimento)

> **💡 Como usar:** Abra estes arquivos no VS Code com a extensão "REST Client" e clique em "Send Request" acima de cada requisição!

#### Exemplo de uso:

1. **Abra o arquivo** `api-dev.http` no VS Code
2. **Instale a extensão** REST Client (se não tiver)
3. **Clique em "Send Request"** acima de cada endpoint
4. **Veja as respostas** diretamente no editor!

```http
### Health Check (Produção)
GET https://desafio-cloud-run-308065563700.us-central1.run.app/health

### Weather API (Produção)
GET https://desafio-cloud-run-308065563700.us-central1.run.app/api/v1/weather/01310-100
```

### Weather API

#### Consultar Clima por CEP
```http
GET /api/v1/weather/{cep}
```

**Parâmetros:**
- `cep` (path parameter): CEP no formato `00000-000` ou `00000000`

**Exemplo de Requisição:**
```bash
curl "http://localhost:8080/api/v1/weather/01310-100"
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

#### Todos os Testes
```bash
# Via Makefile (recomendado)
make test

# Ou diretamente com Go
go test -v ./...
```

#### Testes com Cobertura
```bash
# Gerar relatório de cobertura
go test -coverprofile=coverage.out ./...

# Visualizar cobertura no terminal
go tool cover -func=coverage.out

# Gerar relatório HTML
go tool cover -html=coverage.out -o coverage.html

# Abrir relatório no navegador
open coverage.html
```

#### Testes Específicos
```bash
# Apenas testes unitários
go test ./features/weather/getWeatherByCep/...

# Apenas testes de controllers
go test ./features/weather/...

# Testes com verbose e cobertura
go test -v -cover ./...
```

#### Gerenciamento de Mocks
```bash
# Regenerar todos os mocks
make mocks

# Ou manualmente
mockery --all

# Verificar se mocks estão atualizados
git status
```

### Estrutura de Testes

A aplicação possui cobertura de testes abrangente:

- **✅ Unitários**: Use Cases, Value Objects, Repositories
- **✅ Integração**: Controllers, HTTP endpoints
- **✅ Mocks**: Gerados automaticamente com Mockery
- **✅ Helpers**: Utilitários para configuração de testes

### Cobertura Atual

| Componente | Cobertura | Status |
|------------|-----------|--------|
| **Use Cases** | ~85% | ✅ Bom |
| **Controllers** | ~75% | ✅ Adequado |
| **Value Objects** | 100% | ✅ Completo |
| **Repositories** | ~80% | ✅ Bom |

### Testes E2E

Para testes end-to-end (se implementados):

```bash
# Executar aplicação em background
make run &

# Aguardar inicialização
sleep 5

# Executar testes E2E
go test ./test/e2e/...

# Parar aplicação
pkill -f "go run"
```

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
FROM golang:1.22-alpine AS builder
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

## 🔗 Links Úteis

- **API Produção**: https://desafio-cloud-run-308065563700.us-central1.run.app
- **Documentação WeatherAPI**: https://www.weatherapi.com/docs/
- **Documentação ViaCep**: https://viacep.com.br/

---

**🚀 Desenvolvido com ❤️ usando Go, Clean Architecture e melhores práticas de desenvolvimento!**

