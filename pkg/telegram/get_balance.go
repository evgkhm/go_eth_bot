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

// GetBalanceRequest функция получения текущего баланса eth пользователя
func GetBalanceRequest(cfg *config.Config, address string) *big.Float {
	// godotenv package
	dotenv := cfg.EthScanApiKey

	rawAddress := "&address=" + address

	resp, httpGetErr := http.Get("https://api.etherscan.io/api" +
		"?module=account" +
		"&action=balance" +
		rawAddress +
		"&tag=latest" +
		"&apikey=" + dotenv)
	if httpGetErr != nil {
		log.Fatalln(httpGetErr)
	}

	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			log.Fatalln(err)
		}
	}(resp.Body)

	//парсинг данных, из запроса получаем WEI
	var cResp entity.CryptoUserData
	if decodeJsonErr := json.NewDecoder(resp.Body).Decode(&cResp); decodeJsonErr != nil {
		log.Fatal("ooopsss! an error occurred, please try again")
	}
	wei := new(big.Float)
	wei.SetString(cResp.Result)

	weiDivision := big.NewFloat(1000000000000000000)

	//из WEI в ETH
	ethBalance := new(big.Float).Quo(wei, weiDivision)

	return ethBalance
}

func GetBalance(ChatID int64, usersList map[int64]string, cfg *config.Config, bot *tgbotapi.BotAPI) {
	//ChatID := update.CallbackQuery.Message.Chat.ID //получаем ID пользователя
	var newResp entity.CryptoUserData
	var IsExistAddr bool
	newResp.Address, IsExistAddr = GetAddFromMap(usersList, ChatID)
	if IsExistAddr {
		ethBalance := GetBalanceRequest(cfg, newResp.Address)
		str := fmt.Sprint(ethBalance, " ETH")
		SendTgMess(ChatID, str, bot, Second)
	} else {
		str := "Некорректный адрес"
		SendTgMess(ChatID, str, bot, First)
	}
}

func GetBalanceUSD(ChatID int64, usersList map[int64]string, cfg *config.Config, bot *tgbotapi.BotAPI) {
	//ChatID := update.CallbackQuery.Message.Chat.ID
	//узнаем есть ли у этого ID адрес эфира в мапе
	var newResp entity.CryptoUserData
	var IsExistAddr bool
	newResp.Address, IsExistAddr = GetAddFromMap(usersList, ChatID)
	if IsExistAddr {
		ethBalance := GetBalanceRequest(cfg, newResp.Address)
		ethPrice := GetEthPriceRequest(cfg)
		usdBalance := new(big.Float).Mul(ethBalance, ethPrice)
		str := fmt.Sprintf("%.2f USD", usdBalance)
		SendTgMess(ChatID, str, bot, Second)
	} else {
		str := "Некорректный адрес"
		SendTgMess(ChatID, str, bot, First)
	}
}
