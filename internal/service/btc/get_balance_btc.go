package btc

import (
	"encoding/json"
	"fmt"
	"go_eth_bot/config"
	"go_eth_bot/internal/entity"
	"go_eth_bot/internal/service/util"
	"io"
	"log"
	"math/big"
	"net/http"
	"strconv"
)

// getBTCBalanceRequest Получает баланс биткоин-адреса через Blockchain.com API
func getBTCBalanceRequest(address string) *big.Float {
	resp, httpGetErr := http.Get("https://blockchain.info/q/addressbalance/" + address)
	if httpGetErr != nil {
		log.Fatalln(httpGetErr)
	}

	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			log.Fatalln(err)
		}
	}(resp.Body)

	// Получаем баланс в сатоши (наименьшая единица биткоина)
	var satoshiBalance int64
	if decodeErr := json.NewDecoder(resp.Body).Decode(&satoshiBalance); decodeErr != nil {
		log.Fatal("Не удалось декодировать ответ Blockchain API")
	}

	// Конвертируем сатоши в BTC
	btcBalance := new(big.Float).Quo(big.NewFloat(float64(satoshiBalance)), big.NewFloat(1e8))

	return btcBalance
}

func GetBTCBalance(ChatID int64, usersList map[int64]string) string {
	var newResp entity.BTCUserData
	var IsExistAddr bool
	newResp.Address, IsExistAddr = util.GetAddFromMap(usersList, ChatID)
	if IsExistAddr {
		// Получаем баланс биткоин-адреса
		btcBalance := getBTCBalanceRequest(newResp.Address)

		str := fmt.Sprint(btcBalance, " BTC")
		return str
	}
	return ""
}

func GetBTCBalanceInUSD(ChatID int64, usersList map[int64]string, cfg *config.Config) string {
	var newResp entity.BTCUserData
	var IsExistAddr bool
	newResp.Address, IsExistAddr = util.GetAddFromMap(usersList, ChatID)
	if IsExistAddr {
		// Получаем баланс биткоин-адреса
		btcBalance := getBTCBalanceRequest(newResp.Address)

		// Получаем цену BTC в USD через CoinMarketCap
		btcPrice := GetBTCPriceRequest(cfg)

		btcPriceFloat, err := strconv.ParseFloat(btcPrice, 64)
		if err != nil {
			log.Fatalln(err)
		}

		// Рассчитываем баланс в USD
		usdBalance := new(big.Float).Mul(btcBalance, big.NewFloat(btcPriceFloat))
		str := fmt.Sprintf("%.2f USD", usdBalance)
		return str
		//telegram2.SendTgMess(ChatID, str, bot, telegram2.Second)
	}
	return ""
}
