package api

import (
	"log"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/robfig/cron/v3"
)

// Функция для отправки прогноза погоды каждому пользователю
func SendWeather(bot *tgbotapi.BotAPI) {
	for chatID, city := range userCities {
		weather, err := getWeather(city)
		if err != nil {
			log.Println("Ошибка при получении погоды для города", city, ":", err)
			continue
		}

		msg := tgbotapi.NewMessage(chatID, weather)
		bot.Send(msg)
	}
}

// Функция для настройки планировщика
func SetupCron(bot *tgbotapi.BotAPI) {
	c := cron.New()
	c.AddFunc("38 19 * * *", func() { // каждый день в 8 утра
		SendWeather(bot) // вызов функции отправки прогноза погоды
	})
	c.Start()
}
