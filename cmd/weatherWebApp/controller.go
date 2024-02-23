package main

import (
	"fmt"
	"net/http"

	"github.com/gorilla/mux"

	"github.com/CLesnar/go/pkg/weather"
)

type WeatherWebAppController struct {
}

var (
	controller = &WeatherWebAppController{}
)

func (c *WeatherWebAppController) Router(r *mux.Router) {
	s := r.PathPrefix("v1/weather").Subrouter()
	s.HandleFunc("/data", controller.GetWeatherDataHandler)
}

func (c *WeatherWebAppController) GetWeatherDataHandler(w http.ResponseWriter, r *http.Request) {

	openWeatherMap := weather.OpenWeatherMap{}
	resp, err := openWeatherMap.GetCurrentWeatherData(ctx, criteria, parameters)
	if err != nil {
		fmt.Printf("error: %v", err)
	}
	data, ok := resp.(OpenWeatherMapGetCurrentDataResponse)
}
