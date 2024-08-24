package telegram

import (
	"go_eth_bot/internal/entity"
	"regexp"
	"strings"
)

// IsValidEtAddress функция проверки валидности eth адреса
func IsValidEtAddress(v string) bool {
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

func PutAddToMap(chatID int64, usersList map[int64]string, usersListBTC map[int64]string, text string) {
	var newResp entity.CryptoUserData
	newResp.Address = text
	if IsValidEtAddress(newResp.Address) {
		usersList[chatID] = newResp.Address
		//return "ETH адрес получен. Выберете действие"
		//SendTgMess(ChatID, str, bot, Second)
	} else if IsValidBitcoinAddress(newResp.Address) {
		usersListBTC[chatID] = newResp.Address
		//return "BTC адрес получен. Выберете действие"
		//SendTgMess(ChatID, str, bot, Second)
	} else {
		newResp.Address = ""
		//return "Введите правильный адрес"
		//SendTgMess(ChatID, str, bot, First)
	}
}

func ChangeAddress(chatID int64, usersList map[int64]string) string {
	delete(usersList, chatID)
	return "Введите ETH или BTC адрес"
}
