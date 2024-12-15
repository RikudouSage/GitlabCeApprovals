package main

import (
	"GitlabCeForcedApprovals/controller"
	"GitlabCeForcedApprovals/dto"
	appHttp "GitlabCeForcedApprovals/http"
	"GitlabCeForcedApprovals/router"
	"GitlabCeForcedApprovals/worker"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	gitlab "gitlab.com/gitlab-org/api/client-go"
	"io"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"
)

var AppRouter *router.Router
var WorkerPool worker.Pool
var Gitlab *gitlab.Client

func LambdaHandler(request *events.APIGatewayV2HTTPRequest) (*events.APIGatewayProxyResponse, error) {
	if request == nil {
		return appHttp.ApiGatewayResponse(&appHttp.Response{
			StatusCode: http.StatusBadRequest,
			Body: &dto.Response{
				Success:   false,
				Timestamp: time.Now(),
				Reason:    "Request is empty",
			},
		}), nil
	}

	log.Println("Handling request to " + request.RawPath)

	appRequest := appHttp.NewRequest(
		[]byte(request.Body),
		request.RequestContext.HTTP.Method,
	)

	var err error
	var result *appHttp.Response
	for _, route := range AppRouter.Routes {
		if route.Path == request.RawPath {
			result, err = route.ControllerMethod(appRequest)
			break
		}
	}

	if err != nil {
		return appHttp.ApiGatewayResponse(appHttp.InternalServerError()), err
	}

	if result == nil {
		return appHttp.ApiGatewayResponse(appHttp.NotFoundResponse()), nil
	}

	return appHttp.ApiGatewayResponse(result), nil
}

func main() {
	if _, isLambda := syscall.Getenv("LAMBDA_TASK_ROOT"); isLambda {
		lambda.Start(LambdaHandler)
		return
	}

	gracefulShutdown := make(chan os.Signal, 1)
	signal.Notify(gracefulShutdown, syscall.SIGINT, syscall.SIGTERM)

	http.HandleFunc("/", func(writer http.ResponseWriter, request *http.Request) {
		body, err := io.ReadAll(request.Body)
		if err != nil {
			appHttp.WriteHttpResponse(appHttp.InternalServerError(), writer)
			log.Println(err)
			return
		}
		defer request.Body.Close()

		appRequest := appHttp.NewRequest(
			body,
			request.Method,
		)

		var result *appHttp.Response
		for _, route := range AppRouter.Routes {
			if route.Path == request.URL.Path {
				result, err = route.ControllerMethod(appRequest)
				break
			}
		}

		if err != nil {
			appHttp.WriteHttpResponse(&appHttp.Response{
				StatusCode: http.StatusInternalServerError,
				Body: &dto.Response{
					Success:   false,
					Timestamp: time.Now(),
					Reason:    "Internal server error",
				},
			}, writer)
			log.Println(err)
			return
		}

		if result == nil {
			appHttp.WriteHttpResponse(appHttp.NotFoundResponse(), writer)
			return
		}

		appHttp.WriteHttpResponse(result, writer)
	})

	go func() {
		port, exists := syscall.Getenv("HTTP_PORT")
		if !exists {
			port = "8080"
		}

		log.Println("Starting server at port " + port)
		err := http.ListenAndServe(":"+port, nil)
		if err != nil {
			log.Println(err)
			os.Exit(1)
		}
	}()

	<-gracefulShutdown
	log.Println("Shutting down server...")

	<-WorkerPool.ShutdownAndWait()
	log.Println("All workers finished, server shut down")
}

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
