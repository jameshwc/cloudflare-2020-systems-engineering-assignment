package http

import (
	"bytes"
	"errors"
	"strconv"
	"strings"
)

type Response struct {
	StatusCode   int
	Body         string
	Header       map[string][]string
	ReasonPhrase string
	Size         int
	rawData      []byte
	ErrorCode    string // cloudflare
}

var (
	ErrContentLengthNaN                   = errors.New("Response Content-Length is not a number")
	ErrStatusCodeNaN                      = errors.New("Response Status Code is not a number")
	ErrChunkedNotImplemented              = errors.New("Response Chunked transfer-encoding not implemented yet")
	ErrChunksReadError                    = errors.New("Response Chunks read error, please submit a issue")
	ErrHeaderKeyValue                     = errors.New("Response headers not follow {$key: $value} format")
	ErrHTTPResponseFormatError            = errors.New("Response not follow HTTP/1.1 protocol (RFC 2616)")
	ErrContentLengthNegativeAndNotChunked = errors.New("Response Content-Length negative and Transfer-Encoding not chunked")
	ErrContentLengthError                 = errors.New("Response Content Length incorrect")
)

// https://tools.ietf.org/html/rfc2616#section-6.1
func ParseResponse(b []byte) (*Response, error) {

	dat := bytes.SplitN(b, []byte{'\r', '\n', '\r', '\n'}, 2)
	if len(dat) < 2 {
		return nil, ErrHTTPResponseFormatError
	}

	statusCode, reasonPhrase, err := parseStatusLine(dat[0])
	if err != nil {
		return nil, err
	}

	headers, err := parseHeaders(dat[0])
	if err != nil {
		return nil, err
	}

	var body string

	if v, ok := headers["Transfer-Encoding"]; ok && v[0] == "chunked" {

		body, err = parseChunks(dat[1])
		if err != nil {
			return nil, err
		}

	} else if v, ok := headers["Content-Length"]; ok {

		bodySize, err := strconv.Atoi(v[0])
		if err != nil {
			return nil, ErrContentLengthNaN
		}

		if bodySize < 0 {
			return nil, ErrContentLengthNegativeAndNotChunked
		}

		if len(dat[1]) != bodySize {
			return nil, ErrContentLengthError
		}

		body = string(dat[1])

	}
	errorCode := parseErrorCode(body)
	return &Response{statusCode, body, headers, reasonPhrase, len(b), b, errorCode}, nil
}

func parseStatusLine(b []byte) (int, string, error) {

	s := bytes.SplitN(b, []byte{' '}, 3)
	if len(s) < 3 {
		return 0, "", ErrHTTPResponseFormatError
	}

	statusCode, err := strconv.Atoi(string(s[1]))
	if err != nil {
		return 0, "", ErrStatusCodeNaN
	}

	return statusCode, string(s[2]), nil
}
func parseHeaders(dat []byte) (map[string][]string, error) {

	h := make(map[string][]string)
	headers := bytes.Split(dat, []byte{'\r', '\n'})

	for i := range headers {

		if i == 0 { // skip status line
			continue
		}

		s := bytes.SplitN(headers[i], []byte{':', ' '}, 2)
		if len(s) < 2 {
			return nil, ErrHeaderKeyValue
		}

		key := string(s[0])
		val := string(s[1])

		if _, ok := h[key]; !ok {
			h[key] = make([]string, 0)
		}

		h[key] = append(h[key], val)
	}

	return h, nil
}

func parseChunks(b []byte) (string, error) {
	i := 0
	body := ""
	var hex []byte

	for {
		if b[i] != '\r' {
			hex = append(hex, b[i])
			i++
		} else {
			i += 2
			size, err := strconv.ParseInt(string(hex), 16, 64)
			if err != nil {
				return "", ErrChunksReadError
			}
			if size == 0 {
				break
			}

			if i+int(size)+2 > len(b) { // add len("\r\n")
				return "", ErrChunksReadError
			}

			body = body + string(b[i:i+int(size)])
			i += int(size) + 2 // add size and len("\r\n")
			hex = make([]byte, 0)
		}
	}
	return body, nil
}

func parseErrorCode(body string) (errCode string) {
	if strings.HasPrefix(body, "error code") {
		s := strings.SplitN(body, "error code: ", 2)
		if len(s) == 2 {
			errCode = s[1]
		}
	}
	return
}
