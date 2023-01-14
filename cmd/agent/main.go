package main

import (
	"fmt"
	"net/http"
	"runtime"
)

type Metrics struct {
	Alloc,
	BuckHashSys,
	Frees,
	GCCPUFraction,
	GCSys,
	HeapAlloc,
	HeapIdle,
	HeapInuse,
	HeapObjects,
	HeapReleased,
	HeapSys,
	LastGC,
	Lookups,
	MCacheInuse,
	MCacheSys,
	MSpanInuse,
	MSpanSys,
	Mallocs,
	NextGC,
	NumForcedGC,
	NumGC,
	OtherSys,
	PauseTotalNs,
	StackInuse,
	StackSys,
	Sys,
	TotalAlloc,
	RandomValue float64
	PollCount int
}

func main() {

	var m Metrics
	var rtm runtime.MemStats
	runtime.ReadMemStats(&rtm)
	m.Alloc = float64(rtm.Alloc)
	fmt.Println(m.Alloc)

	client := http.Client{}
	response, err := client.Get("https://golang.org")
	if err == nil {
		fmt.Println(response)
	} else {
		fmt.Println("Error %s", err)
	}
}
