package utils

import (
	"context"
	"fmt"
	"sync"
	"time"
)

type WorkerTask[T any] interface {
	GetTasks(ctx context.Context) ([]T, error)
	Process(ctx context.Context, task T) error
}

type Worker[T any] struct {
	tasks      WorkerTask[T]
	numWorkers int
}

func NewWorker[T any](tasks WorkerTask[T], numWorkers int) *Worker[T] {
	return &Worker[T]{tasks: tasks, numWorkers: numWorkers}
}

func (w *Worker[T]) Start(ctx context.Context, interval time.Duration) {
	ticker := time.NewTicker(interval)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			w.RunBatch(ctx)
		}
	}
}

func (w *Worker[T]) RunBatch(ctx context.Context) error {
	tasks, err := w.tasks.GetTasks(ctx)

	if err != nil {
		return fmt.Errorf("failed to get items: %w", err)
	}

	if len(tasks) == 0 {
		return nil
	}

	sem := make(chan struct{}, w.numWorkers)
	var wg sync.WaitGroup

	var firstErr error
	var errOnce sync.Once

	for _, task := range tasks {
		sem <- struct{}{}
		wg.Add(1)

		go func() {
			defer func() {
				<-sem
				wg.Done()
			}()

			err := w.tasks.Process(ctx, task)
			if err != nil {
				errOnce.Do(func() { firstErr = err })
			}
		}()
	}

	wg.Wait()

	return firstErr
}
