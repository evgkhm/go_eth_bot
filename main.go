package main

import (
	"errors"
	"go_eth_bot/config"
	"go_eth_bot/internal/server"
	"go_eth_bot/internal/telegram"
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

	tgBot := telegram.New(cfg)

	tgBot.Run(cfg)
}
