package config

import (
	"github.com/ilyakaznacheev/cleanenv"
)

type Config struct {
	TgApiKey      string `env-required:"true" env:"TG_API_KEY"`
	EthScanApiKey string `env-required:"true" env:"API_KEY"`
	Port          string `env-required:"true" env:"PORT"`
}

func NewConfig() (*Config, error) {
	cfg := &Config{}

	//cfg.tgApiKey = goDotEnvVariable("TG_API_KEY")
	//cfg.ethScanApiKey = goDotEnvVariable("TG_API_KEY")
	//goDotEnvVariable("TG_API_KEY")
	err := cleanenv.ReadEnv(cfg)
	if err != nil {
		return nil, err
	}
	return cfg, nil
}

// use godot package to load/read the .env file and
// return the value of the key
//func goDotEnvVariable(key string) string {
//	// load .env file
//	err := godotenv.Load(".env")
//
//	if err != nil {
//		errors.New("error loading .env file")
//		return ""
//	}
//
//	return os.Getenv(key)
//}
