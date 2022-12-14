package main

import (
	"errors"
	"go_eth_bot/config"
	"go_eth_bot/pkg/server"
	"go_eth_bot/pkg/telegram"
	"log"
)

func main() {
	cfg, errConfig := config.New()
	if errConfig != nil {
		log.Fatal(errors.New("can't get config"))
	}

	serverErr := server.New(cfg)
	if serverErr != nil {
		log.Fatal(errors.New("can't create server"))
	}

	handlers := telegram.New(cfg)

	handlers.Run(cfg)
}
