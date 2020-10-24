package http

import (
	"errors"
	"fmt"
	"io/ioutil"
	"net"
	"net/url"
	"strconv"
)

type Request struct {
	Header map[string]string
	Host   string
	URL    *url.URL
	Method string
	// Body   io.Reader
}

var (
	ErrURLFormatIncorrect  = errors.New("url format not correct")
	ErrHttpsNotImplemented = errors.New("https not implemented yet")
	ErrVerbNotImplemented  = errors.New("valid verb but not implemented yet")
	ErrVerbInvalid         = errors.New("invalid http verb")
)

func NewRequest(method, URL string) (*Request, error) {
	u, err := url.Parse(URL)
	if err != nil {
		return nil, ErrURLFormatIncorrect
	}

	port := 0
	switch u.Scheme {
	case "http":
		port = 80
	case "https":
		// port = 443
		return nil, ErrHttpsNotImplemented
	}

	switch method {
	case "GET":
	case "POST", "DELETE", "PUT", "HEAD", "PATCH", "OPTIONS", "TRACE", "CONNECT":
		return nil, ErrVerbNotImplemented
	default:
		return nil, ErrVerbInvalid
	}

	return &Request{make(map[string]string), u.Host + ":" + strconv.Itoa(port), u, method}, nil
}

func (r *Request) Send() (*Response, error) {
	conn, err := net.Dial("tcp", r.Host)
	if err != nil {
		return nil, err
	}
	path := r.URL.Path
	if path == "" {
		path = "/"
	}
	dat := fmt.Sprintf("GET %s HTTP/1.1\r\n", path)
	dat += fmt.Sprintf("Host: %v\r\n", r.URL.Host)
	dat += fmt.Sprintf("Connection: close\r\n")
	dat += fmt.Sprintf("\r\n")
	_, err = conn.Write([]byte(dat))

	if err != nil {
		return nil, err
	}

	resp, err := ioutil.ReadAll(conn)
	if err != nil {
		return nil, err
	}
	// fmt.Println(string(resp))
	conn.Close()
	return ParseResponse(resp)
}
