package main

import (
	"log"
	"os"

	"myWeather/internal/api" // Импортируем пакет api

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/joho/godotenv"
	"github.com/robfig/cron/v3"
)

func main() {
	// Загрузка переменных окружения из .env файла
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Ошибка загрузки .env файла")
	}

	// Получаем токен из переменной окружения
	telegramToken := os.Getenv("TELEGRAM_TOKEN")
	if telegramToken == "" {
		log.Fatal("Токен Telegram не найден в .env файле")
	}

	bot, err := tgbotapi.NewBotAPI(telegramToken)
	if err != nil {
		log.Panic(err)
	}

	bot.Debug = true
	log.Printf("Авторизован как %s", bot.Self.UserName)

	// Настроим планировщик для отправки прогноза погоды
	c := cron.New()
	c.AddFunc("38 19 * * *", func() { // каждый день в 8 утра
		api.SendWeather(bot) // вызываем функцию из пакета api
	})
	c.Start()

	// Создаем канал для получения обновлений
	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60
	updates := bot.GetUpdatesChan(u)

	// Слушаем команды
	for update := range updates {
		if update.Message != nil {
			chatID := update.Message.Chat.ID

			switch update.Message.Text {
			case "/start":
				msg := tgbotapi.NewMessage(chatID, "Привет! Я бот прогноза погоды. Каждый день в 8 утра я отправлю прогноз погоды. Вы можете ввести /help чтобы увидеть список команд.")
				bot.Send(msg)
			case "/city":
				msg := tgbotapi.NewMessage(chatID, "Введите название города:")
				bot.Send(msg)
			case "/help":
				msg := tgbotapi.NewMessage(chatID, "Доступные команды:\n/city - Установить город\n/weather - Получить текущую погоду")
				bot.Send(msg)
			default:
				msg := tgbotapi.NewMessage(chatID, "Привет! Я бот прогноза погоды. Каждый день в 8 утра я отправлю прогноз погоды. Вы можете ввести /help чтобы увидеть список команд.")
				bot.Send(msg)
			}
		}
	}
}
