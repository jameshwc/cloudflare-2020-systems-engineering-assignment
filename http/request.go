package http

import (
	"crypto/tls"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"net"
	"net/url"
	"strconv"
	"strings"
)

type Request struct {
	Header map[string]string
	Host   string
	URL    *url.URL
	Method string
	Https  bool
}

var (
	ErrURLFormatIncorrect = errors.New("url format not correct")
	ErrVerbNotImplemented = errors.New("valid verb but not implemented yet")
	ErrVerbInvalid        = errors.New("invalid http verb")
)

func NewRequest(method, URL string) (*Request, error) {
	u, err := url.Parse(URL)
	if err != nil {
		return nil, ErrURLFormatIncorrect
	}

	port := 0
	https := false
	switch u.Scheme {
	case "http":
		port = 80
	case "https":
		port = 443
		https = true
		// return nil, ErrHttpsNotImplemented
	}

	method = strings.ToUpper(method)
	switch method {
	case "GET":
	case "POST", "DELETE", "PUT", "HEAD", "PATCH", "OPTIONS", "TRACE", "CONNECT":
		return nil, ErrVerbNotImplemented
	default:
		return nil, ErrVerbInvalid
	}

	return &Request{make(map[string]string), u.Host + ":" + strconv.Itoa(port), u, method, https}, nil
}

func (r *Request) Send() (*Response, error) {
	if !r.Https {
		conn, err := net.Dial("tcp", r.Host)
		if err != nil {
			return nil, err
		}
		defer conn.Close()
		return r.sendData(conn)
	} else {
		conn, err := tls.Dial("tcp", r.Host, &tls.Config{})
		if err != nil {
			return nil, err
		}
		defer conn.Close()
		return r.sendData(conn)
	}
}

func (r *Request) sendData(conn io.ReadWriteCloser) (*Response, error) {
	path := r.URL.Path
	if path == "" {
		path = "/"
	}
	dat := fmt.Sprintf("GET %s HTTP/1.1\r\n", path)
	dat += fmt.Sprintf("Host: %v\r\n", r.URL.Host)
	dat += fmt.Sprintf("Connection: close\r\n")
	dat += fmt.Sprintf("\r\n")
	_, err := conn.Write([]byte(dat))
	if err != nil {
		return nil, err
	}

	resp, err := ioutil.ReadAll(conn)
	if err != nil {
		return nil, err
	}
	return ParseResponse(resp)
}
