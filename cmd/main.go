package main

import (
	"log"
	"os"
	"strings"

	"myWeather/internal/api"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/joho/godotenv"
	"github.com/robfig/cron/v3"
)

func main() {
	// Загрузка переменных окружения
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Ошибка загрузки .env файла")
	}
	// Открываем файл .env и записываем значение TELEGRAM_TOKEN в переменную telegramToken
	telegramToken := os.Getenv("TELEGRAM_TOKEN")
	if telegramToken == "" {
		log.Fatal("Токен Telegram не найден в .env файле")
	}
	// Создаём бота
	bot, err := tgbotapi.NewBotAPI(telegramToken)
	if err != nil {
		log.Panic(err)
	}
	// Дебаг для тестов
	bot.Debug = true
	log.Printf("Авторизован как %s", bot.Self.UserName)

	// Карта для хранения городов пользователей
	userCities := make(map[int64]string)

	// Планировщик, в данном случае уведомления будут в 8 утра
	c := cron.New()
	c.AddFunc("0 8 * * *", func() {
		api.SendWeather(bot, userCities) // Передаем userCities
	})
	c.Start()

	// Обработчик команд
	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60
	updates := bot.GetUpdatesChan(u)
	// Запускаем цикл по сообщениям пользователей
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
			case "/weather":
				city, ok := userCities[chatID]
				if !ok {
					msg := tgbotapi.NewMessage(chatID, "Вы не указали город. Введите команду /city, чтобы установить ваш город.")
					bot.Send(msg)
					continue
				}
				weather, err := api.GetWeather(city)
				if err != nil {
					msg := tgbotapi.NewMessage(chatID, "Не удалось получить данные о погоде. Попробуйте позже.")
					bot.Send(msg)
					log.Println("Ошибка получения погоды:", err)
					continue
				}
				msg := tgbotapi.NewMessage(chatID, weather)
				bot.Send(msg)
			case "/help":
				msg := tgbotapi.NewMessage(chatID, "Доступные команды:\n/city - Установить город\n/weather - Получить текущую погоду")
				bot.Send(msg)
			default:
				// Проверяем, начинается ли сообщение с "/"
				if strings.HasPrefix(update.Message.Text, "/") {
					msg := tgbotapi.NewMessage(chatID, "Команда не распознана. Введите /help для списка доступных команд.")
					bot.Send(msg)
				} else {
					// Сохраняем введенный город
					userCities[chatID] = update.Message.Text
					msg := tgbotapi.NewMessage(chatID, "Город изменен на "+update.Message.Text)
					bot.Send(msg)
				}
			}
		}
	}
}
