package eth

import (
	"encoding/json"
	"fmt"
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

func GetEthPrice(cfg *config.Config) string {
	//получаем цену эфириума
	ethPrice := GetEthPriceRequest(cfg)
	//str := fmt.Sprint(ethPrice, " USD")
	return fmt.Sprintf("%.0f USD", ethPrice)

	//узнаем есть ли у этого ID адрес эфира в мапе
	//var newResp entity.CryptoUserData
	//var IsExistAddr bool
	//newResp.Address, IsExistAddr = util.GetAddFromMap(usersList, ChatID)
	//if IsExistAddr {
	//	return str
	//}
	//return ""
}
