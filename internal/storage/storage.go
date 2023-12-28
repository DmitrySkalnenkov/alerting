package storage

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"

	"fmt"
	"io"
	"log"
	"os"
	"strconv"
	"time"
)

// var KeyHexStr = "0102030405060708090a0b0c0d0e0f10111213141516171819"
var KeyHexStr = ""
var ServerKeyHexStr = ""
var AgentKeyHexStr = ""

type Metrics struct {
	ID    string   `json:"id"`              // имя метрики
	MType string   `json:"type"`            // параметр, принимающий значение gauge или counter
	Delta *int64   `json:"delta,omitempty"` // значение метрики в случае передачи counter
	Value *float64 `json:"value,omitempty"` // значение метрики в случае передачи gauge
	Hash  string   `json:"hash,omitempty"`  // значение хеш-функции
}

type MetricsStorage []Metrics

var MetStorage = NewMetricStorage()

var NilMetric = Metrics{
	ID:    "",
	MType: "",
	Delta: nil,
	Value: nil,
	Hash:  "",
}

func NewMetricStorage() *MetricsStorage {
	var M = new(MetricsStorage)
	return M
}

func NewMetric() *Metrics {
	N := new(Metrics)
	N.ID = ""
	N.MType = ""
	N.Value = new(float64)
	N.Delta = new(int64)
	N.Hash = ""
	return N
}

func MakeMetric(id string, mType string, mData string) Metrics {
	var N Metrics
	var dataStr string
	N.ID = id
	switch mType {
	case "gauge":
		N.MType = mType
		v, err := strconv.ParseFloat(mData, 64)
		if err != nil {
			fmt.Printf("ERROR: Cannot convert data value to float. Will be used nil metric.")
			return NilMetric
		}
		N.Value = PointOf(v)
		dataStr = fmt.Sprintf("%s:gauge:%f", N.ID, *N.Value)
		N.Hash = HmacSha256(StringToHexStr(dataStr), AgentKeyHexStr)
		return N
	case "counter":
		N.MType = mType
		d, err := strconv.ParseInt(mData, 10, 64)
		if err != nil {
			fmt.Printf("ERROR: Cannot convert data value to int. Will be used nil metric.")
			return NilMetric
		}
		N.Delta = PointOf(d)
		dataStr = fmt.Sprintf("%s:counter:%d", N.ID, *N.Delta)
		N.Hash = HmacSha256(StringToHexStr(dataStr), AgentKeyHexStr)
		return N
	default:
		fmt.Printf("ERROR: Wrong metric type value. Must be 'counter' or 'gague'. Will be used nil metric\n")
		return NilMetric
	}
}

// StringToHexStr() converts ACCII symbol string to HEX string
func StringToHexStr(dataStr string) string {
	return hex.EncodeToString([]byte(dataStr))
}

func isHexString(s string) bool {
	_, err := hex.DecodeString(s)
	return err == nil
}
func SetEncKey(keyVal string) string {
	if isHexString(keyVal) {
		return keyVal
	} else {
		fmt.Printf("WARN: Key value (%s) can't convert into hex. Hash calculation will be disabled.", keyVal)
		return ""
	}
}

// HmacSha256() -- function HMAC-SHA256
func HmacSha256(dataHexStr string, keyHexStr string) string {
	var keyBin, dataBin []byte
	var err error
	keyBin, err = hex.DecodeString(keyHexStr)
	if err != nil {
		fmt.Printf("ERROR: Cannot convert key {%s} into hex string.\n", keyHexStr)
		return ""
	}
	dataBin, err = hex.DecodeString(dataHexStr)
	if err != nil {
		fmt.Printf("ERROR: Cannot convert data {%s} into hex string.\n", dataHexStr)
		return ""
	}
	hmac256 := hmac.New(sha256.New, keyBin)
	hmac256.Write(dataBin)
	dataHmac256 := hmac256.Sum(nil)
	hmac256Hex := hex.EncodeToString(dataHmac256)
	return hmac256Hex
}

