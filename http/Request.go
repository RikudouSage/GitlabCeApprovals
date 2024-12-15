package http

type Request struct {
	Body   []byte
	Method string
}

func NewRequest(body []byte, method string) *Request {
	return &Request{Body: body, Method: method}
}
