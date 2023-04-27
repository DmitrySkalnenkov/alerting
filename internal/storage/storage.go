package storage

import (
	"fmt"
)

type Metrics struct {
	ID    string   `json:"id"`              // имя метрики
	MType string   `json:"type"`            // параметр, принимающий значение gauge или counter
	Delta *int64   `json:"delta,omitempty"` // значение метрики в случае передачи counter
	Value *float64 `json:"value,omitempty"` // значение метрики в случае передачи gauge
}

type MetricsStorage []Metrics

func NewMetricStorage() *MetricsStorage {
	var M = new(MetricsStorage)
	return M
}

var MetStorage = NewMetricStorage()

var NilMetric = Metrics{
	ID:    "",
	MType: "",
	Delta: nil,
	Value: nil,
}

// SetMetric -- Metric setter
func (pm *MetricsStorage) SetMetric(m Metrics) {
	for i := 0; i < len(*pm); i++ {
		if (*pm)[i].ID == m.ID {
			switch m.MType {
			case `gauge`:
				(*pm)[i].MType = "gauge"
				(*pm)[i].Value = m.Value
				(*pm)[i].Delta = nil
			case `counter`:
				(*pm)[i].MType = "counter"
				(*pm)[i].Delta = m.Delta
				(*pm)[i].Value = nil
			}
			return
		}
	}
	*pm = append(*pm, m)
	fmt.Printf("DEBUG: MetricStorage is %v \n", (*pm))
}

// GetMetric -- metric getter, if no metric return nilMetric
func (pm *MetricsStorage) GetMetric(metricID string, metricType string) Metrics {
	for i := 0; i < len(*pm); i++ {
		if (*pm)[i].ID == metricID && (*pm)[i].MType == metricType {
			return (*pm)[i]
		}
	}
	fmt.Printf("DEBUG: MetricName %v with type %v not found.\n", metricID, metricType)
	return NilMetric
}

// Comparing metric if them equal then true
func isMetricsEqual(m1 Metrics, m2 Metrics) (res bool) {
	if m1.ID == m2.ID && m1.MType == m2.MType {
		if m1.Value != nil && m2.Value != nil {
			if *m1.Value == *m2.Value {
				fmt.Printf("DEBUG: Metric1 value is %v, Metric2 value is %v.\n", *m1.Value, *m2.Value)
				return true
			}
		} else if m1.Delta != nil && m2.Delta != nil {
			if *m1.Delta == *m2.Delta {
				fmt.Printf("DEBUG: Metric1 delta is %v, Metric2 delta is %v.\n", *m1.Delta, *m2.Delta)
				return true
			}
		}
		return false
	} else {
		return false
	}
}

func PointOf[T any](value T) *T {
	return &value
}

//////////legacy//////////

type MemStorage struct {
	Gauges   map[string]float64
	Counters map[string]int64
}

func NewMemStorage() *MemStorage {
	ms := &MemStorage{
		Gauges:   make(map[string]float64),
		Counters: make(map[string]int64),
	}
	return ms
}

var Mstorage = NewMemStorage()

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
