package util

import "go_eth_bot/internal/entity"

// GetAddFromMap извлечение из map файла адреса
func GetAddFromMap(usersList map[int64]string, chatID int64) (string, bool) {
	var newResp entity.CryptoUserData

	_, ok := usersList[chatID]
	if ok {
		newResp.Address = usersList[chatID] //извлечение из мапы адрес
		return newResp.Address, true
	}
	return "", false
}
