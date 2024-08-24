package telegram

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"go_eth_bot/internal/entity"
	"regexp"
	"strings"
)

func PutAddToMap(chatID int64, usersListETH map[int64]string, usersListBTC map[int64]string, text string, bot *tgbotapi.BotAPI) {
	var newResp entity.CryptoUserData
	newResp.Address = text
	if isValidEtAddress(newResp.Address) {
		usersListETH[chatID] = newResp.Address
		str := "ETH адрес получен. Выберете действие"
		_, ok := usersListBTC[chatID]
		if ok {
			SendTgMess(chatID, str, bot, Second)
		} else {
			SendTgMess(chatID, str, bot, Third)
		}
	} else if isValidBitcoinAddress(newResp.Address) {
		usersListBTC[chatID] = newResp.Address
		str := "BTC адрес получен. Выберете действие"
		_, ok := usersListETH[chatID]
		if ok {
			SendTgMess(chatID, str, bot, Second)
		} else {
			SendTgMess(chatID, str, bot, Fourth)
		}
	} else {
		newResp.Address = ""
		str := "Введите правильный ETH или BTC адрес"
		SendTgMess(chatID, str, bot, First)
	}
}

func ChangeAddress(chatID int64, usersList map[int64]string) string {
	delete(usersList, chatID)
	return "Введите ETH или BTC адрес"
}

// isValidEtAddress функция проверки валидности eth адреса
func isValidEtAddress(v string) bool {
	re := regexp.MustCompile("^0x[\\da-fA-F]{40}$")
	return re.MatchString(v)
}

// isValidBitcoinAddress проверяет валидность биткоин-адреса
func isValidBitcoinAddress(address string) bool {
	// Регулярные выражения для разных форматов адресов
	reBase58 := regexp.MustCompile("^[13][a-km-zA-HJ-NP-Z1-9]{25,34}$") // P2PKH и P2SH адреса
	reBech32 := regexp.MustCompile("^(bc1)[a-zA-HJ-NP-Z0-9]{25,39}$")   // Bech32 адреса

	address = strings.TrimSpace(address)

	// Проверка на соответствие формату Base58 (P2PKH, P2SH) или Bech32
	return reBase58.MatchString(address) || reBech32.MatchString(address)
}
