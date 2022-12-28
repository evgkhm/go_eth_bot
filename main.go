package main

import (
	"encoding/json"
	"errors"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"go_eth_bot/config"
	"go_eth_bot/internal/controller"
	"go_eth_bot/internal/entity"
	"go_eth_bot/internal/usecase"
	"go_eth_bot/pkg/server"
	"go_eth_bot/pkg/telegram"
	"io"
	"log"
	"math/big"
	"net/http"
	"regexp"
	"strconv"
)

func main() {
	cfg, errConfig := config.NewConfig()
	if errConfig != nil {
		errors.New("can't get config")
	}

	serverErr := server.New(cfg.Port)
	if serverErr != nil {
		errors.New("can't create server")
	}

	updates := telegram.New(cfg.TgApiKey)

	service := usecase.New(updates)

	//var newResp entity.CryptoUserData
	//usersList := make(map[int64]string) //здесь список всех пользователей

	controller.ReadUpdates(updates)
}

// GetGasPrice функция получения текущего газа сети eth
func GetGasPrice() (uint32, uint32, uint32) {
	// godotenv package
	dotenv := goDotEnvVariable("API_KEY")

	resp, err := http.Get("https://api.etherscan.io/api" +
		"?module=gastracker" +
		"&action=gasoracle" +
		"&apikey=" + dotenv)
	if err != nil {
		log.Fatalln(err)
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

// SendTgMess функция отправки сообщения в ТГ
func SendTgMess(id int64, str string, bot *tgbotapi.BotAPI, page telegram.Page) {
	msg := tgbotapi.NewMessage(id, str)

	//Определение клавиатуры
	switch page {
	case telegram.First:
		msg.ReplyMarkup = telegram.FirstKeyboard
	case telegram.Second:
		msg.ReplyMarkup = telegram.SecondKeyboard
	}

	//Отправка сообщения в ТГ
	_, err := bot.Send(msg)
	if err != nil {
		return
	}
}

// IsValidAddress функция проверки валидности eth адреса
func IsValidAddress(v string) bool {
	re := regexp.MustCompile("^0x[\\da-fA-F]{40}$")
	return re.MatchString(v)
}

// GetEthPrice функция получения текущего курса eth
func GetEthPrice() *big.Float {
	// godotenv package
	dotenv := goDotEnvVariable("API_KEY")

	resp, err := http.Get("https://api.etherscan.io/api" +
		"?module=stats" +
		"&action=ethprice" +
		"&apikey=" + dotenv)
	if err != nil {
		log.Fatalln(err)
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

// GetBalanceRequest функция получения текущего баланса eth пользователя
func GetBalanceRequest(address string) *big.Float {
	// godotenv package
	dotenv := goDotEnvVariable("API_KEY")

	rawAddress := "&address=" + address

	resp, err := http.Get("https://api.etherscan.io/api" +
		"?module=account" +
		"&action=balance" +
		rawAddress +
		"&tag=latest" +
		"&apikey=" + dotenv)
	if err != nil {
		log.Fatalln(err)
	}

	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			log.Fatalln(err)
		}
	}(resp.Body)

	//парсинг данных, из запроса получаем WEI
	var cResp entity.CryptoUserData
	if err := json.NewDecoder(resp.Body).Decode(&cResp); err != nil {
		log.Fatal("ooopsss! an error occurred, please try again")
	}
	wei := new(big.Float)
	wei.SetString(cResp.Result)

	weiDivision := big.NewFloat(1000000000000000000)

	//из WEI в ETH
	ethBalance := new(big.Float).Quo(wei, weiDivision)

	return ethBalance
}
