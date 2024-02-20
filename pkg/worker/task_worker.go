package worker

import (
	"context"
	"fmt"
	"sync"
)

type Result struct {
	Val interface{}
	Err error
}

type Pool struct {
	jobs         chan Job
	Results      chan Result
	WorkersCount int
}

type ExecutionFn func(ctx context.Context) (interface{}, error)

type Job struct {
	Descriptor string
	ExecFn     ExecutionFn
}

func (job Job) execute(ctx context.Context) Result {
	value, err := job.ExecFn(ctx)
	if err != nil {
		return Result{Err: err}
	}

	return Result{Val: value}
}

func NewPool(workersCount int) *Pool {
	return &Pool{
		WorkersCount: workersCount,
		jobs:         make(chan Job, workersCount),
		Results:      make(chan Result, workersCount),
	}
}

func (worker *Pool) listen(ctx context.Context, jobs <-chan Job) {
	for {
		select {
		case job, ok := <-jobs:
			if !ok {
				return
			}
			worker.jobs <- job
		case <-ctx.Done():
			return
		}
	}
}

func (worker *Pool) Run(ctx context.Context, jobs <-chan Job) {
	var wg sync.WaitGroup

	wg.Add(1)
	go func() {
		defer wg.Done()
		defer close(worker.jobs)
		worker.listen(ctx, jobs)

		fmt.Println("Finishing listening")
	}()

	wg.Add(worker.WorkersCount)
	for i := 0; i < worker.WorkersCount; i++ {
		go runJobs(ctx, &wg, worker.jobs, worker.Results)
	}

	wg.Wait()
}

func runJobs(ctx context.Context, wg *sync.WaitGroup, jobs <-chan Job, results chan<- Result) {
	defer wg.Done()
	for {
		select {
		case job, ok := <-jobs:
			if !ok {
				return
			}

			results <- job.execute(ctx)
		case <-ctx.Done():
			results <- Result{Err: ctx.Err()}
			return
		}
	}
}

func (worker *Pool) Close() {
	close(worker.Results)
}
