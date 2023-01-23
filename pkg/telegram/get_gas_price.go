package telegram

import (
	"encoding/json"
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"go_eth_bot/config"
	"go_eth_bot/internal/entity"
	"io"
	"log"
	"net/http"
	"strconv"
)

// GetGasPrice функция получения текущего газа сети eth
func GetEthGasRequest(cfg *config.Config) (uint32, uint32, uint32) {
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

func GetEthGas(ChatID int64, usersList map[int64]string, cfg *config.Config, bot *tgbotapi.BotAPI) {
	lowGas, averageGas, highGas := GetEthGasRequest(cfg)
	str := fmt.Sprintf("Low %d gwei \nAverage %d gwei \nHigh %d gwei", lowGas, averageGas, highGas)

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
