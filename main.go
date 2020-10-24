package main

import (
	"flag"
	"fmt"
	"log"
	"runtime"
	"time"

	"github.com/jameshwc/go-stress/http"
	"github.com/jameshwc/go-stress/profile"
	"github.com/jameshwc/go-stress/stress"
)

func main() {
	runtime.GOMAXPROCS(1)

	var (
		concurrency uint64
		totalNumber uint64
		requestURL  string
	)

	flag.Uint64Var(&totalNumber, "profile", 0, "profile mode: assign the number of requests")
	flag.Uint64Var(&concurrency, "c", 1, "profile mode: number of concurrency")
	flag.StringVar(&requestURL, "url", "", "url")
	flag.Parse()
	if requestURL == "" {
		fmt.Printf("Usage: go run main.go -c 1 -n 1 -u https://www.google.com/ \n")
		fmt.Printf("The parameter you assigned: -c %d -n %d -u %s \n", concurrency, totalNumber, requestURL)
		flag.Usage()
		return
	}
	if totalNumber > 0 {
		stress.StartTest(concurrency, totalNumber, profile.NewRequest(requestURL, "GET", "", 30*time.Second))
	} else {
		req, err := http.NewRequest("GET", requestURL)
		if err != nil {
			log.Fatal(err)
		}
		resp, err := req.Send()
		if err != nil {
			log.Fatal(err)
		}
		if resp.ErrorCode != "" {
			fmt.Println("error code: ", resp.ErrorCode)
		} else {
			fmt.Println(resp.Body)
		}
	}
}
