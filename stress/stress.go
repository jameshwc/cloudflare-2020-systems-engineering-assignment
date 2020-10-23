package stress

import (
	"fmt"
	"sync"
	"time"

	"github.com/jameshwc/go-stress/http"
	"github.com/jameshwc/go-stress/stat"
)

func StartTest(concurrency, totalNumber uint64, request *http.Request) {
	ch := make(chan *http.Response, 1000)
	var wg sync.WaitGroup
	var receive sync.WaitGroup
	receive.Add(1)
	go stat.Receive(concurrency, ch, &receive)
	for i := uint64(0); i < concurrency; i++ {
		fmt.Println(i)
		wg.Add(1)
		go http.ConcurrencyRequest(i, ch, totalNumber, &wg, request)
	}
	fmt.Println("HIIIIIIIIIIIIIIIIIIIIIIIIIIIIIi")
	wg.Wait()
	fmt.Println("HIIIIIIIIIIIIIIIIIIIIIIIIIIIIIi")
	time.Sleep(1 * time.Millisecond)
	close(ch)
	receive.Wait()
	return
}
