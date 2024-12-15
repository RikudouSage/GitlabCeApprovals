package http

import (
	"GitlabCeForcedApprovals/json"
	"github.com/aws/aws-lambda-go/events"
	"log"
	"net/http"
)

func ApiGatewayResponse(response *Response) *events.APIGatewayProxyResponse {
	headers := response.Headers
	body := response.Body
	statusCode := response.StatusCode

	if headers == nil {
		headers = make(map[string]string)
	}
	if body == nil {
		body = make(map[string]string)
	}
	if statusCode == 0 {
		statusCode = http.StatusOK
	}

	_, ok := headers["Content-Type"]
	if !ok {
		headers["Content-Type"] = "application/json"
	}

	var err error
	if _, ok = body.(string); !ok {
		body, err = json.ToJsonString(body)
		if err != nil {
			body, _ = json.ToJsonString(map[string]string{
				"error": "Internal request error",
			})
			statusCode = http.StatusInternalServerError
			log.Println(err)
		}
	}

	return &events.APIGatewayProxyResponse{
		StatusCode: statusCode,
		Body:       body.(string),
		Headers:    headers,
	}
}
