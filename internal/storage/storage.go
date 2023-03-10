package storage

//import (
//	"fmt"
//)

type MemStorage struct {
	Gauges   map[string]float64
	Counters map[string]int64
}

var Mstorage = NewMemStorage()

func NewMemStorage() *MemStorage {
	ms := &MemStorage{
		Gauges:   make(map[string]float64),
		Counters: make(map[string]int64),
	}
	return ms
}

func (m MemStorage) PushGauge(metricName string, value float64) {
	m.Gauges[metricName] = value
}

func (m MemStorage) PopGauge(metricName string) float64 {
	_, ok := m.Gauges[metricName]
	if ok {
		//fmt.Printf("DEBUG: Gauge metric with name %v is found.\n", metricName)
		return m.Gauges[metricName]
	} else {
		//fmt.Printf("DEBUG: Gauge metric with name %v is not found.\n", metricName)
		return 0
	}
}

func (m MemStorage) PushCounter(metricName string, value int64) {
	_, ok := m.Counters[metricName]
	if ok {
		m.Counters[metricName] = m.Counters[metricName] + value
	} else {
		m.Counters[metricName] = value
	}
}

func (m MemStorage) PopCounter(metricName string) int64 {
	_, ok := m.Counters[metricName]
	if ok {
		//fmt.Printf("DEBUG: Counter metric with name %v is found.\n", metricName)
		return m.Counters[metricName]
	} else {
		//fmt.Printf("DEBUG: Counter metric with name  %s is not found.\n", metricName)
		return 0
	}
}
