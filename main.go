package main

import (
	"flag"
	"fmt"
	"log"
	"runtime"
	"time"

	"github.com/jameshwc/simple-http/http"
	"github.com/jameshwc/simple-http/profile"
	"github.com/jameshwc/simple-http/stress"
)

func main() {
	runtime.GOMAXPROCS(1)

	var (
		concurrency uint64
		totalNumber uint64
		requestURL  string
	)

	flag.StringVar(&requestURL, "url", "", "(required) url")
	flag.Uint64Var(&totalNumber, "profile", 0, "(optional) profile mode: assign the number of requests")
	flag.Uint64Var(&concurrency, "c", 1, "(optional) profile mode: number of concurrency")
	flag.Parse()
	if requestURL == "" {
		fmt.Printf("Usage: ./simple-http -u http://cloudflare-2020-general-engineering-assignment.jameshwc.workers.dev/links \n")
		fmt.Printf("Profile Usage: ./simple-http -u http://cloudflare-2020-general-engineering-assignment.jameshwc.workers.dev/links -profile 100 [-c 10 (default 1)]\n")
		fmt.Printf("The parameter you assigned: -c %d -profile %d -u %s \n", concurrency, totalNumber, requestURL)
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
