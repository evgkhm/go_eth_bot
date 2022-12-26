package model

// CryptoUserData содержит данные о балансе пользователя
type CryptoUserData struct {
	Address string //ETH address
	Result  string `json:"result"`
}

// CryptoResponsePrice сожержит данные о текущем курсе eth
type CryptoResponsePrice struct {
	Result struct {
		Ethusd string `json:"ethusd"`
	} `json:"result"`
}

// CryptoResponseGas содержит данные о текущем газе
type CryptoResponseGas struct {
	Result struct {
		SafeGasPrice    string `json:"SafeGasPrice"`
		ProposeGasPrice string `json:"ProposeGasPrice"`
		FastGasPrice    string `json:"FastGasPrice"`
	} `json:"result"`
}
