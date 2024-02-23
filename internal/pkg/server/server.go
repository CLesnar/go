package server

import (
	"context"

	"github.com/gorilla/mux"
)

const (
	DefaultPort = 8080
)

func Serve(ctx context.Context, scopes []ScopeRouter, errChan chan<- error) {
	r := mux.NewRouter()
	// r.Use(various built-in middleware)
}