// SetMetric -- Metric setter
func (pm *MetricsStorage) SetMetric(m Metrics) {
	if *pm != nil {
		for i := 0; i < len(*pm); i++ {
			if (*pm)[i].ID == m.ID && (*pm)[i].MType == m.MType {
				switch m.MType {
				case "gauge":
					(*pm)[i].Value = m.Value
					(*pm)[i].Delta = new(int64)
					return
				case "counter":
					*(*pm)[i].Delta = *(*pm)[i].Delta + *m.Delta
					(*pm)[i].Value = new(float64)
					return
				}
			}
		}
		if m.ID != "" && (m.MType == "gauge" || m.MType == "counter") {
			*pm = append(*pm, m)
			return
		}
	} else if m.ID != "" && (m.MType == "gauge" || m.MType == "counter") {
		*pm = append(*pm, m)
		return
	}
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
func IsMetricsEqual(m1 Metrics, m2 Metrics) (res bool) {
	if m1.ID == m2.ID && m1.MType == m2.MType {
		switch m1.MType {
		case "gauge":
			if m1.Value == nil && m2.Value == nil {
				return true
			}
			if *m1.Value == *m2.Value && m1.Hash == m2.Hash {
				//fmt.Printf("DEBUG: Metric1 value is %v, Metric2 value is %v.\n", *m1.Value, *m2.Value)
				return true
			} else {
				return false
			}
		case "counter":
			if m1.Delta == nil && m2.Delta == nil {
				return true
			}
			if *m1.Delta == *m2.Delta && m1.Hash == m2.Hash {
				//fmt.Printf("DEBUG: Metric1 value is %v, Metric2 value is %v.\n", *m1.Value, *m2.Value)
				return true
			} else {
				return false
			}
		default:
			return false
		}
	} else {
		return false
	}
}

func PointOf[T any](value T) *T {
	return &value
}

// Update metrics values in channel
func UpdateMetricsInChannel(ch chan MetricsStorage) {
	ms := MetStorage
	for i := 0; ; i++ {
		ch <- *ms
		time.Sleep(1 * time.Second)
	}
}

// Restore metrics from file in MetStorage
func RestoreMetricsFromFile(fileStoragePath string, ms *MetricsStorage) {
	if fileStoragePath != "" {
		fileMetricStorage, err := os.OpenFile(fileStoragePath, os.O_RDONLY|os.O_CREATE, 0777)
		if err != nil {
			fmt.Printf("ERROR: Cannot open file '%s'.\n", fileStoragePath)
			log.Fatal(err)
		}
		defer fileMetricStorage.Close()
		fromFile, err := io.ReadAll(fileMetricStorage)
		if err != nil {
			fmt.Printf("ERROR: Cannot read file '%s'.\n", fileStoragePath)
			log.Fatal(err)
		}
		var tmp MetricsStorage
		err = json.Unmarshal(fromFile, &tmp)
		if err == nil {
			fmt.Printf("INFO: Metrics from file were restored succesfully.\n")
			*ms = tmp
		} else {
			fmt.Printf("ERROR: %s.\n", err)
		}
	}
}

// Writing metrics to file metric storage
func WriteMetricsToFile(filePath string, ch chan MetricsStorage, st time.Duration) {
	fileMetricStorage, err := os.OpenFile(filePath, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0777)
	if err != nil {
		fmt.Printf("ERROR: Cannot open file '%s'.\n", filePath)
		log.Fatal(err)
	}
	defer fileMetricStorage.Close()
	for {
		curMetricStorage := <-ch
		//fmt.Printf("DEBUG: Current metric string is '%s'.\n", curMetricStorage)
		if len(curMetricStorage) > 0 {
			toFile, _ := json.Marshal(curMetricStorage)
			fileMetricStorage.Truncate(0)
			fileMetricStorage.Seek(0, 0)
			_, err := fileMetricStorage.WriteString(string(toFile))
			if err == nil {
				fmt.Printf("INFO: Metrics from the server were dumped to the file.\n")
			} else {
				fmt.Printf("ERROR: %s.\n", err)
				return
			}
		}
		if st != 0 {
			time.Sleep(st)
		}
	}
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
