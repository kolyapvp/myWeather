package api

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/joho/godotenv"
)

// Структура для хранения ответа от OpenWeatherMap
type WeatherResponse struct {
	Main struct {
		Temp float64 `json:"temp"`
	} `json:"main"`
	Weather []struct {
		Description string `json:"description"`
	} `json:"weather"`
}

// Функция для получения прогноза погоды
func GetWeather(city string) (string, error) {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Ошибка загрузки .env файла")
	}
	// Берём и вставляем Апи ключ для сервиса погоды
	weatherAPIKey := os.Getenv("WeatherAPIKey")
	if weatherAPIKey == "" {
		log.Fatal("Токен WeatherAPIKey не найден в .env файле")
	}
	// Формируем строку для http запроса
	url := fmt.Sprintf("https://api.openweathermap.org/data/2.5/weather?q=%s&units=metric&appid=%s", city, weatherAPIKey)
	resp, err := http.Get(url)
	if err != nil {
		return "", fmt.Errorf("ошибка запроса к API: %w", err)
	}
	defer resp.Body.Close()

	var weather WeatherResponse
	if err := json.NewDecoder(resp.Body).Decode(&weather); err != nil {
		return "", fmt.Errorf("ошибка декодирования ответа API: %w", err)
	}

	// Проверка, что массив weather.Weather не пуст
	if len(weather.Weather) == 0 {
		return "", fmt.Errorf("нет данных о погоде для города %s", city)
	}

	return fmt.Sprintf("Погода в %s: %.1f°C, %s", city, weather.Main.Temp, weather.Weather[0].Description), nil
}
