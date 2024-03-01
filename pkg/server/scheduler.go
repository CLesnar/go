package server

import (
	"context"
	"sync"
	"time"

	"github.com/CLesnar/go/internal/pkg/util"
)

type TaskFunc func(context.Context) bool

type Task struct {
	Name        string
	Periodicity time.Duration
	Func        TaskFunc
	Ctx         context.Context
}

type Scheduler struct {
	WG    sync.WaitGroup
	Tasks []Task
}

func (t *Task) makeScheduledTaskFunc(f TaskFunc) TaskFunc {
	return func(ctx context.Context) bool {
		defer util.Recover(ctx, t.Name)
		for {
			next := time.Now().UnixNano() + t.Periodicity.Nanoseconds()
			nextTime := time.Unix(0, next)
			if exit := t.Func(t.Ctx); exit {
				return exit
			}
			time.Sleep(time.Until(nextTime))
		}
	}
}

func (s *Scheduler) Add(ctx context.Context, name string, periodicity time.Duration, Func TaskFunc) {
	t := Task{Ctx: ctx, Name: name, Periodicity: periodicity}
	f := TaskFunc(func(ctx context.Context) bool {
		defer s.WG.Done()
		return Func(ctx)
	})
	t.Func = t.makeScheduledTaskFunc(f)
	s.WG.Add(1)
	s.Tasks = append(s.Tasks, t)
}

func (t *Task) Start() {
	go t.Func(t.Ctx)
}

func (s *Scheduler) Start(ctx context.Context) {
	for _, task := range s.Tasks {
		task.Start()
	}
	s.WG.Wait()
}
