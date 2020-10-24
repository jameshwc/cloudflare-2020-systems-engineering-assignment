package stress

import (
	"sync"
	"time"

	"github.com/jameshwc/go-stress/profile"
	"github.com/jameshwc/go-stress/stat"
)

func StartTest(concurrency, totalNumber uint64, request *profile.Request) {
	ch := make(chan *profile.Response, 10000)
	var wg sync.WaitGroup
	var receive sync.WaitGroup
	receive.Add(1)
	go stat.Receive(concurrency, ch, &receive)
	for i := uint64(0); i < concurrency; i++ {
		wg.Add(1)
		go profile.ConcurrencyRequest(i, ch, totalNumber, &wg, request)
	}
	wg.Wait()
	time.Sleep(1 * time.Millisecond)
	close(ch)
	receive.Wait()
	return
}
