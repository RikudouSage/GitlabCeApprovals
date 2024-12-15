package http

import (
	"GitlabCeForcedApprovals/json"
	"log"
	"net/http"
)

func WriteHttpResponse(response *Response, writer http.ResponseWriter) {
	body := response.Body
	headers := response.Headers
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
		body, err = json.ToJson(body)
		if err != nil {
			body, _ = json.ToJson(map[string]string{
				"error": "Internal request error",
			})
			statusCode = http.StatusInternalServerError
			log.Println(err)
		}
	}

	for key, value := range headers {
		writer.Header().Set(key, value)
	}
	writer.WriteHeader(statusCode)

	_, err = writer.Write(body.([]byte))
	if err != nil {
		log.Println(err)
	}
}
