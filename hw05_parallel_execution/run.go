package hw05parallelexecution

import (
	"errors"
	"sync"
	"sync/atomic"
)

var ErrErrorsLimitExceeded = errors.New("errors limit exceeded")

type Task func() error

// Run starts tasks in n goroutines and stops its work when receiving m errors from tasks.
func Run(tasks []Task, n, m int) error {
	if n <= 0 {
		return errors.New("specify the correct number of goroutines")
	}

	var counter int64
	var resultingError error

	maxGoroutines := n

	if len(tasks) < maxGoroutines {
		maxGoroutines = len(tasks)
	}

	maxErrors := int64(m)

	wg := &sync.WaitGroup{}
	wg.Add(maxGoroutines)

	done := make(chan struct{})
	jobs := make(chan Task)

	for i := 0; i < maxGoroutines; i++ {
		go func(done <-chan struct{}, jobs <-chan Task) {
			defer wg.Done()

			for {
				select {
				case <-done:
					return
				default:
					job, ok := <-jobs

					if !ok {
						return
					}

					if err := job(); err != nil {
						atomic.AddInt64(&counter, 1)
					}
				}
			}
		}(done, jobs)
	}

	for _, task := range tasks {
		if maxErrors > 0 && atomic.LoadInt64(&counter) >= maxErrors {
			close(done)
			resultingError = ErrErrorsLimitExceeded

			break
		}

		jobs <- task
	}

	close(jobs)

	wg.Wait()

	return resultingError
}
