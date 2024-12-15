package controller

import (
	"GitlabCeForcedApprovals/http"
	"GitlabCeForcedApprovals/json"
	"GitlabCeForcedApprovals/worker"
	"GitlabCeForcedApprovals/worker/job"
	gitlab "gitlab.com/gitlab-org/api/client-go"
	goHttp "net/http"
)

type WebhookController struct {
	Pool   worker.Pool
	Gitlab *gitlab.Client
}

func (receiver *WebhookController) MergeRequestEvent(request *http.Request) (*http.Response, error) {
	if request.Method != goHttp.MethodPost {
		return http.MethodNotAllowed(request.Method, goHttp.MethodPost), nil
	}

	var event gitlab.MergeEvent
	err := json.Map(request.Body, &event)
	if err != nil {
		return nil, err
	}

	receiver.Pool.EnqueueJob(&job.MergeEventHandler{Event: &event, Gitlab: receiver.Gitlab})

	return http.Success("Success"), nil
}
