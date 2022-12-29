package config

import (
	"errors"
	"github.com/joho/godotenv"
	"log"
	"os"
)

type Config struct {
	TgApiKey      string `env:"TG_API_KEY"`
	EthScanApiKey string `env:"API_KEY"`
	Port          string `env:"PORT"`
}

func New() (*Config, error) {
	cfg := &Config{}
	var err error
	cfg.TgApiKey, err = goDotEnvVariable("TG_API_KEY")
	cfg.EthScanApiKey, err = goDotEnvVariable("API_KEY")
	cfg.Port, err = goDotEnvVariable("PORT")
	if err != nil {
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
