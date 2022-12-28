package telegram

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
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
}

// FirstKeyboard firstKeyboard первая клавиатура для отображения в ТГ
var FirstKeyboard = tgbotapi.NewInlineKeyboardMarkup(
	tgbotapi.NewInlineKeyboardRow(
		tgbotapi.NewInlineKeyboardButtonData("📊Цена ETH", "/get_price"),
		tgbotapi.NewInlineKeyboardButtonData("⛽Цена Gas", "/get_gas"),
	),
)

// SecondKeyboard secondKeyboard вторая клавиатура для отображения в ТГ
var SecondKeyboard = tgbotapi.NewInlineKeyboardMarkup(
	tgbotapi.NewInlineKeyboardRow(
		tgbotapi.NewInlineKeyboardButtonData("🔷Баланс ETH", "/get_balance"),
		tgbotapi.NewInlineKeyboardButtonData("💲Баланс в USD", "/get_balance_usd"),
	),
	tgbotapi.NewInlineKeyboardRow(
		tgbotapi.NewInlineKeyboardButtonData("📊Цена ETH", "/get_price"),
		tgbotapi.NewInlineKeyboardButtonData("⛽Цена Gas", "/get_gas"),
	),
	tgbotapi.NewInlineKeyboardRow(
		tgbotapi.NewInlineKeyboardButtonData("🔙Другой адрес ETH", "/change_addr"),
	),
)

func New(tgApiKey string) *Updates {
	upd := &Updates{}
	// подключаемся к телеграм боту с помощью токена
	bot, err := tgbotapi.NewBotAPI(tgApiKey)
	if err != nil {
		log.Panic(err)
	}

	bot.Debug = false

	upd.updates = bot.ListenForWebhook("/" + bot.Token)

	return upd
}
