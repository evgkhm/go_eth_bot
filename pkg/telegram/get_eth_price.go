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

// GetEthPrice функция получения текущего курса eth
func GetEthPriceRequest(cfg *config.Config) *big.Float {
	// godotenv package
	dotenv := cfg.EthScanApiKey

	resp, httpGetErr := http.Get("https://api.etherscan.io/api" +
		"?module=stats" +
		"&action=ethprice" +
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

	//парсинг данных
	var cResp entity.CryptoResponsePrice
	if err := json.NewDecoder(resp.Body).Decode(&cResp); err != nil {
		log.Fatal("error while decode data from get eth price")
	}

	ethPrice := new(big.Float)
	ethPrice.SetString(cResp.Result.Ethusd)

	return ethPrice
}

func GetEthPrice(ChatID int64, usersList map[int64]string, cfg *config.Config, bot *tgbotapi.BotAPI) {
	//получаем цену эфириума
	ethPrice := GetEthPriceRequest(cfg)
	str := fmt.Sprint(ethPrice, " USD")

	//получаем ID пользователя
	//ChatID := update.CallbackQuery.Message.Chat.ID
	//узнаем есть ли у этого ID адрес эфира в мапе
	var newResp entity.CryptoUserData
	var IsExistAddr bool
	newResp.Address, IsExistAddr = GetAddFromMap(usersList, ChatID)
	if IsExistAddr {
		SendTgMess(ChatID, str, bot, Second)
	} else { //Если адреса нет вызов первой клавиатуры
		SendTgMess(ChatID, str, bot, First)
	}
}
