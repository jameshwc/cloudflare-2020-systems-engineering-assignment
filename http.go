package main

import (
	"crypto/tls"
	"fmt"
	"io"
	"net/http"
	"sync"
	"time"

	"github.com/sirupsen/logrus"
)

func HttpRequest(method, url string, body io.Reader, timeout time.Duration) (resp *http.Response, requestTime uint64, err error) {
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
	resp, err = client.Do(req)
	requestTime = uint64(time.Since(startTime))
	if err != nil {
		logrus.Error("Request failed: ", err)
	}
	return
}

func Http(chanID uint64, ch chan<- *RequestResult, totalNumber uint64, wg *sync.WaitGroup, request *Request) {
	defer func() {
		wg.Done()
	}()
	for i := uint64(0); i < totalNumber; i++ {
		isSucceed, errCode, requestTime := send(request)
		result := &RequestResult{
			Id:        fmt.Sprintf("%d_%d", chanID, i),
			ChanId:    i,
			Time:      requestTime,
			IsSucceed: isSucceed,
			ErrCode:   errCode,
		}
		ch <- result
	}
}

func send(request *Request) (bool, int, uint64) {
	var (
		// startTime = time.Now()
		isSucceed = false
		errCode   = 200
	)

	resp, requestTime, err := HttpRequest(request.Method, request.URL, request.BodyReader, request.Timeout)
	if err != nil {
		logrus.Error(err)
		errCode = 509 // 请求错误
	} else {
		defer resp.Body.Close()
		if resp.StatusCode == http.StatusOK {
			isSucceed = true
		} else {
			errCode = resp.StatusCode
		}
	}
	return isSucceed, errCode, requestTime
}
