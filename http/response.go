package http

import (
	"errors"
	"strconv"
	"strings"
)

type Response struct {
	StatusCode       int
	ContentLength    int
	Body             string
	Header           map[string][]string
	TransferEncoding []string
	Size             int
	rawData          []byte
	ErrorCode        string // cloudflare
}

var (
	ErrContentLengthNaN      = errors.New("Content-Length is not a number")
	ErrStatusCodeNaN         = errors.New("Status Code is not a number")
	ErrChunkedNotImplemented = errors.New("Chunked transfer-encoding not implemented yet")
)

func ParseResponse(b []byte) (*Response, error) {
	split := strings.Split(string(b), "\r\n")
	body := ""
	header := make(map[string][]string)
	statusCode := 0
	contentLength := -1
	var err error
	for id, h := range split {
		if id == 0 {
			s := strings.Split(h, " ")
			statusCode, err = strconv.Atoi(s[1])
			if err != nil {
				return nil, ErrStatusCodeNaN
			}
			continue
		}
		if len(h) == 0 {
			body = split[id+1]
			break
		}
		s := strings.Split(h, ":")
		name := s[0]
		value := strings.TrimSpace(s[1])
		if _, ok := header[name]; !ok {
			header[name] = make([]string, 0)
		}
		header[name] = append(header[name], value)
		if name == "Content-Length" {
			contentLength, err = strconv.Atoi(value)
			if err != nil {
				return nil, ErrContentLengthNaN
			}
		}
		if name == "Transfer-Encoding" && value == "chunked" {
			return nil, ErrChunkedNotImplemented
		}
	}
	errCode := ""
	if strings.HasPrefix(body, "error code") {
		errCode = strings.Split(body, "error code: ")[1]
	}
	return &Response{statusCode, contentLength, body, header, []string{}, len(b), b, errCode}, nil
}
