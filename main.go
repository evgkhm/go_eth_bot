package main

import (
	"encoding/json"
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/joho/godotenv"
	"go_eth_bot/internal/entity"
	"go_eth_bot/pkg/telegram"
	"io"
	"log"
	"math/big"
	"net/http"
	"os"
	"regexp"
	"strconv"
)

func main() {
	dotenv := goDotEnvVariable("TG_API_KEY")
	// подключаемся к телеграм боту с помощью токена
	bot, err := tgbotapi.NewBotAPI(dotenv)
	if err != nil {
		log.Panic(err)
	}

	bot.Debug = false

	updates := bot.ListenForWebhook("/" + bot.Token)

	//создание сервера, чтобы heroku не ругался на port
	http.HandleFunc("/", MainHandler) //закинуть в pkg->httpserver-server.go
	go func() {
		err := http.ListenAndServe(":"+goDotEnvVariable("PORT"), nil)
		if err != nil {
			log.Panic(err)
		}
	}()

	var newResp entity.CryptoUserData
	usersList := make(map[int64]string) //здесь список всех пользователей

	// читаем обновления из канала
	for update := range updates {
		if update.Message != nil && update.Message.Text == "/start" {
			//в чат вошел новый пользователь. Поприветствуем его
			str := fmt.Sprintf(`Привет %s! Этот бот показывает стоимость эфира, газа и текущий баланс.
Для проверки баланса введите адрес кошелька`, update.Message.From.FirstName)
			SendTgMess(update.Message.Chat.ID, str, bot, telegram.First)
		} else if update.Message != nil {
			//если получили обычное сообщение сообщение от пользователя в ТГ
			newResp.Address = update.Message.Text
			if IsValidAddress(newResp.Address) {
				ChatID := update.Message.Chat.ID    //получаем ID пользователя
				usersList[ChatID] = newResp.Address //проверить что уникальный ID добавляется 1 раз!!!!!!!!
				str := "Адрес получен. Выберете действие"
				SendTgMess(update.Message.Chat.ID, str, bot, telegram.Second)
			} else {
				newResp.Address = ""
				str := "Введите ETH адрес"
				SendTgMess(update.Message.Chat.ID, str, bot, telegram.First)
			}
		}

		//если получили нажатие кнопки
		if update.CallbackQuery != nil {
			switch update.CallbackQuery.Data {
			case "/get_balance":
				if IsValidAddress(newResp.Address) { //проверка на валидность адреса
					ChatID := update.CallbackQuery.Message.Chat.ID //получаем ID пользователя
					newResp.Address = usersList[ChatID]            //извлечение из мапы адрес эфира
					ethBalance := GetBalanceRequest(newResp.Address)
					str := fmt.Sprint(ethBalance, " ETH")
					SendTgMess(update.CallbackQuery.Message.Chat.ID, str, bot, telegram.Second)
				} else {
					newResp.Address = ""
					str := "Некорректный адрес"
					SendTgMess(update.CallbackQuery.Message.Chat.ID, str, bot, telegram.First)
				}
			case "/get_balance_usd":
				if IsValidAddress(newResp.Address) { //проверка на валидность
					ChatID := update.CallbackQuery.Message.Chat.ID //получаем ID пользователя
					newResp.Address = usersList[ChatID]            //извлечение из мапы адрес эфира
					ethBalance := GetBalanceRequest(newResp.Address)
					ethPrice := GetEthPrice()
					usdBalance := new(big.Float).Mul(ethBalance, ethPrice)
					str := fmt.Sprintf("%.2f USD", usdBalance)
					SendTgMess(update.CallbackQuery.Message.Chat.ID, str, bot, telegram.Second)
				} else {
					newResp.Address = ""
					str := "Некорректный адрес"
					SendTgMess(update.CallbackQuery.Message.Chat.ID, str, bot, telegram.First)
				}
			case "/get_price":
				ethPrice := GetEthPrice()
				str := fmt.Sprint(ethPrice, " USD")

				//Если адреса нет вызов первой клавиатуры
				if newResp.Address == "" {
					SendTgMess(update.CallbackQuery.Message.Chat.ID, str, bot, telegram.First)
				} else {
					SendTgMess(update.CallbackQuery.Message.Chat.ID, str, bot, telegram.Second)
				}
			case "/get_gas":
				lowGas, averageGas, highGas := GetGasPrice()
				str := fmt.Sprintf("Low %d gwei \nAverage %d gwei \nHigh %d gwei", lowGas, averageGas, highGas)

				//Если адреса нет вызов первой клавиатуры
				if newResp.Address == "" {
					SendTgMess(update.CallbackQuery.Message.Chat.ID, str, bot, telegram.First)
				} else {
					SendTgMess(update.CallbackQuery.Message.Chat.ID, str, bot, telegram.Second)
				}
			case "/change_addr":
				ChatID := update.CallbackQuery.Message.Chat.ID //получаем ID пользователя
				delete(usersList, ChatID)                      //удаляем из мапы пользователя
				newResp.Address = ""
				fallthrough
			default:
				newResp.Address = update.CallbackQuery.Message.Text
				if IsValidAddress(newResp.Address) {
					str := "Адрес получен. Выберете действие"
					SendTgMess(update.CallbackQuery.Message.Chat.ID, str, bot, telegram.Second)
				} else {
					newResp.Address = ""
					str := "Введите ETH адрес"
					SendTgMess(update.CallbackQuery.Message.Chat.ID, str, bot, telegram.First)
				}
			}
		}
	}
}

// MainHandler функция приветствия для правильной работы с heroku
func MainHandler(resp http.ResponseWriter, _ *http.Request) {
	_, err := resp.Write([]byte("Hi there! I'm Bot!"))
	if err != nil {
		log.Panic(err)
	}
}

// use godot package to load/read the .env file and
// return the value of the key
func goDotEnvVariable(key string) string {
	// load .env file
	err := godotenv.Load(".env")

	if err != nil {
		log.Fatalf("Error loading .env file")
	}

	return os.Getenv(key)
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
