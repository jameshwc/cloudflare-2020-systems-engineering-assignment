package http

import (
	"crypto/tls"
	"io"
	"io/ioutil"
	"net/http"
	"sync"
	"time"

	"github.com/sirupsen/logrus"
)

func ConcurrencyRequest(chanID uint64, ch chan<- *Response, totalNumber uint64, wg *sync.WaitGroup, request *Request) {
	defer func() {
		wg.Done()
	}()
	for i := uint64(0); i < totalNumber; i++ {
		isSucceed, statusCode, requestTime, size := getResponse(request)
		result := NewResponse(i, chanID, requestTime, isSucceed, statusCode, size)
		ch <- result
	}
	return
}

func request(method, url string, body io.Reader, timeout time.Duration) (statusCode, size int, requestTime uint64, err error) {
	req, err := http.NewRequest(method, url, body)
	if err != nil {
		logrus.Info(err)
		return
	}
	client := &http.Client{
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		},
		Timeout: timeout,
	}
	startTime := time.Now()
	resp, err := client.Do(req)
	if err != nil {
		logrus.Error("Request failed: ", err)
	}
	defer resp.Body.Close()
	statusCode = resp.StatusCode
	size = calcResponseSize(resp)
	requestTime = uint64(time.Since(startTime))
	return
}

func getResponse(r *Request) (bool, int, uint64, int) {
	size := 0
	isSucceed := false
	statusCode, size, requestTime, err := request(r.Method, r.URL, r.BodyReader, r.Timeout)
	if err != nil {
		logrus.Error(err)
	} else {
		if statusCode == http.StatusOK {
			isSucceed = true
		}
	}
	return isSucceed, statusCode, requestTime, size
}

func calcResponseSize(r *http.Response) int {
	n := 0
	if r.ContentLength != -1 {
		n += int(r.ContentLength)
	} else {
		body, err := io.Copy(ioutil.Discard, r.Body)
		if err != nil {
			return 0
		}
		n += int(body)
	}
	header := 0
	for name, values := range r.Header {
		header += len(name)
		for _, value := range values {
			header += len(value)
		}
	}
	if len(r.TransferEncoding) != 0 {
		header = header + len(r.TransferEncoding[0])
		header = header + len("Transfer-Encoding")
	}
	return n + header
}
