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
	"net/url"
)

// GetBTCPriceRequest функция получения текущего курса eth
func GetBTCPriceRequest(cfg *config.Config) *big.Float {
	client := &http.Client{}

	//resp, httpGetErr := http.Get("https://pro-api.coinmarketcap.com/v1/cryptocurrency/listings/latest?start=1&convert=USD&limit=1")

	req, reqErr := http.NewRequest("GET", "https://pro-api.coinmarketcap.com/v1/cryptocurrency/listings/latest", nil)
	if reqErr != nil {
		log.Print(reqErr)
	}
	q := url.Values{}
	q.Add("start", "1")
	q.Add("limit", "1")
	q.Add("convert", "USD")

	req.Header.Set("Accept", "application/json")
	// godotenv package
	dotenv := cfg.CoinMarketCapApiKey
	req.Header.Add("X-CMC_PRO_API_KEY", dotenv)

	req.URL.RawQuery = q.Encode()

	resp, httpGetErr := client.Do(req)
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
	var cResp entity.CryptoResponseBTC
	if err := json.NewDecoder(resp.Body).Decode(&cResp); err != nil {
		log.Fatal("error while decode data from get eth price")
	}

	btcPrice := new(big.Float)
	btcPrice.SetString(cResp.Data.Quote.USD.Price)

	return btcPrice
}

func GetBTCPrice(ChatID int64, usersListBTC map[int64]string, cfg *config.Config, bot *tgbotapi.BotAPI) {
	//получаем цену
	btcPrice := GetBTCPriceRequest(cfg)
	str := fmt.Sprint(btcPrice, " USD")

	//получаем ID пользователя
	//ChatID := update.CallbackQuery.Message.Chat.ID
	//узнаем есть ли у этого ID адрес эфира в мапе
	//var newResp entity.CryptoUserData!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!
	//var IsExistAddr bool
	//newResp.Address, IsExistAddr = GetAddFromMap(usersList, ChatID)
	//if IsExistAddr {
	//	SendTgMess(ChatID, str, bot, Second)
	//} else { //Если адреса нет вызов первой клавиатуры
	//	SendTgMess(ChatID, str, bot, First)
	//}
	SendTgMess(ChatID, str, bot, First)
}
