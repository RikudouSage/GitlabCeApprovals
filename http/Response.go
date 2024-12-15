package http

import (
	"GitlabCeForcedApprovals/dto"
	"fmt"
	"net/http"
	"time"
)

type Response struct {
	StatusCode int
	Body       any
	Headers    map[string]string
}

func NotFoundResponse() *Response {
	return &Response{
		StatusCode: http.StatusNotFound,
		Body: &dto.Response{
			Success:   false,
			Timestamp: time.Now(),
			Reason:    "Not found",
		},
	}
}

func MethodNotAllowed(provided string, expected string) *Response {
	return &Response{
		StatusCode: http.StatusMethodNotAllowed,
		Body: &dto.Response{
			Success:   false,
			Timestamp: time.Now(),
			Reason:    fmt.Sprintf("Method %s is not allowed, %s expected", provided, expected),
		},
	}
}

func Success(message string) *Response {
	return &Response{
		StatusCode: http.StatusOK,
		Body: &dto.Response{
			Success:   true,
			Timestamp: time.Now(),
			Reason:    message,
		},
	}
}

func InternalServerError() *Response {
	return &Response{
		StatusCode: http.StatusInternalServerError,
		Body: &dto.Response{
			Success:   false,
			Timestamp: time.Now(),
			Reason:    "Internal server error",
		},
	}
}
