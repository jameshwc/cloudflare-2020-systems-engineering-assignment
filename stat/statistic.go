package stat

import (
	"fmt"
	"sort"
	"strings"
	"sync"
	"time"

	"github.com/jameshwc/go-stress/http"
)

var exportPeriod = 1 * time.Second

func Receive(concurrent uint64, ch <-chan *http.Response, wg *sync.WaitGroup) {
	defer func() {
		wg.Done()
	}()

	stopChan := make(chan bool)

	var (
		processingTime uint64
		requestTime    uint64
		maxTime        uint64
		minTime        uint64
		successNum     uint64
		failureNum     uint64
		chanSize       int
		chanIDs        = make(map[uint64]bool)
	)

	startTime := uint64(time.Now().UnixNano())

	var statusCode = make(map[int]int)

	ticker := time.NewTicker(exportPeriod)
	go func() {
		for {
			select {
			case <-ticker.C:
				endTime := uint64(time.Now().UnixNano())
				requestTime = endTime - startTime
				go calculateData(concurrent, processingTime, requestTime, maxTime, minTime, successNum, failureNum, chanSize, statusCode)
			case <-stopChan:
				return
			}
		}
	}()

	header()

	for data := range ch {
		processingTime = processingTime + data.Time

		if maxTime <= data.Time {
			maxTime = data.Time
		}

		if minTime == 0 {
			minTime = data.Time
		} else if minTime > data.Time {
			minTime = data.Time
		}

		if data.IsSucceed == true {
			successNum = successNum + 1
		} else {
			failureNum = failureNum + 1
		}
		code := data.StatusCode
		if value, ok := statusCode[code]; ok {
			statusCode[code] = value + 1
		} else {
			statusCode[code] = 1
		}

		if _, ok := chanIDs[data.ChanID]; !ok {
			chanIDs[data.ChanID] = true
			chanSize++
		}
	}

	stopChan <- true

	endTime := uint64(time.Now().UnixNano())
	requestTime = endTime - startTime

	calculateData(concurrent, processingTime, requestTime, maxTime, minTime, successNum, failureNum, chanSize, statusCode)

	fmt.Printf("\n\n")

	fmt.Println("*************************  结果 stat  ****************************")
	fmt.Println("处理协程数量:", concurrent)
	// fmt.Println("处理协程数量:", concurrent, "程序处理总时长:", fmt.Sprintf("%.3f", float64(processingTime/concurrent)/1e9), "秒")
	fmt.Println("请求总数（并发数*请求数 -c * -n）:", successNum+failureNum, "总请求时间:", fmt.Sprintf("%.3f", float64(requestTime)/1e9),
		"秒", "successNum:", successNum, "failureNum:", failureNum)

	fmt.Println("*************************  结果 end   ****************************")

	fmt.Printf("\n\n")
}

func calculateData(concurrent, processingTime, requestTime, maxTime, minTime, successNum, failureNum uint64, chanSize int, statusCode map[int]int) {
	if processingTime == 0 {
		processingTime = 1
	}

	var (
		averageTime      float64
		maxTimeFloat     float64
		minTimeFloat     float64
		requestTimeFloat float64
	)

	if successNum != 0 && concurrent != 0 {
		averageTime = float64(processingTime) / float64(successNum*1e6)
	}

	maxTimeFloat = float64(maxTime) / 1e6
	minTimeFloat = float64(minTime) / 1e6
	requestTimeFloat = float64(requestTime) / 1e9

	table(successNum, failureNum, statusCode, averageTime, maxTimeFloat, minTimeFloat, requestTimeFloat, chanSize)
}

func header() {
	fmt.Printf("\n\n")
	fmt.Println("─────┬───────┬───────┬────────┬────────┬────────┬────────────┬────────────┬────────────┬───────────")
	fmt.Println(" Time│ Concur│  Total| Success│ Failure│ size   │Slowest Time│Fastest Time│Average Time│Status Code")
	fmt.Println("─────┼───────┼───────┼────────┼────────┼────────┼────────────┼────────────┼────────────┼────────────")
}
func table(successNum, failureNum uint64, statusCode map[int]int, averageTime, maxTimeFloat, minTimeFloat, requestTimeFloat float64, chanSize int) {
	result := fmt.Sprintf("%4.0fs│%7d│%7d|%7d│%7d│%8.2f│%8.2f│%8.2f│%8.2f│%v", requestTimeFloat, chanSize, successNum+failureNum, successNum, failureNum, maxTimeFloat, minTimeFloat, averageTime, printMap(statusCode))
	fmt.Println(result)
}

func printMap(statusCode map[int]int) (mapStr string) {

	var mapArr []string
	for key, value := range statusCode {
		mapArr = append(mapArr, fmt.Sprintf("%d:%d", key, value))
	}

	sort.Strings(mapArr)

	mapStr = strings.Join(mapArr, ";")

	return
}
