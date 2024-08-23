package entity

type BTCUserData struct {
	Balance int64  `json:"balance"` // баланс в сатоши
	Address string `json:"address"` // биткоин-адрес
}
