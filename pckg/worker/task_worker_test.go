package worker

import (
	"context"
	"fmt"
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func MockWorkerFunction(_ context.Context) (interface{}, error) {
	return "ok", nil
}

func Test_WorkerCompletesJobs(t *testing.T) {
	worker := NewPool(2)
	jobsQueue := make(chan Job)
	ctx := context.TODO()
	testingContext, cancel := context.WithTimeout(ctx, 2*time.Second)
	defer cancel()
	var wg sync.WaitGroup

	wg.Add(1)
	go func() {
		defer wg.Done()
		defer close(jobsQueue)

		for i := 0; i < 10; i++ {
			jobsQueue <- Job{
				Descriptor: fmt.Sprintf("MockFunc: %d", i),
				ExecFn:     MockWorkerFunction,
			}
		}

		fmt.Println("Finished writing jobs")
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		defer worker.Close()

		worker.Run(testingContext, jobsQueue)
		fmt.Println("Finished running jobs")
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		workerJobCount := 0
		for result := range worker.Results {
			workerJobCount++
			assert.NoError(t, result.Err)

			assert.Equal(t, "ok", result.Val)
		}
		assert.Equal(t, 10, workerJobCount)
	}()

	wg.Wait()
}
