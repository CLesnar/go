package server

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	environment "github.com/CLesnar/go/internal/pkg/enviornment"
	"github.com/CLesnar/go/pkg/middleware"
	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
)

const (
	DefaultPort = 8080
)

type Scope interface {
	Route(*mux.Router)
}

func Serve(ctx context.Context, s Scope, errChan chan<- error) {
	r := NewRouter()
	s.Route(r)

	// CORS
	headersAllowed := handlers.AllowedHeaders([]string{"Content-Type", "Authorization"}) // Add more headers such as ones defined in datadog logs or similar
	methodsAllowed := handlers.AllowedMethods([]string{http.MethodDelete, http.MethodGet, http.MethodPost, http.MethodPut})
	originsAllowed := handlers.AllowedOrigins([]string{"*"}) // TODO: allows select domains only for a real server - maybe use AllowedOriginsValidator() ?
	cors := handlers.CORS(headersAllowed, methodsAllowed, originsAllowed)
	corshandler := cors(r)

	timeout := environment.GetEnvVar("TIMEOUT", "30")
	port := environment.GetEnvVar("PORT", "8111")
	if devenv := environment.GetEnvVar("DEVENV", ""); len(devenv) > 0 {
		timeout = "500" // increase timeout for debugging
	}
	timeoutDuration, err := time.ParseDuration(timeout + "s")
	if err != nil {
		fmt.Printf("error parsing timeout duration: %v", err) // TODO replace struct logger
	}
	httpServe := &http.Server{
		Handler:      corshandler,
		Addr:         port,
		WriteTimeout: timeoutDuration,
		ReadTimeout:  timeoutDuration,
	}
	go func() {
		errChan <- httpServe.ListenAndServe()
	}()
}

func NewRouter() *mux.Router {
	r := mux.NewRouter().StrictSlash(true)
	r.HandleFunc("/health", HealthCheckHandler)
	r.HandleFunc("/docs/openapi.json", OpenAPIDocHandler)
	r.Use(middleware.AuthenticationMiddleware,
		middleware.PermissionsMiddleware,
		middleware.LogRouteMiddleware,
		middleware.LogRouteMiddleware,
		middleware.AddCtxMapMiddleware)
	return r
}

func OpenAPIDocHandler(w http.ResponseWriter, r *http.Request) {
	openapijson := map[string]interface{}{
		"openapi": "3.0.2",
		"info": map[string]interface{}{
			"title":       "'TODO' Service",
			"description": "TODO: get swagger openapi.json",
		},
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(openapijson) // nolint
}

func HealthCheckHandler(w http.ResponseWriter, r *http.Request) {
	healthCheck := map[string]bool{"ok": true}
	// TODO: check health of process / go routines
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(healthCheck) // nolint
}
