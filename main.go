package main

import (
	"flag"
	"fmt"
	"runtime"
	"time"
)

func main() {
	runtime.GOMAXPROCS(1)

	var (
		concurrency uint64
		totalNumber uint64
		requestURL  string
	)

	flag.Uint64Var(&concurrency, "c", 1, "number of concurrency")
	flag.Uint64Var(&totalNumber, "n", 1, "number of requests")
	flag.StringVar(&requestURL, "u", "", "url")
	flag.Parse()

	if concurrency == 0 || totalNumber == 0 || requestURL == "" {
		fmt.Printf("Usage: go run main.go -c 1 -n 1 -u https://www.google.com/ \n")
		fmt.Printf("The parameter you assigned: -c %d -n %d -u %s \n", concurrency, totalNumber, requestURL)
		flag.Usage()
		return
	}

	StartStressTest(concurrency, totalNumber, NewRequest(requestURL, "GET", "hello, world", 30*time.Second))
}
