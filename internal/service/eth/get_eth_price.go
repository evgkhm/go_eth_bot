package eth

import (
	"encoding/json"
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"go_eth_bot/config"
	"go_eth_bot/internal/entity"
	telegram2 "go_eth_bot/internal/telegram"
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
	//str := fmt.Sprint(ethPrice, " USD")
	str := fmt.Sprintf("%.0f USD", ethPrice)

	//получаем ID пользователя
	//ChatID := update.CallbackQuery.Message.Chat.ID
	//узнаем есть ли у этого ID адрес эфира в мапе
	var newResp entity.CryptoUserData
	var IsExistAddr bool
	newResp.Address, IsExistAddr = telegram2.GetAddFromMap(usersList, ChatID)
	if IsExistAddr {
		telegram2.SendTgMess(ChatID, str, bot, telegram2.Second)
	} else { //Если адреса нет вызов первой клавиатуры
		telegram2.SendTgMess(ChatID, str, bot, telegram2.First)
	}
}
