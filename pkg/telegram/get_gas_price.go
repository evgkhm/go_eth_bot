package telegram

import (
	"encoding/json"
	"go_eth_bot/config"
	"go_eth_bot/internal/entity"
	"io"
	"log"
	"net/http"
	"strconv"
)

// GetGasPrice функция получения текущего газа сети eth
func GetGasPrice(cfg *config.Config) (uint32, uint32, uint32) {
	// godotenv package
	dotenv := cfg.EthScanApiKey

	resp, httpGetErr := http.Get("https://api.etherscan.io/api" +
		"?module=gastracker" +
		"&action=gasoracle" +
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

	//Decode the data
	var cResp entity.CryptoResponseGas
	if err := json.NewDecoder(resp.Body).Decode(&cResp); err != nil {
		log.Fatal("error while decode data from get eth price")
	}

	//From string to uint32
	safeGasPrice, err := strconv.ParseUint(cResp.Result.SafeGasPrice, 10, 32)
	if err != nil {
		panic(err)
	}
	proposeGasPrice, err := strconv.ParseUint(cResp.Result.ProposeGasPrice, 10, 32)
	if err != nil {
		panic(err)
	}
	fastGasPrice, err := strconv.ParseUint(cResp.Result.FastGasPrice, 10, 32)
	if err != nil {
		panic(err)
	}

	return uint32(safeGasPrice), uint32(proposeGasPrice), uint32(fastGasPrice)
}
