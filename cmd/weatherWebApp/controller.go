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
	ctx := r.Context()
	openWeatherMap := weather.OpenWeatherMap{}

	resp, err := openWeatherMap.GetCurrentWeatherData(ctx, criteria, parameters)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = w.Write([]byte(fmt.Sprintf("error: %v", err)))
		return
	}
	data, ok := resp.(OpenWeatherMaResponsepGetCurrentData)
}
