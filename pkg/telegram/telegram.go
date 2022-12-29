package telegram

import (
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"go_eth_bot/config"
	"go_eth_bot/internal/entity"
	"log"
	"math/big"
)

// Page для выбора клавиатуры ТГ
type Page int

const (
	First Page = iota + 1
	Second
)

type Updates struct {
	updates tgbotapi.UpdatesChannel
	bot     *tgbotapi.BotAPI
}

func (u Updates) Run(cfg *config.Config) {
	usersList := make(map[int64]string) //здесь список всех пользователей

	for update := range u.updates {
		if update.Message != nil && update.Message.Text == "/start" {
			//в чат вошел новый пользователь. Поприветствуем его
			str := fmt.Sprintf(`Привет %s! Этот бот показывает стоимость эфира, газа и текущий баланс.
Для проверки баланса введите адрес кошелька`, update.Message.From.FirstName)
			SendTgMess(update.Message.Chat.ID, str, u.bot, First)
		} else if update.Message != nil {
			//если получили обычное сообщение сообщение от пользователя в ТГ
			var newResp entity.CryptoUserData
			newResp.Address = update.Message.Text
			if IsValidAddress(newResp.Address) {
				ChatID := update.Message.Chat.ID    //получаем ID пользователя
				usersList[ChatID] = newResp.Address //проверить что уникальный ID добавляется 1 раз!!!!!!!!
				str := "Адрес получен. Выберете действие"
				SendTgMess(update.Message.Chat.ID, str, u.bot, Second)
			} else {
				newResp.Address = ""
				str := "Введите ETH адрес"
				SendTgMess(update.Message.Chat.ID, str, u.bot, First)
			}
		}

		//если получили нажатие кнопки
		if update.CallbackQuery != nil {
			switch update.CallbackQuery.Data {
			case "/get_balance":
				ChatID := update.CallbackQuery.Message.Chat.ID //получаем ID пользователя
				var newResp entity.CryptoUserData
				var IsExistAddr bool
				newResp.Address, IsExistAddr = GetAddFromMap(usersList, ChatID)
				if IsExistAddr {
					ethBalance := GetBalanceRequest(cfg, newResp.Address)
					str := fmt.Sprint(ethBalance, " ETH")
					SendTgMess(update.CallbackQuery.Message.Chat.ID, str, u.bot, Second)
				} else {
					str := "Некорректный адрес"
					SendTgMess(update.CallbackQuery.Message.Chat.ID, str, u.bot, First)
				}

				//_, ok := usersList[ChatID]
				//if ok {
				//	var newResp entity.CryptoUserData
				//	newResp.Address = usersList[ChatID] //извлечение из мапы адрес эфира
				//	ethBalance := GetBalanceRequest(cfg, newResp.Address)
				//	str := fmt.Sprint(ethBalance, " ETH")
				//	SendTgMess(update.CallbackQuery.Message.Chat.ID, str, u.bot, Second)
				//} else {
				//	str := "Некорректный адрес"
				//	SendTgMess(update.CallbackQuery.Message.Chat.ID, str, u.bot, First)
				//}

			case "/get_balance_usd":
				ChatID := update.CallbackQuery.Message.Chat.ID //получаем ID пользователя
				_, ok := usersList[ChatID]
				if ok {
					var newResp entity.CryptoUserData
					newResp.Address = usersList[ChatID] //извлечение из мапы адрес эфира
					ethBalance := GetBalanceRequest(cfg, newResp.Address)
					ethPrice := GetEthPrice(cfg)
					usdBalance := new(big.Float).Mul(ethBalance, ethPrice)
					str := fmt.Sprintf("%.2f USD", usdBalance)
					SendTgMess(update.CallbackQuery.Message.Chat.ID, str, u.bot, Second)
				} else {
					str := "Некорректный адрес"
					SendTgMess(update.CallbackQuery.Message.Chat.ID, str, u.bot, First)
				}

			case "/get_price":
				ethPrice := GetEthPrice(cfg)
				str := fmt.Sprint(ethPrice, " USD")

				var newResp entity.CryptoUserData
				//Если адреса нет вызов первой клавиатуры
				if newResp.Address == "" {
					SendTgMess(update.CallbackQuery.Message.Chat.ID, str, u.bot, First)
				} else {
					SendTgMess(update.CallbackQuery.Message.Chat.ID, str, u.bot, Second)
				}
			case "/get_gas":
				lowGas, averageGas, highGas := GetGasPrice(cfg)
				str := fmt.Sprintf("Low %d gwei \nAverage %d gwei \nHigh %d gwei", lowGas, averageGas, highGas)

				var newResp entity.CryptoUserData
				//Если адреса нет вызов первой клавиатуры
				if newResp.Address == "" {
					SendTgMess(update.CallbackQuery.Message.Chat.ID, str, u.bot, First)
				} else {
					SendTgMess(update.CallbackQuery.Message.Chat.ID, str, u.bot, Second)
				}
			case "/change_addr":
				ChatID := update.CallbackQuery.Message.Chat.ID //получаем ID пользователя
				delete(usersList, ChatID)                      //удаляем из мапы пользователя
				var newResp entity.CryptoUserData              //?????????????????????????????????
				newResp.Address = ""
				fallthrough
			default:
				var newResp entity.CryptoUserData
				newResp.Address = update.CallbackQuery.Message.Text
				if IsValidAddress(newResp.Address) {
					str := "Адрес получен. Выберете действие"
					SendTgMess(update.CallbackQuery.Message.Chat.ID, str, u.bot, Second)
				} else {
					newResp.Address = ""
					str := "Введите ETH адрес"
					SendTgMess(update.CallbackQuery.Message.Chat.ID, str, u.bot, First)
				}
			}
		}
	}
}

func New(cfg *config.Config) *Updates {
	u := &Updates{}
	// подключаемся к телеграм боту с помощью токена
	bot, err := tgbotapi.NewBotAPI(cfg.TgApiKey)
	if err != nil {
		log.Panic(err)
	}

	bot.Debug = false

	u.bot = bot
	u.updates = bot.ListenForWebhook("/" + bot.Token)

	return u
}
