package main

import (
	//"encoding/json"
	"fmt"
	"math/rand"
	"net/http"
	"runtime"
	"strconv"
	"time"
)

/*
type Metrics struct {
	Alloc         float64
	BuckHashSys   float64
	Frees         float64
	GCCPUFraction float64
	GCSys         float64
	HeapAlloc     float64
	HeapIdle      float64
	HeapInuse     float64
	HeapObjects   float64
	HeapReleased  float64
	HeapSys       float64
	LastGC        float64
	Lookups       float64
	MCacheInuse   float64
	MCacheSys     float64
	MSpanInuse    float64
	MSpanSys      float64
	Mallocs       float64
	NextGC        float64
	NumForcedGC   float64
	NumGC         float64
	OtherSys      float64
	PollCount     int64
	PauseTotalNs  float64
	RandomValue   float64
	StackInuse    float64
	StackSys      float64
	Sys           float64
	TotalAlloc    float64
}
*/

func main() {
	pollInterval := 2 * time.Second
	//reportInterval := 10 * time.Second

	serverIpAddress := "127.0.0.1"
	serverTcpPort := 8080
	serverUrl := ""
	var PollCount int64
	var RandomValue float64
	var rtm runtime.MemStats

	var mArray [29][3]string

	client := http.Client{}

	for {
		<-time.After(pollInterval)
		runtime.ReadMemStats(&rtm)
		PollCount = PollCount + 1
		RandomValue = float64(rand.Float64())
		//1
		mArray[0][0] = "Alloc"
		mArray[0][1] = "gauge"
		mArray[0][2] = strconv.FormatUint(rtm.Alloc, 10)
		//2
		mArray[1][0] = "BuckHashSys"
		mArray[1][1] = "gauge"
		mArray[1][2] = strconv.FormatUint(rtm.BuckHashSys, 10)
		//3
		mArray[2][0] = "Frees"
		mArray[2][1] = "gauge"
		mArray[2][2] = strconv.FormatUint(rtm.Frees, 10)
		//4
		mArray[3][0] = "GCCPUFraction"
		mArray[3][1] = "gauge"
		mArray[3][2] = strconv.FormatFloat(rtm.GCCPUFraction, 'G', -1, 64)
		//5
		mArray[4][0] = "GCSys"
		mArray[4][1] = "gauge"
		mArray[4][2] = strconv.FormatUint(rtm.GCSys, 10)
		//6
		mArray[5][0] = "HeapAlloc"
		mArray[5][1] = "gauge"
		mArray[5][2] = strconv.FormatUint(rtm.HeapAlloc, 10)
		//7
		mArray[6][0] = "HeapIdle"
		mArray[6][1] = "gauge"
		mArray[6][2] = strconv.FormatUint(rtm.HeapIdle, 10)
		//8
		mArray[7][0] = "HeapInuse"
		mArray[7][1] = "gauge"
		mArray[7][2] = strconv.FormatUint(rtm.HeapInuse, 10)
		//9
		mArray[8][0] = "HeapObjects"
		mArray[8][1] = "gauge"
		mArray[8][2] = strconv.FormatUint(rtm.HeapObjects, 10)
		//10
		mArray[9][0] = "HeapReleased"
		mArray[9][1] = "gauge"
		mArray[9][2] = strconv.FormatUint(rtm.HeapReleased, 10)
		//11
		mArray[10][0] = "HeapSys"
		mArray[10][1] = "gauge"
		mArray[10][2] = strconv.FormatUint(rtm.HeapSys, 10)
		//12
		mArray[11][0] = "LastGC"
		mArray[11][1] = "gauge"
		mArray[11][2] = strconv.FormatUint(rtm.LastGC, 10)
		//13
		mArray[12][0] = "Lookups"
		mArray[12][1] = "gauge"
		mArray[12][2] = strconv.FormatUint(rtm.Lookups, 10)
		//14
		mArray[13][0] = "MCacheInuse"
		mArray[13][1] = "gauge"
		mArray[13][2] = strconv.FormatUint(rtm.MCacheInuse, 10)
		//15
		mArray[14][0] = "MCacheSys"
		mArray[14][1] = "gauge"
		mArray[14][2] = strconv.FormatUint(rtm.MCacheSys, 10)
		//16
		mArray[15][0] = "MSpanInuse"
		mArray[15][1] = "gauge"
		mArray[15][2] = strconv.FormatUint(rtm.MSpanInuse, 10)
		//17
		mArray[16][0] = "MSpanSys"
		mArray[16][1] = "gauge"
		mArray[16][2] = strconv.FormatUint(rtm.MSpanSys, 10)
		//18
		mArray[17][0] = "Mallocs"
		mArray[17][1] = "gauge"
		mArray[17][2] = strconv.FormatUint(rtm.Mallocs, 10)
		//19
		mArray[18][0] = "NextGC"
		mArray[18][1] = "gauge"
		mArray[18][2] = strconv.FormatUint(rtm.NextGC, 10)
		//20
		mArray[19][0] = "NumForcedGC"
		mArray[19][1] = "gauge"
		mArray[19][2] = strconv.FormatUint(uint64(rtm.NumForcedGC), 10)
		//21
		mArray[20][0] = "NumGC"
		mArray[20][1] = "gauge"
		mArray[20][2] = strconv.FormatUint(uint64(rtm.NumGC), 10)
		//22
		mArray[21][0] = "OtherSys"
		mArray[21][1] = "gauge"
		mArray[21][2] = strconv.FormatUint(rtm.OtherSys, 10)
		//23
		mArray[22][0] = "PollCount"
		mArray[22][1] = "counter"
		mArray[22][2] = strconv.FormatInt(PollCount, 10)
		//24
		mArray[23][0] = "PauseTotalNs"
		mArray[23][1] = "gauge"
		mArray[23][2] = strconv.FormatUint(rtm.PauseTotalNs, 10)
		//25
		mArray[24][0] = "RandomValue"
		mArray[24][1] = "gauge"
		mArray[24][2] = strconv.FormatFloat(RandomValue, 'G', -1, 64)
		//26
		mArray[25][0] = "StackInuse"
		mArray[25][1] = "gauge"
		mArray[25][2] = strconv.FormatUint(rtm.StackInuse, 10)
		//27
		mArray[26][0] = "StackSys"
		mArray[26][1] = "gauge"
		mArray[26][2] = strconv.FormatUint(rtm.StackSys, 10)
		//28
		mArray[27][0] = "Sys"
		mArray[27][1] = "gauge"
		mArray[27][2] = strconv.FormatUint(rtm.Sys, 10)
		//29
		mArray[28][0] = "TotalAlloc"
		mArray[28][1] = "gauge"
		mArray[28][2] = strconv.FormatUint(rtm.TotalAlloc, 10)

		fmt.Println()
		fmt.Println(mArray)

		for row := 0; row < len(mArray); row++ {
			if mArray[row][0] != "" {
				serverUrl = fmt.Sprintf("http://%s:%s/update/%s/%s/%s",
					serverIpAddress, strconv.Itoa(serverTcpPort), mArray[row][1],
					mArray[row][0], mArray[row][2])
				fmt.Println(serverUrl)
			}
		}
	}

	response, err := client.Get(serverUrl)
	if err == nil {
		fmt.Println(response)
	} else {
		fmt.Println("Error %s", err)
	}
}
