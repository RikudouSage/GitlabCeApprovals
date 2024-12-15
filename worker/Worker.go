package worker

import (
	"GitlabCeForcedApprovals/worker/job"
	"fmt"
)

type worker struct {
	job job.HandleableJob
}

func (worker *worker) IsIdle() bool {
	return worker.job == nil
}

func (worker *worker) RunJob(job job.HandleableJob) {
	if job == nil {
		panic("job cannot be nil")
	}
	worker.job = job

	go worker.run()
}

func (worker *worker) run() {
	defer func() {
		worker.job = nil
	}()
	if !<-worker.job.Handle() {
		fmt.Println("There was an error handling a job")
	}
}
