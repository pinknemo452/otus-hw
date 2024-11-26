package hw05parallelexecution

import (
	"context"
	"errors"
	"sync"
)

var (
	ErrErrorsLimitExceeded = errors.New("errors limit exceeded")
	errJobResultsChClosed  = errors.New("job results channel closed")
)

type Task func() error

func doTask(ctx context.Context, jobResultsCh chan error, task Task) error {
	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
			err := task()
			if err != nil {
				jobResultsCh <- err
			}
			return nil
		}
	}
}

func jobResultCollector(ctxCancel context.CancelCauseFunc,
	jobResultsCh <-chan error, errorLimit int, wg *sync.WaitGroup,
) error {
	defer wg.Done()
	errCounter := 0
	done := false
	for err := range jobResultsCh {
		if err != nil {
			errCounter++

			if errCounter >= errorLimit && !done {
				ctxCancel(ErrErrorsLimitExceeded)
				done = true
			}
		}
	}

	return errJobResultsChClosed
}

// Run starts tasks in n goroutines and stops its work when receiving m errors from tasks.
func Run(tasks []Task, n, m int) error {
	wg := sync.WaitGroup{}

	jobResultsCh := make(chan error)
	ctx, cancel := context.WithCancelCause(context.Background())

	chunkLength := len(tasks) / n

	jobResultCollectorWg := sync.WaitGroup{}
	jobResultCollectorWg.Add(1)
	go jobResultCollector(cancel, jobResultsCh, m, &jobResultCollectorWg)

	for i := range n {
		wg.Add(1)
		chunkStart := i * chunkLength
		chunkEnd := chunkStart + chunkLength
		taskChunk := tasks[chunkStart:chunkEnd]
		go func() {
			defer wg.Done()
			for _, task := range taskChunk {
				err := doTask(ctx, jobResultsCh, task)
				if err != nil {
					return
				}
			}
		}()
	}

	wg.Wait()
	close(jobResultsCh)
	jobResultCollectorWg.Wait()
	return context.Cause(ctx)
}
