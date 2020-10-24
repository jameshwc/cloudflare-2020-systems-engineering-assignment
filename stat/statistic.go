package stat

import (
	"fmt"
	"sort"
	"strconv"
	"strings"
	"sync"
	"time"

	"log"

	"github.com/jameshwc/simple-http/profile"
	"github.com/jameshwc/simple-http/stat/median"
)

var exportPeriod = 1 * time.Second

func Receive(concurrent uint64, ch <-chan *profile.Response, wg *sync.WaitGroup) {
	defer func() {
		wg.Done()
	}()

	stopChan := make(chan bool)

	var (
		processingTime uint64
		requestTime    uint64
		maxTime        uint64
		minTime        uint64
		minSize        uint64
		maxSize        uint64
		successNum     uint64
		failureNum     uint64
		medianTime     float64
		chanSize       int
		chanIDs        = make(map[uint64]bool)
	)

	medianFinder := median.NewMedianFinder()
	startTime := uint64(time.Now().UnixNano())

	var statusCode = make(map[int]int)
	var errorCode = make(map[int]int)

	ticker := time.NewTicker(exportPeriod)
	go func() {
		for {
			select {
			case <-ticker.C:
				endTime := uint64(time.Now().UnixNano())
				requestTime = endTime - startTime
				go calculateData(concurrent, processingTime, requestTime, maxTime, minTime, successNum, failureNum, maxSize, minSize, chanSize, medianTime, statusCode, errorCode)
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

		medianFinder.AddNum(median.HeapType(data.Time))

		if maxSize <= data.Size {
			maxSize = data.Size
		}

		if minSize == 0 {
			minSize = data.Size
		} else if minSize > data.Size {
			minSize = data.Size
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
		if data.ErrorCode != "" {

			errCode, err := strconv.Atoi(data.ErrorCode)
			if err != nil {
				log.Println(err)
			}

			if value, ok := errorCode[errCode]; ok {
				errorCode[errCode] = value + 1
			} else {
				errorCode[errCode] = 1
			}

		}

		if _, ok := chanIDs[data.ChanID]; !ok {
			chanIDs[data.ChanID] = true
			chanSize++
		}
		medianTime = medianFinder.FindMedian()
	}

	stopChan <- true

	endTime := uint64(time.Now().UnixNano())
	requestTime = endTime - startTime
	medianTime = medianFinder.FindMedian()
	calculateData(concurrent, processingTime, requestTime, maxTime, minTime, successNum, failureNum, maxSize, minSize, chanSize, medianTime, statusCode, errorCode)
	total := successNum + failureNum

	fmt.Printf("\n\n")

	fmt.Printf("*************************  Statistics  ****************************\n")
	fmt.Printf("Total Concurrent Number: %d\n", concurrent)
	fmt.Printf("Total Requests:( -c %d * -n %d）: %d\n", concurrent, total/concurrent, total)
	fmt.Printf("Fastest/Slowest Time (ms): %.2f/%.2f\n", float64(minTime)/1e6, float64(maxTime)/1e6)
	fmt.Printf("Mean/Median Time (ms): %.2f/%.2f\n", float64(processingTime)/float64(total*1e6), medianTime/1e6)
	fmt.Printf("Success Requests: %d (%3.f%%)\n", successNum, float64(successNum)/float64(total)*100)
	fmt.Printf("Smallest/Largest Response Size (bytes): %d/%d\n", minSize, maxSize)
	fmt.Printf("All Status Code: %s\n", printMap(statusCode))
	fmt.Printf("All Error Code: %s\n", printMap(errorCode))
	fmt.Println("*******************************************************************")
	fmt.Printf("\n\n")
}

// * The number of requests
// * The fastest time
// * The slowest time
// * The mean & median times
// * The percentage of requests that succeeded
// * Any error codes returned that weren't a success
// * The size in bytes of the smallest response
// * The size in bytes of the largest response
func calculateData(concurrent, processingTime, requestTime, maxTime, minTime, successNum, failureNum, maxSize, minSize uint64, chanSize int, medianTime float64, statusCode map[int]int, errorCode map[int]int) {
	if processingTime == 0 {
		processingTime = 1
	}

	var (
		averageTime      float64
		maxTimeFloat     float64
		minTimeFloat     float64
		requestTimeFloat float64
	)
	total := successNum + failureNum
	if total != 0 && concurrent != 0 {
		averageTime = float64(processingTime) / float64(total*1e6)
	}

	maxTimeFloat = float64(maxTime) / 1e6
	minTimeFloat = float64(minTime) / 1e6
	medianTime = medianTime / 1e6
	requestTimeFloat = float64(requestTime) / 1e9
	table(successNum, failureNum, maxSize, minSize, averageTime, medianTime, maxTimeFloat, minTimeFloat, requestTimeFloat, chanSize, statusCode, errorCode)
}

func header() {
	fmt.Printf("\n\n")
	fmt.Println("──────┬────────────┬───────┬───────────┬───────────┬──────────────┬───────────────┬──────────────┬──────────────┬─────────────────────┬────────────")
	fmt.Println(" Time │ Concurrent │ Total |  Success  │  Failure  │ Largest Size │ Smallest Size │ Fastest Time │ Slowest Time │ Average/Median Time │ Status/Error Code")
	fmt.Println("──────┼────────────┼───────┼───────────┼───────────┼──────────────┼───────────────┼──────────────┼──────────────┼─────────────────────┼────────────")
}
func table(successNum, failureNum, maxSize, minSize uint64,
	averageTime, medianTime, maxTimeFloat, minTimeFloat, requestTimeFloat float64, chanSize int, statusCode map[int]int, errorCode map[int]int) {
	total := successNum + failureNum
	result := fmt.Sprintf("%5.0fs│%12d│%7d|%4d (%3.f%%)│%4d (%3.f%%)│%14d|%15d|%14.0f|%14.0f|%10.0f/%-10.0f|%v / %v",
		requestTimeFloat, chanSize, total, successNum, float64(successNum)/float64(total)*100, failureNum, float64(failureNum)/float64(total)*100,
		maxSize, minSize, minTimeFloat, maxTimeFloat, averageTime, medianTime, printMap(statusCode), printMap(errorCode))
	fmt.Println(result)
}

func printMap(statusCode map[int]int) (mapStr string) {

	var mapArr []string
	for key, value := range statusCode {
		mapArr = append(mapArr, fmt.Sprintf("%d:%d", key, value))
	}

	sort.Strings(mapArr)

	mapStr = strings.Join(mapArr, ";")

	if mapStr == "" {
		mapStr = "null"
	}

	return
}
