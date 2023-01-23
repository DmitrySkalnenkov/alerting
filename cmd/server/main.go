package main

import (
	"fmt"
	"io"
	"net/http"
)

type MemStorage struct {
	gaguges  map[string]float64
	counters map[string]int64
}

func (m MemStorage) PushGauge(metricName string, value float64) {
	m.gaguges[metricName] = value
}

func (m MemStorage) PopGauge(metricName string) float64 {
	value, ok := m.gaguges[metricName]
	if ok == true {
		value = m.gaguges[metricName]
		return value
	} else {
		fmt.Println("Metric with name  %s is not found", metricName)
		return 0
	}
}

func (m MemStorage) PushCounter(metricName string, value int64) {
	_, ok := m.counters[metricName]
	if ok {
		m.counters[metricName] = m.counters[metricName] + value
	} else {
		m.counters[metricName] = value
	}
}

func (m MemStorage) PopCounter(metricName string) int64 {
	var value int64
	value = m.counters[metricName]
	return value
}

// HelloWorld — обработчик запроса.
func HelloWorld(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("<h1>Hello, World</h1>"))
}

func main() {

	server := &http.Server{
		Addr: "127.0.0.1:8080",
	}
	server.ListenAndServe()
	//var mstorage MemStorage

	hg := func(w http.ResponseWriter, _ *http.Request) {
		io.WriteString(w, "Hello from a HandleFunc gauge\n")
	}

	hc := func(w http.ResponseWriter, _ *http.Request) {
		io.WriteString(w, "Hello from a HandleFunc counter\n")
	}

	http.HandleFunc("/update/gaguge", hg)
	http.HandleFunc("/update/counter", hc)

}
