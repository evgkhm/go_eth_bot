package telegram

import (
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func Greeting(chatID int64, name string, bot *tgbotapi.BotAPI) {
	//в чат вошел новый пользователь. Поприветствуем его
	str := fmt.Sprintf(`Привет %s! Этот бот показывает стоимость криптовалют ETH, BTC и текущий баланс в кошельке.
Для проверки баланса введите адрес кошелька`, name)
	SendTgMess(chatID, str, bot, First)
}
