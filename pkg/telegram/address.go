package telegram

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"go_eth_bot/internal/entity"
	"regexp"
	"strings"
)

// IsValidAddress функция проверки валидности eth адреса
func IsValidAddress(v string) bool {
	re := regexp.MustCompile("^0x[\\da-fA-F]{40}$")
	return re.MatchString(v)
}

// IsValidBitcoinAddress проверяет валидность биткоин-адреса
func IsValidBitcoinAddress(address string) bool {
	// Регулярные выражения для разных форматов адресов
	reBase58 := regexp.MustCompile("^[13][a-km-zA-HJ-NP-Z1-9]{25,34}$") // P2PKH и P2SH адреса
	reBech32 := regexp.MustCompile("^(bc1)[a-zA-HJ-NP-Z0-9]{25,39}$")   // Bech32 адреса

	address = strings.TrimSpace(address)

	// Проверка на соответствие формату Base58 (P2PKH, P2SH) или Bech32
	return reBase58.MatchString(address) || reBech32.MatchString(address)
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

func PutAddToMap(ChatID int64, usersList map[int64]string, usersListBTC map[int64]string, text string, bot *tgbotapi.BotAPI) {
	var newResp entity.CryptoUserData
	newResp.Address = text
	if IsValidAddress(newResp.Address) {
		usersList[ChatID] = newResp.Address
		str := "ETH адрес получен. Выберете действие"
		SendTgMess(ChatID, str, bot, Second)
	} else if IsValidBitcoinAddress(newResp.Address) {
		usersListBTC[ChatID] = newResp.Address
		str := "BTC адрес получен. Выберете действие"
		SendTgMess(ChatID, str, bot, Second)
	} else {
		newResp.Address = ""
		str := "Введите ETH адрес"
		SendTgMess(ChatID, str, bot, First)
	}
}

func ChangeAddress(ChatID int64, usersList map[int64]string, bot *tgbotapi.BotAPI) {
	delete(usersList, ChatID)
	str := "Введите ETH или BTC адрес"
	SendTgMess(ChatID, str, bot, First)
}
