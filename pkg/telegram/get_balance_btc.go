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
)

// GetBTCBalanceRequest функция получения текущего баланса BTC пользователя
func GetBTCBalanceRequest(cfg *config.Config, address string) *big.Float {
	// Используем BlockCypher API для получения баланса BTC
	resp, httpGetErr := http.Get("https://api.blockcypher.com/v1/btc/main/addrs/" + address + "/balance")
	if httpGetErr != nil {
		log.Fatalln(httpGetErr)
	}

	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			log.Fatalln(err)
		}
	}(resp.Body)

	// Парсинг данных, из запроса получаем баланс в сатоши
	var cResp entity.BTCUserData
	if decodeJsonErr := json.NewDecoder(resp.Body).Decode(&cResp); decodeJsonErr != nil {
		log.Fatal("ooopsss! an error occurred, please try again")
	}

	satoshis := new(big.Float)
	satoshis.SetString(fmt.Sprintf("%d", cResp.Balance))

	// Перевод из сатоши в BTC
	satoshiToBTC := big.NewFloat(1e8)
	btcBalance := new(big.Float).Quo(satoshis, satoshiToBTC)

	return btcBalance
}

// GetBTCBalanceUSD функция получения баланса в долларах США
func GetBTCBalanceUSD(cfg *config.Config, btcBalance *big.Float) *big.Float {
	// Используем CoinGecko API для получения курса BTC к USD
	resp, httpGetErr := http.Get("https://api.coingecko.com/api/v3/simple/price?ids=bitcoin&vs_currencies=usd")
	if httpGetErr != nil {
		log.Fatalln(httpGetErr)
	}

	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			log.Fatalln(err)
		}
	}(resp.Body)

	var priceData map[string]map[string]float64
	if decodeJsonErr := json.NewDecoder(resp.Body).Decode(&priceData); decodeJsonErr != nil {
		log.Fatal("ooopsss! an error occurred, please try again")
	}

	// Получаем курс BTC к USD
	btcToUSD := priceData["bitcoin"]["usd"]

	// Рассчитываем баланс в USD
	usdBalance := new(big.Float).Mul(btcBalance, big.NewFloat(btcToUSD))

	return usdBalance
}

func GetBTCBalance(ChatID int64, usersList map[int64]string, cfg *config.Config, bot *tgbotapi.BotAPI) {
	var newResp entity.BTCUserData
	var IsExistAddr bool
	newResp.Address, IsExistAddr = GetAddFromMap(usersList, ChatID)
	if IsExistAddr {
		btcBalance := GetBTCBalanceRequest(cfg, newResp.Address)
		str := fmt.Sprint(btcBalance, " BTC")
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
		btcBalance := GetBTCBalanceRequest(cfg, newResp.Address)
		usdBalance := GetBTCBalanceUSD(cfg, btcBalance)
		str := fmt.Sprintf("%.2f USD", usdBalance)
		SendTgMess(ChatID, str, bot, Second)
	} else {
		str := "Некорректный адрес"
		SendTgMess(ChatID, str, bot, First)
	}
}
