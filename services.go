package main

import (
	"GitlabCeForcedApprovals/controller"
	"GitlabCeForcedApprovals/router"
	"GitlabCeForcedApprovals/worker"
	gitlab "gitlab.com/gitlab-org/api/client-go"
	"log"
	"os"
	"strconv"
	"syscall"
)

var AppRouter *router.Router
var WorkerPool worker.Pool
var Gitlab *gitlab.Client

func init() {
	var err error

	if _, isLambda := syscall.Getenv("LAMBDA_TASK_ROOT"); isLambda {
		WorkerPool = &worker.LambdaPool{}
	} else {
		workerCount, exists := syscall.Getenv("WORKER_COUNT")
		if !exists {
			workerCount = "100"
		}
		workerCountInt, err := strconv.Atoi(workerCount)
		if err != nil {
			panic(err)
		}

		log.Println("Starting a worker pool with " + workerCount + " workers")
		WorkerPool = worker.NewStandardPool(workerCountInt)
	}

	Gitlab, err = gitlab.NewClient(os.Getenv("GITLAB_ACCESS_TOKEN"), gitlab.WithBaseURL(os.Getenv("GITLAB_BASE_URL")))
	if err != nil {
		panic(err)
	}

	webhookController := &controller.WebhookController{Pool: WorkerPool, Gitlab: Gitlab}

	AppRouter = router.NewRouter()
	AppRouter.AddRoute(router.NewRoute("/webhooks", webhookController.MergeRequestEvent))
}
