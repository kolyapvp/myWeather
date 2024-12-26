package api

import (
	"log"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

// Функция для отправки прогноза погоды каждому пользователю
func SendWeather(bot *tgbotapi.BotAPI, userCities map[int64]string) {
	for chatID, city := range userCities {
		// Используем GetWeather из api
		weather, err := GetWeather(city)
		if err != nil {
			log.Println("Ошибка при получении погоды для города", city, ":", err)
			continue
		}

		msg := tgbotapi.NewMessage(chatID, weather)
		bot.Send(msg)
	}
}
