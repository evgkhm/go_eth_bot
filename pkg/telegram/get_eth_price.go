package telegram

import (
	"encoding/json"
	"go_eth_bot/config"
	"go_eth_bot/internal/entity"
	"io"
	"log"
	"math/big"
	"net/http"
)

// GetEthPrice функция получения текущего курса eth
func GetEthPrice(cfg *config.Config) *big.Float {
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
