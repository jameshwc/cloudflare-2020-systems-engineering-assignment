package main

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

type RequestResult struct {
	Id        string
	ChanId    uint64
	Time      uint64
	IsSucceed bool
	ErrCode   int
}

func NewRequest(url, method, body string, timeout time.Duration) *Request {
	return &Request{url, method, body, strings.NewReader(body), timeout}
}
