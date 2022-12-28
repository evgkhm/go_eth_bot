package controller

import (
	"fmt"
	"go_eth_bot/internal/usecase"
	"go_eth_bot/pkg/telegram"
	"math/big"
)

type Handler struct {
	services *usecase.Service
}

func (h Handler) Run(updates *usecase.Service) {
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

func New(service *usecase.Service) *Handler {
	return &Handler{services: service}
}

//func ReadUpdates(updates *telegram.Updates) {
//	// читаем обновления из канала
//	for update := range updates {
//		if update.Message != nil && update.Message.Text == "/start" {
//			//в чат вошел новый пользователь. Поприветствуем его
//			str := fmt.Sprintf(`Привет %s! Этот бот показывает стоимость эфира, газа и текущий баланс.
//Для проверки баланса введите адрес кошелька`, update.Message.From.FirstName)
//			SendTgMess(update.Message.Chat.ID, str, bot, telegram.First)
//		} else if update.Message != nil {
//			//если получили обычное сообщение сообщение от пользователя в ТГ
//			newResp.Address = update.Message.Text
//			if IsValidAddress(newResp.Address) {
//				ChatID := update.Message.Chat.ID    //получаем ID пользователя
//				usersList[ChatID] = newResp.Address //проверить что уникальный ID добавляется 1 раз!!!!!!!!
//				str := "Адрес получен. Выберете действие"
//				SendTgMess(update.Message.Chat.ID, str, bot, telegram.Second)
//			} else {
//				newResp.Address = ""
//				str := "Введите ETH адрес"
//				SendTgMess(update.Message.Chat.ID, str, bot, telegram.First)
//			}
//		}
//
//		//если получили нажатие кнопки
//		if update.CallbackQuery != nil {
//			switch update.CallbackQuery.Data {
//			case "/get_balance":
//				if IsValidAddress(newResp.Address) { //проверка на валидность адреса
//					ChatID := update.CallbackQuery.Message.Chat.ID //получаем ID пользователя
//					newResp.Address = usersList[ChatID]            //извлечение из мапы адрес эфира
//					ethBalance := GetBalanceRequest(newResp.Address)
//					str := fmt.Sprint(ethBalance, " ETH")
//					SendTgMess(update.CallbackQuery.Message.Chat.ID, str, bot, telegram.Second)
//				} else {
//					newResp.Address = ""
//					str := "Некорректный адрес"
//					SendTgMess(update.CallbackQuery.Message.Chat.ID, str, bot, telegram.First)
//				}
//			case "/get_balance_usd":
//				if IsValidAddress(newResp.Address) { //проверка на валидность
//					ChatID := update.CallbackQuery.Message.Chat.ID //получаем ID пользователя
//					newResp.Address = usersList[ChatID]            //извлечение из мапы адрес эфира
//					ethBalance := GetBalanceRequest(newResp.Address)
//					ethPrice := GetEthPrice()
//					usdBalance := new(big.Float).Mul(ethBalance, ethPrice)
//					str := fmt.Sprintf("%.2f USD", usdBalance)
//					SendTgMess(update.CallbackQuery.Message.Chat.ID, str, bot, telegram.Second)
//				} else {
//					newResp.Address = ""
//					str := "Некорректный адрес"
//					SendTgMess(update.CallbackQuery.Message.Chat.ID, str, bot, telegram.First)
//				}
//			case "/get_price":
//				ethPrice := GetEthPrice()
//				str := fmt.Sprint(ethPrice, " USD")
//
//				//Если адреса нет вызов первой клавиатуры
//				if newResp.Address == "" {
//					SendTgMess(update.CallbackQuery.Message.Chat.ID, str, bot, telegram.First)
//				} else {
//					SendTgMess(update.CallbackQuery.Message.Chat.ID, str, bot, telegram.Second)
//				}
//			case "/get_gas":
//				lowGas, averageGas, highGas := GetGasPrice()
//				str := fmt.Sprintf("Low %d gwei \nAverage %d gwei \nHigh %d gwei", lowGas, averageGas, highGas)
//
//				//Если адреса нет вызов первой клавиатуры
//				if newResp.Address == "" {
//					SendTgMess(update.CallbackQuery.Message.Chat.ID, str, bot, telegram.First)
//				} else {
//					SendTgMess(update.CallbackQuery.Message.Chat.ID, str, bot, telegram.Second)
//				}
//			case "/change_addr":
//				ChatID := update.CallbackQuery.Message.Chat.ID //получаем ID пользователя
//				delete(usersList, ChatID)                      //удаляем из мапы пользователя
//				newResp.Address = ""
//				fallthrough
//			default:
//				newResp.Address = update.CallbackQuery.Message.Text
//				if IsValidAddress(newResp.Address) {
//					str := "Адрес получен. Выберете действие"
//					SendTgMess(update.CallbackQuery.Message.Chat.ID, str, bot, telegram.Second)
//				} else {
//					newResp.Address = ""
//					str := "Введите ETH адрес"
//					SendTgMess(update.CallbackQuery.Message.Chat.ID, str, bot, telegram.First)
//				}
//			}
//		}
//	}
//}
