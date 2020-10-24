package stat

import (
	"fmt"
	"sort"
	"strings"
	"sync"
	"time"

	"github.com/jameshwc/go-stress/profile"
	"github.com/jameshwc/go-stress/stat/median"
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

	ticker := time.NewTicker(exportPeriod)
	go func() {
		for {
			select {
			case <-ticker.C:
				endTime := uint64(time.Now().UnixNano())
				requestTime = endTime - startTime
				go calculateData(concurrent, processingTime, requestTime, maxTime, minTime, successNum, failureNum, maxSize, minSize, chanSize, medianTime, statusCode)
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

		if _, ok := chanIDs[data.ChanID]; !ok {
			chanIDs[data.ChanID] = true
			chanSize++
		}
		medianTime = medianFinder.FindMedian()
	}

	stopChan <- true

	endTime := uint64(time.Now().UnixNano())
	requestTime = endTime - startTime

	calculateData(concurrent, processingTime, requestTime, maxTime, minTime, successNum, failureNum, maxSize, minSize, chanSize, medianFinder.FindMedian(), statusCode)

	fmt.Printf("\n\n")

	fmt.Println("*************************  Statistics  ****************************")
	fmt.Println("Total Concurrent Number:", concurrent)
	fmt.Println("Total Requests: -c * -n）:", successNum+failureNum, "Total Request Time: ", fmt.Sprintf("%.3f", float64(requestTime)/1e9),
		"s", "# of success:", successNum, "# of failure:", failureNum)
	fmt.Println("*******************************************************************")
	fmt.Printf("\n\n")
}

func calculateData(concurrent, processingTime, requestTime, maxTime, minTime, successNum, failureNum, maxSize, minSize uint64, chanSize int, medianTime float64, statusCode map[int]int) {
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
	table(successNum, failureNum, maxSize, minSize, averageTime, medianTime, maxTimeFloat, minTimeFloat, requestTimeFloat, chanSize, statusCode)
}

func header() {
	fmt.Printf("\n\n")
	fmt.Println("──────┬────────────┬───────┬─────────┬─────────┬──────────────┬───────────────┬──────────────┬──────────────┬─────────────────────┬────────────")
	fmt.Println(" Time │ Concurrent │ Total | Success │ Failure │ Largest Size │ Smallest Size │ Fastest Time │ Slowest Time │ Average/Median Time │ Status Code")
	fmt.Println("──────┼────────────┼───────┼─────────┼─────────┼──────────────┼───────────────┼──────────────┼──────────────┼─────────────────────┼────────────")
}
func table(successNum, failureNum, maxSize, minSize uint64,
	averageTime, medianTime, maxTimeFloat, minTimeFloat, requestTimeFloat float64, chanSize int, statusCode map[int]int) {
	result := fmt.Sprintf("%5.0fs│%12d│%7d|%9d│%9d│%14d|%15d|%14.0f|%14.0f|%10.0f/%-10.0f|%v",
		requestTimeFloat, chanSize, successNum+failureNum, successNum, failureNum, maxSize, minSize, minTimeFloat, maxTimeFloat, averageTime, medianTime, printMap(statusCode))
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
