package config

import (
	"log"
	"os"

	"github.com/spf13/viper"
)

type Config struct {
	Server       ServerConfig       `mapstructure:"server"`
	App          AppConfig          `mapstructure:"app"`
	ExternalAPIs ExternalAPIsConfig `mapstructure:"external_apis"`
}

type ServerConfig struct {
	Port string `mapstructure:"port"`
	Host string `mapstructure:"host"`
}

type AppConfig struct {
	Env                string `mapstructure:"env"`
	RequestTimeoutSec  int    `mapstructure:"request_timeout_sec"`
}

type ExternalAPIsConfig struct {
	ViaCep  ViaCepConfig  `mapstructure:"viacep"`
	Weather WeatherConfig `mapstructure:"weather"`
}

type ViaCepConfig struct {
	BaseURL string `mapstructure:"base_url"`
}

type WeatherConfig struct {
	BaseURL string `mapstructure:"base_url"`
	APIKey  string `mapstructure:"api_key"`
}

func Load() *Config {
	// Enable automatic environment variable reading
	viper.AutomaticEnv()

	// Set default values
	viper.SetDefault("PORT", "8080")
	viper.SetDefault("HOST", "localhost")
	viper.SetDefault("ENV", "development")
	viper.SetDefault("REQUEST_TIMEOUT_SEC", 300) // 5 minutos
	viper.SetDefault("VIACEP_BASE_URL", "https://viacep.com.br/ws/")
	viper.SetDefault("WEATHER_BASE_URL", "http://api.weatherapi.com/v1/current.json?key=")
	viper.SetDefault("WEATHER_API_KEY", "")

	// Try to load .env file if it exists
	if _, err := os.Stat(".env"); err == nil {
		viper.SetConfigFile(".env")
		viper.SetConfigType("env")
		
		if err := viper.ReadInConfig(); err != nil {
			log.Printf("Warning: Found .env file but could not read it: %v", err)
		} else {
			log.Println("Configuration loaded from .env file")
		}
	} else {
		log.Println("No .env file found, using environment variables and defaults")
	}

	var config Config
	config.Server.Port = viper.GetString("PORT")
	config.Server.Host = viper.GetString("HOST")
	config.App.Env = viper.GetString("ENV")
	config.App.RequestTimeoutSec = viper.GetInt("REQUEST_TIMEOUT_SEC")
	config.ExternalAPIs.ViaCep.BaseURL = viper.GetString("VIACEP_BASE_URL")
	config.ExternalAPIs.Weather.BaseURL = viper.GetString("WEATHER_BASE_URL")
	config.ExternalAPIs.Weather.APIKey = viper.GetString("WEATHER_API_KEY")

	return &config
}
