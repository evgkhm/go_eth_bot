package telegram

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"go_eth_bot/internal/entity"
	"regexp"
)

// IsValidAddress функция проверки валидности eth адреса
func IsValidAddress(v string) bool {
	re := regexp.MustCompile("^0x[\\da-fA-F]{40}$")
	return re.MatchString(v)
}

// GetAddFromMap извлечение из map файла эфир адреса
func GetAddFromMap(usersList map[int64]string, chatID int64) (string, bool) {
	var newResp entity.CryptoUserData

	_, ok := usersList[chatID]
	if ok {
		newResp.Address = usersList[chatID] //извлечение из мапы адрес эфира

		return newResp.Address, true
	}
	return "", false
}

func PutAddToMap(ChatID int64, usersList map[int64]string, text string, bot *tgbotapi.BotAPI) {
	var newResp entity.CryptoUserData
	newResp.Address = text
	if IsValidAddress(newResp.Address) {
		usersList[ChatID] = newResp.Address
		str := "Адрес получен. Выберете действие"
		SendTgMess(ChatID, str, bot, Second)
	} else {
		newResp.Address = ""
		str := "Введите ETH адрес"
		SendTgMess(ChatID, str, bot, First)
	}
}

func ChangeAddress(ChatID int64, usersList map[int64]string, bot *tgbotapi.BotAPI) {
	delete(usersList, ChatID)
	str := "Введите ETH адрес"
	SendTgMess(ChatID, str, bot, First)
}
