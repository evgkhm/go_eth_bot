package telegram

import (
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"go_eth_bot/config"
	"go_eth_bot/internal/entity"
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
			//в чат вошел новый пользователь. Поприветствуем его
			str := fmt.Sprintf(`Привет %s! Этот бот показывает стоимость эфира, газа и текущий баланс.
Для проверки баланса введите адрес кошелька`, update.Message.From.FirstName)
			SendTgMess(update.Message.Chat.ID, str, u.bot, First)
		} else if update.Message != nil {
			//если получили обычное сообщение сообщение от пользователя в ТГ
			var newResp entity.CryptoUserData
			newResp.Address = update.Message.Text
			if IsValidAddress(newResp.Address) {
				//получаем ID пользователя
				ChatID := update.Message.Chat.ID
				usersList[ChatID] = newResp.Address
				str := "Адрес получен. Выберете действие"
				SendTgMess(update.Message.Chat.ID, str, u.bot, Second)
			} else {
				newResp.Address = ""
				str := "Введите ETH адрес"
				SendTgMess(update.Message.Chat.ID, str, u.bot, First)
			}
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
				//получаем ID пользователя
				ChatID := update.CallbackQuery.Message.Chat.ID
				//удаляем из мапы пользователя
				delete(usersList, ChatID)
				str := "Введите ETH адрес"
				SendTgMess(update.CallbackQuery.Message.Chat.ID, str, u.bot, First)

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
