package test

import (
	"os"
	"testing"
)

// SetupTestEnvironment configura variáveis de ambiente para testes
func SetupTestEnvironment(t *testing.T) {
	t.Helper()
	
	// Configurações básicas para testes
	os.Setenv("PORT", "8080")
	os.Setenv("HOST", "localhost")
	os.Setenv("VIACEP_BASE_URL", "https://viacep.com.br/ws")
	os.Setenv("WEATHER_BASE_URL", "http://api.weatherapi.com/v1")
	os.Setenv("WEATHER_API_KEY", "test-api-key")
	os.Setenv("REQUEST_TIMEOUT_SEC", "30")
}

// CleanupTestEnvironment limpa variáveis de ambiente após testes
func CleanupTestEnvironment(t *testing.T) {
	t.Helper()
	
	envVars := []string{
		"PORT",
		"HOST", 
		"VIACEP_BASE_URL",
		"WEATHER_BASE_URL",
		"WEATHER_API_KEY",
		"REQUEST_TIMEOUT_SEC",
	}
	
	for _, env := range envVars {
		os.Unsetenv(env)
	}
}
