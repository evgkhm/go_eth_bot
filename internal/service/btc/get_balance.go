package btc

import (
	"encoding/json"
	"fmt"
	"go_eth_bot/config"
	"go_eth_bot/internal/entity"
	"go_eth_bot/internal/service/util"
	"log"
	"math/big"
	"net/http"
	"strconv"
)

// Получает баланс биткоин-адреса через Blockchain.com API
func getBTCBalanceRequest(address string) *big.Float {
	if address == "" {
		log.Println("Адрес не может быть пустым")
		return nil
	}
	//log.Println("Адрес: " + address)

	resp, httpGetErr := http.Get("https://blockchain.info/q/addressbalance/" + address)
	if httpGetErr != nil {
		log.Println(httpGetErr)
		return nil
	}
	defer resp.Body.Close()

	// Получаем баланс в сатоши (наименьшая единица биткоина)
	var satoshiBalance int64
	if decodeErr := json.NewDecoder(resp.Body).Decode(&satoshiBalance); decodeErr != nil {
		log.Println("Не удалось декодировать ответ Blockchain API", decodeErr)
		return nil
	}

	// Конвертируем сатоши в BTC
	btcBalance := new(big.Float).Quo(big.NewFloat(float64(satoshiBalance)), big.NewFloat(1e8))

	return btcBalance
}

func GetBTCBalance(chatID int64, usersList map[int64]string) string {
	var newResp entity.BTCUserData
	var IsExistAddr bool
	newResp.Address, IsExistAddr = util.GetAddFromMap(usersList, chatID)
	if IsExistAddr {
		// Получаем баланс биткоин-адреса
		btcBalance := getBTCBalanceRequest(newResp.Address)

		str := fmt.Sprint(btcBalance, " BTC")
		return str
	}
	return ""
}

func GetBTCBalanceInUSD(currBtcBalance string, cfg *config.Config) string {
	if currBtcBalance == "" {
		return "0 USD"
	}

	// из string во float64
	btcBalanceFloat, err := strconv.ParseFloat(currBtcBalance, 64)
	if err != nil {
		log.Println(err)
		return "0 USD"
	}

	// Получаем цену BTC в USD через CoinMarketCap
	btcPrice := getBTCPriceRequest(cfg)

	btcPriceFloat, err := strconv.ParseFloat(btcPrice, 64)
	if err != nil {
		log.Println(err)
		return "0 USD"
	}

	// Рассчитываем баланс в USD
	usdBalance := new(big.Float).Mul(big.NewFloat(btcBalanceFloat), big.NewFloat(btcPriceFloat))
	str := fmt.Sprintf("%.2f USD", usdBalance)
	return str
}
