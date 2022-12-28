package usecase

import "go_eth_bot/pkg/telegram"

type Service struct {
	updates *telegram.Updates
}

func New(updates *telegram.Updates) *Service {
	return &Service{updates: updates}
}
