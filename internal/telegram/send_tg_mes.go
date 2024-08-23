package telegram

import tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"

// SendTgMess функция отправки сообщения в ТГ
func SendTgMess(id int64, str string, bot *tgbotapi.BotAPI, page Page) {
	msg := tgbotapi.NewMessage(id, str)

	//Определение клавиатуры
	switch page {
	case First:
		msg.ReplyMarkup = FirstKeyboard
	case Second:
		msg.ReplyMarkup = SecondKeyboard
	}

	//Отправка сообщения в ТГ
	_, err := bot.Send(msg)
	if err != nil {
		return
	}
}
