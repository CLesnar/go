package weather

import (
	"context"
)

const (
	version = "0.1.0"

	EnvVarOpenWeatherMapBaeUrl   = "OPEN_WEATHER_MAP_BASE_URL"
	DefaultOpenWeatherMapBaseUrl = "https://api.openweathermap.org/data/2.5/weather"
)

var (
	OpenWeatherMapBaseUrl = DefaultOpenWeatherMapBaseUrl
)

type Weather interface {
	GetWeatherCondition(ctx context.Context, lat float64, lon float64, apiKey string) (string, error)
}

func Version() string {
	return version
}

func temperatureCondition(temp float64) string {
	tempCondition := "normal"
	if temp >= 100 {
		tempCondition = "Extreme Heat"
	} else if temp >= 80 {
		tempCondition = "Hot"
	} else if temp >= 60 {
		tempCondition = "Nice"
	} else if temp >= 40 {
		tempCondition = "Cold"
	} else if temp > 32 {
		tempCondition = "Almost Freezing"
	} else if temp >= 0 {
		tempCondition = "Freezing"
	} else {
		tempCondition = "Extreme Freeze"
	}
	return tempCondition
}
