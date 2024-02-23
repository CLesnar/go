package main

import (
	"context"
	"fmt"

	"github.com/CLesnar/go/pkg/server"
	"github.com/CLesnar/go/pkg/weather"
)

func main() {
	fmt.Printf("starting weatherWebApp service...\n")
	ctx := context.Background()
	fmt.Printf("version: %v\n", weather.Version()) // TODO replace with struct logger
	s := WeatherWebAppScope{}
	errChan := make(chan error)
	server.Serve(ctx, &s, errChan)
	fmt.Printf("exiting weatherWebApp service.\n")
}
