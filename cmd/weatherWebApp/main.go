package main

import (
	"context"
	"fmt"

	"github.com/CLesnar/go/pkg/server"
	"github.com/CLesnar/go/pkg/weather"
)

func Main() {
	ctx := context.Background()
	fmt.Printf("version: %v\n", weather.Version()) // TODO replace with struct logger
	s := WeatherWebAppScope{}
	errChan := make(chan error)
	server.Serve(ctx, &s, errChan)
}
