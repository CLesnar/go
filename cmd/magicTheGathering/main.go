package main

import (
	"context"
	"fmt"
	"time"

	"github.com/CLesnar/go/pkg/magicTheGathering"
	"github.com/CLesnar/go/pkg/server"
)

func main() {
	fmt.Printf("starting magicTheGathering service...\n")
	ctx := context.WithValue(context.Background(), server.ServeCtxKeyEnvPort, "MTG_PORT")
	fmt.Printf("version: %v\n", magicTheGathering.Version()) // TODO replace with struct logger
	s := MtgScope{}
	errChan := make(chan error)
	server.Serve(ctx, &s, errChan)
	sched := server.Scheduler{}
	sched.Add(ctx, "game-manager", time.Duration(1*time.Second), server.TaskFunc(GameManager))
	sched.Start(ctx)
	fmt.Printf("exiting weatherWebApp service.\n")
}
