package worker

import "GitlabCeForcedApprovals/worker/job"

type Pool interface {
	ShutdownAndWait() chan int
	EnqueueJob(job job.HandleableJob) chan bool
}
