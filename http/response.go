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
}

var (
	ErrContentLengthNaN = errors.New("Content-Length is not a number")
	ErrStatusCodeNaN    = errors.New("Status Code is not a number")
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
	}
	return &Response{statusCode, contentLength, body, header, []string{}, len(b), b}, nil
}
