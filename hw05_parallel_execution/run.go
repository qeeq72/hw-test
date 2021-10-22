package hw05parallelexecution

import (
	"errors"
	"sync"
	"sync/atomic"
)

var ErrErrorsLimitExceeded = errors.New("errors limit exceeded")

type Task func() error

func getTryCount(errLimit int, taskCount int) int32 {
	if errLimit < 0 {
		return int32(taskCount)
	}

	if errLimit > taskCount {
		return int32(taskCount)
	}

	return int32(errLimit)
}

func Run(tasks []Task, n, m int) error {
	tryCount := getTryCount(m, len(tasks))
	tasksCh := make(chan Task, n)
	wg := &sync.WaitGroup{}

	for i := 0; i < n; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for task := range tasksCh {
				if atomic.LoadInt32(&tryCount) <= 0 {
					return
				}
				err := task()
				if err != nil {
					atomic.AddInt32(&tryCount, -1)
				}
			}
		}()
	}

	for i := range tasks {
		if atomic.LoadInt32(&tryCount) <= 0 {
			break
		}
		tasksCh <- tasks[i]
	}
	close(tasksCh)
	wg.Wait()
	if tryCount <= 0 {
		return ErrErrorsLimitExceeded
	}
	return nil
}
