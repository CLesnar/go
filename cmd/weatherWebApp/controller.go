package main

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/gorilla/mux"

	chttp "github.com/CLesnar/go/internal/pkg/http"
	"github.com/CLesnar/go/internal/pkg/model"
	"github.com/CLesnar/go/pkg/weather"
)

type WeatherWebAppScope struct {
}

var (
	controller = &WeatherWebAppScope{}
)

func (c *WeatherWebAppScope) Route(r *mux.Router) {
	s := r.PathPrefix("/v1/weather").Subrouter()
	s.HandleFunc("/data", controller.GetWeatherDataHandler).Methods(http.MethodGet)
	// Add middleware specific to this controller
}

func (c *WeatherWebAppScope) GetWeatherDataHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	openWeatherMap, params := weather.OpenWeatherMap{}, model.WeatherConditionGetParameters{}
	if err := chttp.DecodeQueryString(r, &params); err != nil {
		chttp.RespondError(w, http.StatusBadRequest, err)
		return
	}
	if err := params.Validate(); err != nil {
		chttp.RespondError(w, http.StatusInternalServerError, err)
		return
	}
	resp, err := openWeatherMap.GetWeatherCondition(ctx, *params.Latitude, *params.Longitude, *params.ApiId)
	if err != nil {
		chttp.RespondError(w, http.StatusInternalServerError, err)
		return
	}
	data, err := json.Marshal(resp)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(fmt.Sprintf("error: %v", err))) // nolint
		return
	}
	w.Header().Set("Content-Type", "application/json")
	if _, err := w.Write(data); err != nil { // failed to write response
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
}
