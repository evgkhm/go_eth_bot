package telegram

import (
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
