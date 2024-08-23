package config

import (
	"errors"
	"github.com/joho/godotenv"
	"log"
	"os"
)

type Config struct {
	TgApiKey            string `env:"TG_API_KEY"`
	EthScanApiKey       string `env:"API_KEY"`
	Port                string `env:"PORT"`
	CoinMarketCapApiKey string `env:"CMC_API_KEY"`
}

func New() (*Config, error) {
	cfg := &Config{}
	var err error
	cfg.TgApiKey = os.Getenv("TG_API_KEY")
	if cfg.TgApiKey == "" {
		return nil, err
	}

	cfg.EthScanApiKey = os.Getenv("API_KEY")
	if cfg.EthScanApiKey == "" {
		return nil, err
	}

	cfg.Port = os.Getenv("PORT")
	if cfg.Port == "" {
		return nil, err
	}

	cfg.CoinMarketCapApiKey = os.Getenv("CMC_API_KEY")
	if cfg.CoinMarketCapApiKey == "" {
		return nil, err
	}

	return cfg, nil
}

// use godot package to load/read the .env file and return the value of the key
func goDotEnvVariable(key string) (string, error) {
	// load .env file
	err := godotenv.Load(".env")

	if err != nil {
		log.Fatal(errors.New("error loading .env file"))
		return "", err
	}

	return os.Getenv(key), nil
}
