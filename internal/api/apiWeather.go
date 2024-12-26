package api

import (
	"encoding/json"
	"fmt"
	"net/http"
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

// Словарь для хранения выбранного города для каждого пользователя
var userCities = make(map[int64]string)

// Функция для получения прогноза погоды
func getWeather(city string) (string, error) {
	weatherAPIKey := "your_api_key_here" // Убедитесь, что у вас есть правильный API ключ
	url := fmt.Sprintf("https://api.openweathermap.org/data/2.5/weather?q=%s&units=metric&appid=%s", city, weatherAPIKey)
	resp, err := http.Get(url)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	var weather WeatherResponse
	if err := json.NewDecoder(resp.Body).Decode(&weather); err != nil {
		return "", err
	}

	return fmt.Sprintf("Погода в %s: %.1f°C, %s", city, weather.Main.Temp, weather.Weather[0].Description), nil
}
