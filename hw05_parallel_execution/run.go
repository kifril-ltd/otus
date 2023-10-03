package hw05parallelexecution

import (
	"errors"
	"sync"
	"sync/atomic"
)

var (
	ErrInvalidWorkersCount = errors.New("invalid number of workers")
	ErrErrorsLimitExceeded = errors.New("errors limit exceeded")
)

type Task func() error

// Run starts tasks in n goroutines and stops its work when receiving m errors from tasks.
func Run(tasks []Task, workersCount, errorsLimit int) error {
	if workersCount <= 0 {
		return ErrInvalidWorkersCount
	}

	tasksChan := make(chan Task, len(tasks))
	wg := sync.WaitGroup{}

	var errCounter int32
	for i := 0; i < workersCount; i++ {
		wg.Add(1)
		go worker(&wg, tasksChan, errorsLimit, &errCounter)
	}

	for _, task := range tasks {
		tasksChan <- task
	}
	close(tasksChan)

	wg.Wait()

	if errCounter >= int32(errorsLimit) {
		return ErrErrorsLimitExceeded
	}

	return nil
}

func worker(wg *sync.WaitGroup, tasks <-chan Task, errLimit int, errCounter *int32) {
	defer wg.Done()

	for task := range tasks {
		if atomic.LoadInt32(errCounter) >= int32(errLimit) {
			return
		}

		if err := task(); err != nil {
			atomic.AddInt32(errCounter, 1)
		}
	}
}
