package telegram

import (
	"encoding/json"
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"go_eth_bot/config"
	"go_eth_bot/internal/entity"
	"io"
	"log"
	"math/big"
	"net/http"
	"strconv"
)

// GetBTCBalanceRequest получает баланс биткоин-адреса через Blockchain.com API
func GetBTCBalanceRequest(address string) *big.Float {
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

//// GetBTCPriceRequest получает текущую цену биткоина в USD через CoinMarketCap API
//func GetBTCPriceRequest(cfg *config.Config) float64 {
//	client := &http.Client{}
//	req, err := http.NewRequest("GET", "https://pro-api.coinmarketcap.com/v1/cryptocurrency/quotes/latest?symbol=BTC", nil)
//	if err != nil {
//		log.Fatalln(err)
//	}
//
//	req.Header.Add("X-CMC_PRO_API_KEY", cfg.CoinMarketCapApiKey)
//	req.Header.Add("Accept", "application/json")
//
//	resp, httpGetErr := client.Do(req)
//	if httpGetErr != nil {
//		log.Fatalln(httpGetErr)
//	}
//
//	defer func(Body io.ReadCloser) {
//		err := Body.Close()
//		if err != nil {
//			log.Fatalln(err)
//		}
//	}(resp.Body)
//
//	var cmcResponse entity.CoinMarketCapResponse
//	if decodeJsonErr := json.NewDecoder(resp.Body).Decode(&cmcResponse); decodeJsonErr != nil {
//		log.Fatal("Не удалось декодировать ответ CoinMarketCap API")
//	}
//
//	// Получаем цену BTC в USD
//	btcPrice := cmcResponse.Data["BTC"].Quote["USD"].Price
//
//	return btcPrice
//}

func GetBTCBalance(ChatID int64, usersList map[int64]string, cfg *config.Config, bot *tgbotapi.BotAPI) {
	var newResp entity.BTCUserData
	var IsExistAddr bool
	newResp.Address, IsExistAddr = GetAddFromMap(usersList, ChatID)
	if IsExistAddr {
		// Получаем баланс биткоин-адреса
		btcBalance := GetBTCBalanceRequest(newResp.Address)

		// Получаем цену BTC в USD через CoinMarketCap
		btcPrice := GetBTCPriceRequest(cfg)

		btcPriceFloat, err := strconv.ParseFloat(btcPrice, 64)
		if err != nil {
			log.Fatalln(err)
		}

		// Рассчитываем баланс в USD
		usdBalance := new(big.Float).Mul(btcBalance, big.NewFloat(btcPriceFloat))
		str := fmt.Sprintf("%.2f USD", usdBalance)

		SendTgMess(ChatID, str, bot, Second)
	} else {
		str := "Некорректный адрес"
		SendTgMess(ChatID, str, bot, First)
	}
}

func GetBTCBalanceInUSD(ChatID int64, usersList map[int64]string, cfg *config.Config, bot *tgbotapi.BotAPI) {
	var newResp entity.BTCUserData
	var IsExistAddr bool
	newResp.Address, IsExistAddr = GetAddFromMap(usersList, ChatID)
	if IsExistAddr {
		// Получаем баланс биткоин-адреса
		btcBalance := GetBTCBalanceRequest(newResp.Address)

		// Получаем цену BTC в USD через CoinMarketCap
		btcPrice := GetBTCPriceRequest(cfg)

		btcPriceFloat, err := strconv.ParseFloat(btcPrice, 64)
		if err != nil {
			log.Fatalln(err)
		}

		// Рассчитываем баланс в USD
		usdBalance := new(big.Float).Mul(btcBalance, big.NewFloat(btcPriceFloat))
		str := fmt.Sprintf("%.2f USD", usdBalance)

		SendTgMess(ChatID, str, bot, Second)
	} else {
		str := "Некорректный адрес"
		SendTgMess(ChatID, str, bot, First)
	}
}
