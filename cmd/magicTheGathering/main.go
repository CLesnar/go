package main

import (
	"context"
	"fmt"

	"github.com/CLesnar/go/pkg/magicTheGathering"
	"github.com/CLesnar/go/pkg/server"
)

func main() {
	fmt.Printf("starting magicTheGathering service...\n")
	ctx := context.WithValue(context.Background(), server.ServeCtxKeyEnvPort, "MTG_PORT")
	fmt.Printf("version: %v\n", magicTheGathering.Version()) // TODO replace with struct logger
	s := WeatherWebAppScope{}
	errChan := make(chan error)
	server.Serve(ctx, &s, errChan)
	fmt.Printf("exiting weatherWebApp service.\n")
}
