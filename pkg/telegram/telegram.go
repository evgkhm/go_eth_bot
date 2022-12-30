package telegram

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"go_eth_bot/config"
	"log"
)

// Page для выбора клавиатуры ТГ
type Page int

const (
	First Page = iota + 1
	Second
)

type Updates struct {
	updates tgbotapi.UpdatesChannel
	bot     *tgbotapi.BotAPI
}

func (u Updates) Run(cfg *config.Config) {
	usersList := make(map[int64]string) //здесь список всех пользователей

	for update := range u.updates {
		if update.Message != nil && update.Message.Text == "/start" {
			Greeting(update.Message.Chat.ID, update.Message.From.FirstName, u.bot)
		} else if update.Message != nil {
			PutAddToMap(update.Message.Chat.ID, usersList, update.Message.Text, u.bot)
		}

		//если получили нажатие кнопки
		if update.CallbackQuery != nil {
			switch update.CallbackQuery.Data {
			case "/get_balance":
				GetBalance(update.CallbackQuery.Message.Chat.ID, usersList, cfg, u.bot)

			case "/get_balance_usd":
				GetBalanceUSD(update.CallbackQuery.Message.Chat.ID, usersList, cfg, u.bot)

			case "/get_price":
				GetEthPrice(update.CallbackQuery.Message.Chat.ID, usersList, cfg, u.bot)

			case "/get_gas":
				GetEthGas(update.CallbackQuery.Message.Chat.ID, usersList, cfg, u.bot)

			case "/change_addr":
				ChangeAddress(update.CallbackQuery.Message.Chat.ID, usersList, u.bot)
			}
		}
	}
}

func New(cfg *config.Config) *Updates {
	u := &Updates{}
	// подключаемся к телеграм боту с помощью токена
	bot, err := tgbotapi.NewBotAPI(cfg.TgApiKey)
	if err != nil {
		log.Panic(err)
	}

	bot.Debug = false

	u.bot = bot
	u.updates = bot.ListenForWebhook("/" + bot.Token)

	return u
}
