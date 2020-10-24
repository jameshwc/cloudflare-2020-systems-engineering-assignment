package profile

import (
	"io"
	"sync"
	"time"

	"log"

	"github.com/jameshwc/simple-http/http"
	// "net/http"
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

func request(method, url string, body io.Reader, timeout time.Duration) (statusCode int, size, requestTime uint64, errCode string, err error) {
	req, err := http.NewRequest(method, url)
	if err != nil {
		log.Println(err)
		return
	}
	req.Header["Content-Type"] = "application/x-www-form-urlencoded; charset=utf-8"
	startTime := time.Now()
	// resp, err := http.Get(url)
	resp, err := req.Send()
	requestTime = uint64(time.Since(startTime))
	if err != nil {
		return
	}
	statusCode = resp.StatusCode
	errCode = resp.ErrorCode
	size = uint64(resp.Size)
	// size = calcResponseSize(resp)
	return
}

func getResponse(r *Request) (bool, int, uint64, uint64, string) {
	isSucceed := false
	statusCode, size, requestTime, errCode, err := request(r.Method, r.URL, r.BodyReader, r.Timeout)
	if err != nil {
		log.Println(err)
	} else {
		if statusCode == 200 && errCode == "" {
			isSucceed = true
		}
	}
	return isSucceed, statusCode, requestTime, size, errCode
}

// func calcResponseSize(r *http.Response) uint64 {
// 	n := uint64(0)
// 	if r.ContentLength != -1 {
// 		n += uint64(r.ContentLength)
// 	} else if r.Body != nil {
// 		body, err := io.Copy(ioutil.Discard, r.Body)
// 		if err != nil {
// 			return 0
// 		}
// 		n += uint64(body)
// 	}
// 	header := 0
// 	for name, values := range r.Header {
// 		header += len(name)
// 		for _, value := range values {
// 			header += len(value)
// 		}
// 	}
// 	if len(r.TransferEncoding) != 0 {
// 		header = header + len(r.TransferEncoding[0])
// 		header = header + len("Transfer-Encoding")
// 	}
// 	return n + uint64(header)
// }
