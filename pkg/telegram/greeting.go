package telegram

import (
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func Greeting(ChatID int64, name string, bot *tgbotapi.BotAPI) {
	//в чат вошел новый пользователь. Поприветствуем его
	str := fmt.Sprintf(`Привет %s! Этот бот показывает стоимость эфира, газа и текущий баланс.
Для проверки баланса введите адрес кошелька`, name)
	SendTgMess(ChatID, str, bot, First)
}
