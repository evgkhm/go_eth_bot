package eth

import (
	"encoding/json"
	"fmt"
	"go_eth_bot/config"
	"go_eth_bot/internal/entity"
	"go_eth_bot/internal/service/util"
	"io"
	"log"
	"math/big"
	"net/http"
	"strconv"
)

// getBalanceRequest функция получения текущего баланса eth пользователя
func getBalanceRequest(cfg *config.Config, address string) *big.Float {
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
		log.Println(httpGetErr)
	}

	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			log.Println(err)
		}
	}(resp.Body)

	//парсинг данных, из запроса получаем WEI
	var cResp entity.CryptoUserData
	if decodeJsonErr := json.NewDecoder(resp.Body).Decode(&cResp); decodeJsonErr != nil {
		log.Println("ooopsss! an error occurred, please try again")
	}
	wei := new(big.Float)
	wei.SetString(cResp.Result)

	weiDivision := big.NewFloat(1000000000000000000)

	//из WEI в ETH
	ethBalance := new(big.Float).Quo(wei, weiDivision)

	return ethBalance
}

func GetBalance(chatID int64, usersList map[int64]string, cfg *config.Config) (string, string) {
	var newResp entity.CryptoUserData
	var IsExistAddr bool
	newResp.Address, IsExistAddr = util.GetAddFromMap(usersList, chatID)
	if IsExistAddr {
		ethBalance := getBalanceRequest(cfg, newResp.Address)
		ethBalanceWithoutName := fmt.Sprint(ethBalance)
		ethBalanceWithName := fmt.Sprint(ethBalance, " ETH")

		return ethBalanceWithName, ethBalanceWithoutName
	}
	return "", ""
}

func GetBalanceUSD(ethBalance string, cfg *config.Config) string {
	if ethBalance == "" {
		return "0 USD"
	}

	// из string во float64
	ethBalanceFloat, err := strconv.ParseFloat(ethBalance, 64)
	if err != nil {
		log.Println(err)
		return "0 USD"
	}

	ethPrice := getEthPriceRequest(cfg)

	usdBalance := new(big.Float).Mul(big.NewFloat(ethBalanceFloat), ethPrice)
	return fmt.Sprintf("%.0f USD", usdBalance)
}
