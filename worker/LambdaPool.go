package worker

import (
	"GitlabCeForcedApprovals/worker/job"
)

type LambdaPool struct {
}

func (receiver *LambdaPool) ShutdownAndWait() chan int {
	result := make(chan int)

	go func() {
		defer close(result)
		result <- 0
	}()

	return result
}

func (receiver *LambdaPool) EnqueueJob(job job.HandleableJob) chan bool {
	jobResult := <-job.Handle()

	result := make(chan bool)

	go func() {
		defer close(result)
		result <- jobResult
	}()

	return result
}
