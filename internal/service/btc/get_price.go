package btc

import (
	"github.com/tidwall/gjson"
	"go_eth_bot/config"
	"io"
	"log"
	"net/http"
	"net/url"
	"strings"
)

// getBTCPriceRequest функция получения текущего курса btc
func getBTCPriceRequest(cfg *config.Config) string {
	client := &http.Client{}

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
		log.Println(httpGetErr)
	}
	defer resp.Body.Close()

	respBody, _ := io.ReadAll(resp.Body)

	value := gjson.Get(string(respBody), "data.#.quote.USD.price|0").String()
	btcPrice := removeExtn(value)

	return btcPrice
}

// removeExtn удаление символов после точки
func removeExtn(input string) string {
	if len(input) > 0 {
		if i := strings.LastIndex(input, "."); i > 0 {
			input = input[:i]
		}
	}
	return input
}

func GetBTCPrice(cfg *config.Config) string {
	//получаем цену
	btcPrice := getBTCPriceRequest(cfg)

	return btcPrice + " USD"
}
