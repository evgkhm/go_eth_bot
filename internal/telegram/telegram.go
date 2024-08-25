package telegram

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"go_eth_bot/config"
	"go_eth_bot/internal/service/btc"
	"go_eth_bot/internal/service/eth"
	"log"
)

// Page для выбора клавиатуры ТГ
type Page int

const (
	First Page = iota + 1
	Second
	Third
	Fourth
)

type Updates struct {
	updates tgbotapi.UpdatesChannel
	bot     *tgbotapi.BotAPI
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

func (u Updates) Run(cfg *config.Config) {
	usersListETH := make(map[int64]string) //здесь список всех пользователей ETH
	usersListBTC := make(map[int64]string) //здесь список всех пользователей BTC
	for update := range u.updates {
		if update.Message != nil && update.Message.Text == "/start" {
			Greeting(update.Message.Chat.ID, update.Message.From.FirstName, u.bot)
		} else if update.Message != nil {
			PutAddToMap(update.Message.Chat.ID, usersListETH, usersListBTC, update.Message.Text)
		}

		//если получили нажатие кнопки
		if update.CallbackQuery != nil {
			var keyboard Page
			currEthBalance := eth.GetBalance(update.CallbackQuery.Message.Chat.ID, usersListETH, cfg)
			currEthUSDTBalance := eth.GetBalanceUSD(update.CallbackQuery.Message.Chat.ID, usersListETH, cfg)
			currBtcBalance := btc.GetBTCBalance(update.CallbackQuery.Message.Chat.ID, usersListBTC)
			currBtcUSDTBalance := btc.GetBTCBalanceInUSD(update.CallbackQuery.Message.Chat.ID, usersListBTC, cfg)

			if currEthBalance == "" && currBtcBalance == "" {
				keyboard = First
			} else if currEthBalance != "" && currBtcBalance != "" {
				keyboard = Second
			} else if currEthBalance != "" {
				keyboard = Third
			} else {
				keyboard = Fourth
			}

			switch update.CallbackQuery.Data {
			case "/get_balance":
				SendTgMess(update.CallbackQuery.Message.Chat.ID, currEthBalance, u.bot, keyboard)

			case "/get_balance_usd":
				SendTgMess(update.CallbackQuery.Message.Chat.ID, currEthUSDTBalance, u.bot, keyboard)

			case "/get_price":
				currPrice := eth.GetEthPrice(cfg)
				SendTgMess(update.CallbackQuery.Message.Chat.ID, currPrice, u.bot, keyboard)

			case "/get_btc_price":
				currPrice := btc.GetBTCPrice(cfg)
				SendTgMess(update.CallbackQuery.Message.Chat.ID, currPrice, u.bot, keyboard)

			case "/get_balance_btc":
				SendTgMess(update.CallbackQuery.Message.Chat.ID, currBtcBalance, u.bot, keyboard)

			case "/get_balance_btc_usd":
				SendTgMess(update.CallbackQuery.Message.Chat.ID, currBtcUSDTBalance, u.bot, keyboard)

			case "/change_addr_eth":
				resp := ChangeAddress(update.CallbackQuery.Message.Chat.ID, usersListETH)
				SendTgMess(update.CallbackQuery.Message.Chat.ID, resp, u.bot, First)

			case "/change_addr_btc":
				resp := ChangeAddress(update.CallbackQuery.Message.Chat.ID, usersListBTC)
				SendTgMess(update.CallbackQuery.Message.Chat.ID, resp, u.bot, First)
			}
		}
	}
}
