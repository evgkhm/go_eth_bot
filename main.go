package main

import (
	"errors"
	"go_eth_bot/config"
	"go_eth_bot/pkg/server"
	"go_eth_bot/pkg/telegram"
)

func main() {
	cfg, errConfig := config.New()
	if errConfig != nil {
		errors.New("can't get config")
	}

	serverErr := server.New(cfg.Port)
	if serverErr != nil {
		errors.New("can't create server")
	}

	handlers := telegram.New(cfg.TgApiKey)

	handlers.Run(cfg)
}
