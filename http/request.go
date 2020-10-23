package http

import (
	"io"
	"strings"
	"time"
)

type Request struct {
	URL        string
	Method     string
	Body       string
	BodyReader io.Reader
	Timeout    time.Duration
}

func NewRequest(url, method, body string, timeout time.Duration) *Request {
	return &Request{url, method, body, strings.NewReader(body), timeout}
}
