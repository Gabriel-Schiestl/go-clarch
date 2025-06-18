package workerpool

import (
	"errors"
	"sync"
)

type Result[TResult any] struct {
	Value TResult
	Err   error
}

type Job[P, R any] func(job P) (R, error)

type WorkerPool[P, R any] interface {
	Start() error
	Stop() error
	Results() <-chan Result[R]
	AddTasks(tasks ...P)
}

type workerPool[P, R any] struct {
	numWorkers int
	job        Job[P, R]
	resultChan chan Result[R]
	tasksChan  chan P
	wg         *sync.WaitGroup
	once       *sync.Once
	started    bool
}

func NewWorkerPool[P, R any](numWorkers int, job Job[P, R], bufferSize int) WorkerPool[P, R] {
	return &workerPool[P, R]{
		numWorkers: numWorkers,
		job:        job,
		resultChan: make(chan Result[R], bufferSize),
		tasksChan:  make(chan P, bufferSize),
		wg:         &sync.WaitGroup{},
		once:       &sync.Once{},
	}
}

func (wp *workerPool[P, R]) Start() error {
	if wp.started {
		return errors.New("worker pool already started")
	}

	wp.once.Do(func() {
		for i := 0; i < wp.numWorkers; i++ {
			wp.wg.Add(1)
			go func() {
				for job := range wp.tasksChan {
					wp.work(job)
				}
			}()
		}
		wp.started = true
	})

	return nil
}

func (wp *workerPool[P, R]) work(job P) {
	defer wp.wg.Done()
	result, err := wp.job(job)
	if err != nil {
		wp.resultChan <- Result[R]{Err: err}
		return
	}

	wp.resultChan <- Result[R]{Value: result}
}

func (wp *workerPool[P, R]) Results() <-chan Result[R] {
	return wp.resultChan
}

func (wp *workerPool[P, R]) Stop() error {
	if !wp.started {
		return errors.New("worker pool not started")
	}

	close(wp.tasksChan)
	wp.wg.Wait()
	close(wp.resultChan)

	wp.once = &sync.Once{}
	wp.started = false
	return nil
}

func (wp *workerPool[P, R]) AddTasks(tasks ...P) {
	for _, task := range tasks {
		wp.tasksChan <- task
	}
}
