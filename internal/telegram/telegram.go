package telegram

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"go_eth_bot/config"
	btc2 "go_eth_bot/internal/service/btc"
	eth2 "go_eth_bot/internal/service/eth"
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
	usersList := make(map[int64]string)    //здесь список всех пользователей
	usersListBTC := make(map[int64]string) //здесь список всех пользователей BTC
	for update := range u.updates {
		if update.Message != nil && update.Message.Text == "/start" {
			Greeting(update.Message.Chat.ID, update.Message.From.FirstName, u.bot)
		} else if update.Message != nil {
			PutAddToMap(update.Message.Chat.ID, usersList, usersListBTC, update.Message.Text, u.bot)
			//PutAddToMapBTC(update.Message.Chat.ID, usersListBTC, update.Message.Text, u.bot)
		}

		//если получили нажатие кнопки
		if update.CallbackQuery != nil {
			switch update.CallbackQuery.Data {
			case "/get_balance":
				eth2.GetBalance(update.CallbackQuery.Message.Chat.ID, usersList, cfg, u.bot)

			case "/get_balance_usd":
				eth2.GetBalanceUSD(update.CallbackQuery.Message.Chat.ID, usersList, cfg, u.bot)

			case "/get_price":
				eth2.GetEthPrice(update.CallbackQuery.Message.Chat.ID, usersList, cfg, u.bot)

			case "/get_gas":
				eth2.GetEthGas(update.CallbackQuery.Message.Chat.ID, usersList, cfg, u.bot)

			case "/change_addr":
				ChangeAddress(update.CallbackQuery.Message.Chat.ID, usersList, u.bot)

			case "/get_btc_price":
				resp := btc2.GetBTCPrice(cfg)
				SendTgMess(update.CallbackQuery.Message.Chat.ID, resp, u.bot, First)

			case "/get_balance_btc":
				resp := btc2.GetBTCBalance(update.CallbackQuery.Message.Chat.ID, usersListBTC)
				if resp != "" {
					SendTgMess(update.CallbackQuery.Message.Chat.ID, resp, u.bot, Second)
				} else {
					SendTgMess(update.CallbackQuery.Message.Chat.ID, "Некорректный адрес", u.bot, First)
				}
			case "/get_balance_btc_usd":
				resp := btc2.GetBTCBalanceInUSD(update.CallbackQuery.Message.Chat.ID, usersListBTC, cfg)
				if resp != "" {
					SendTgMess(update.CallbackQuery.Message.Chat.ID, resp, u.bot, Second)
				} else {
					SendTgMess(update.CallbackQuery.Message.Chat.ID, "Некорректный адрес", u.bot, First)
				}
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
