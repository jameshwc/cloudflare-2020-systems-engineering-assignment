package main

import (
	"fmt"
	"sync"
	"time"
)

func StartStressTest(concurrency, totalNumber uint64, request *Request) {
	ch := make(chan *RequestResult, 1000)
	var wg sync.WaitGroup
	var receive sync.WaitGroup
	receive.Add(1)
	go ReceivingResults(concurrency, ch, &receive)
	for i := uint64(0); i < concurrency; i++ {
		fmt.Println(i)
		wg.Add(1)
		go Http(i, ch, totalNumber, &wg, request)
	}
	wg.Wait()
	time.Sleep(1 * time.Millisecond)
	close(ch)
	receive.Wait()
	return
}
