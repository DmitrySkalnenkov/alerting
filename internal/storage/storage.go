package storage

type MemStorage struct {
	Gauges   map[string]float64
	Counters map[string]int64
}

func (m MemStorage) PushGauge(metricName string, value float64) {
	m.Gauges[metricName] = value
}

/*func (m MemStorage) PopGauge(metricName string) float64 {
	_, ok := m.gauges[metricName]
	if ok {
		return m.gauges[metricName]
	} else {
		fmt.Printf("Gauge metric with name  %s is not found.\n", metricName)
		return 0
	}
}*/

func (m MemStorage) PushCounter(metricName string, value int64) {
	_, ok := m.Counters[metricName]
	if ok {
		m.Counters[metricName] = m.Counters[metricName] + value
	} else {
		m.Counters[metricName] = value
	}
}

/*func (m MemStorage) PopCounter(metricName string) int64 {
	_, ok := m.counters[metricName]
	if ok {
		return m.counters[metricName]
	} else {
		fmt.Printf("Gauge metric with name  %s is not found.\n", metricName)
		return 0
	}
}*/
