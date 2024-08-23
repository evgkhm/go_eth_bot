package telegram

import tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"

// FirstKeyboard firstKeyboard первая клавиатура для отображения в ТГ
var FirstKeyboard = tgbotapi.NewInlineKeyboardMarkup(
	tgbotapi.NewInlineKeyboardRow(
		tgbotapi.NewInlineKeyboardButtonData("📊Цена ETH", "/get_price"),
		tgbotapi.NewInlineKeyboardButtonData("📊Цена BTC", "/get_btc_price"),
		//tgbotapi.NewInlineKeyboardButtonData("⛽Цена Gas", "/get_gas"),
	),
)

// SecondKeyboard secondKeyboard вторая клавиатура для отображения в ТГ
var SecondKeyboard = tgbotapi.NewInlineKeyboardMarkup(
	tgbotapi.NewInlineKeyboardRow(
		tgbotapi.NewInlineKeyboardButtonData("🔷Баланс ETH", "/get_balance"),
		tgbotapi.NewInlineKeyboardButtonData("💲Баланс в USD", "/get_balance_usd"),
	),
	tgbotapi.NewInlineKeyboardRow(
		tgbotapi.NewInlineKeyboardButtonData("₿ Баланс BTC", "/get_balance_btc"),
		tgbotapi.NewInlineKeyboardButtonData("💲Баланс в USD", "/get_balance_btc_usd"),
	),
	tgbotapi.NewInlineKeyboardRow(
		tgbotapi.NewInlineKeyboardButtonData("📊Цена ETH", "/get_price"),
		tgbotapi.NewInlineKeyboardButtonData("📊Цена BTC", "/get_btc_price"),
	),
	tgbotapi.NewInlineKeyboardRow(
		tgbotapi.NewInlineKeyboardButtonData("🔙Другой адрес ETH", "/change_addr"),
	),
)
