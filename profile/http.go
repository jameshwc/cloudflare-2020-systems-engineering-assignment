package profile

import (
	"sync"
	"time"

	"log"

	"github.com/jameshwc/simple-http/http"
)

func ConcurrencyRequest(chanID uint64, ch chan<- *Response, totalNumber uint64, wg *sync.WaitGroup, request *Request) {
	defer func() {
		wg.Done()
	}()
	for i := uint64(0); i < totalNumber; i++ {
		isSucceed, statusCode, requestTime, size, errCode := getResponse(request)
		result := NewResponse(i, chanID, requestTime, size, isSucceed, statusCode, errCode)
		ch <- result
	}
	return
}

func request(method, url string) (statusCode int, size, requestTime uint64, errCode string, err error) {
	req, err := http.NewRequest(method, url)
	if err != nil {
		log.Println(err)
		return
	}
	req.Header["Content-Type"] = "application/x-www-form-urlencoded; charset=utf-8"
	startTime := time.Now()
	resp, err := req.Send()
	requestTime = uint64(time.Since(startTime))
	if err != nil {
		return
	}
	statusCode = resp.StatusCode
	errCode = resp.ErrorCode
	size = uint64(resp.Size)
	return
}

func getResponse(r *Request) (bool, int, uint64, uint64, string) {
	isSucceed := false
	statusCode, size, requestTime, errCode, err := request(r.Method, r.URL)
	if err != nil {
		log.Println(err)
	} else {
		if statusCode == 200 && errCode == "" {
			isSucceed = true
		}
	}
	return isSucceed, statusCode, requestTime, size, errCode
}
