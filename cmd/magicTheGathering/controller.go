package main

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/gorilla/mux"

	chttp "github.com/CLesnar/go/internal/pkg/http"
	"github.com/CLesnar/go/pkg/model"
	"github.com/CLesnar/go/pkg/weather"
)

type MtgScope struct {
}

func (c *MtgScope) Route(r *mux.Router) {
	s := r.PathPrefix("/v1/magic").Subrouter()

	// Setup
	s.HandleFunc("/sets", c.HandleGetSets).Methods(http.MethodGet)
	s.HandleFunc("/sets/{sid}", c.HandleGetSets).Methods(http.MethodGet)
	s.HandleFunc("/cards", c.HandleGetCards).Methods(http.MethodGet)
	s.HandleFunc("/cards/{cid}", c.HandleGetCards).Methods(http.MethodGet)
	s.HandleFunc("/pools", c.HandleGetPoolCards).Methods(http.MethodGet)
	s.HandleFunc("/pools/{ppid}", c.HandleGetPoolCards).Methods(http.MethodGet)
	s.HandleFunc("/decks/{pid}", c.HandlePostCardToPlayerDeck).Methods(http.MethodPost)
	s.HandleFunc("/decks", c.HandleGetPlayerDeck).Methods(http.MethodGet)
	s.HandleFunc("/decks/{pid}", c.HandleGetPlayerDeck).Methods(http.MethodGet)

	// Game State
	s.HandleFunc("/games", c.HandleGetGameState).Methods(http.MethodGet)

	// Player Interactions
	s.HandleFunc("/games/{gid}/players/{pid}", c.HandlePlayerAction).Methods(http.MethodGet)

	// Add middleware specific to this controller
}

func (c *MtgScope) HandleGetSets(w http.ResponseWriter, r *http.Request) {
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
	resp, err := openWeatherMap.GetWeatherCondition(ctx, *params.Latitude, *params.Longitude, *params.AppId)
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
