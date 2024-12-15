package worker

import (
	"GitlabCeForcedApprovals/worker/job"
	"time"
)

type StandardPool struct {
	Workers []*worker
}

func NewStandardPool(workers int) *StandardPool {
	pool := &StandardPool{
		Workers: make([]*worker, workers),
	}
	for i := range workers {
		pool.Workers[i] = &worker{}
	}

	return pool
}

func (pool *StandardPool) ShutdownAndWait() chan int {
	channel := make(chan int)

	go func() {
		for {
			allIdle := true
			for _, workerObject := range pool.Workers {
				if !workerObject.IsIdle() {
					allIdle = false
					break
				}
			}
			if allIdle {
				break
			}

			time.Sleep(100 * time.Millisecond)
		}

		channel <- 1
		close(channel)
	}()

	return channel
}

func (pool *StandardPool) EnqueueJob(job job.HandleableJob) chan bool {
	channel := make(chan bool)

	go func() {
		for {
			for _, worker := range pool.Workers {
				if worker.IsIdle() {
					worker.RunJob(job)
					channel <- true
					close(channel)
					return
				}
			}
			time.Sleep(100 * time.Millisecond)
		}
	}()

	return channel
}
