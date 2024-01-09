package main

import (
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"math/rand"
	"net"
	"net/http"
	"os"
	"runtime"
	"strconv"
	"time"

	"github.com/DmitrySkalnenkov/alerting/internal/auxiliary"
	"github.com/DmitrySkalnenkov/alerting/internal/storage"
)

type Client struct {
	IP     string
	Port   string
	Client *http.Client
}

//Sends agent metric storage to server by POST with Content-Type 'plain/text'
func (cl Client) metricStoragePlainSending(mA *storage.MetricsStorage) {
	curURL := ""
	for row := 0; row < len(*mA); row++ {
		if (*mA)[row].ID != "" {
			switch (*mA)[row].MType {
			case "gauge":
				pv := (*mA)[row].Value
				sv := strconv.FormatFloat(*pv, 'f', 0, 64)
				curURL = fmt.Sprintf("http://%s:%s/update/%s/%s/%s", cl.IP, cl.Port, (*mA)[row].MType, (*mA)[row].ID, sv)
				fmt.Printf("INFO[A]: metricStoragePlainSending(), URL is: %s \n", curURL)
				_, err := cl.sendPlainPostRequest(curURL)
				if err != nil {
					fmt.Printf("ERROR[A]: sendPostRequest() error -- %v. \n", err)
				}
			case "counter":
				pd := (*mA)[row].Delta
				sd := strconv.FormatInt(*pd, 10)
				curURL = fmt.Sprintf("http://%s:%s/update/%s/%s/%s", cl.IP, cl.Port, (*mA)[row].MType, (*mA)[row].ID, sd)
				fmt.Printf("INFO[A]: metricStoragePlainSending(), URL is: %s \n", curURL)
				_, err := cl.sendPlainPostRequest(curURL)
				if err != nil {
					fmt.Printf("ERROR[A]: sendPostRequest() error -- %v.\n", err)
				}
			default:
				fmt.Printf("ERROR[A]: metricStoragePlainSending() wrong metric type. It must be `gauge` or `counter`.\n")
			}
		}
	}
}

// Sends  metrics to server by POST ("application/json") with metric type and value in JSON.
func (cl Client) metricStorageJsonSending(mA *storage.MetricsStorage) {
	curURL := ""
	var curMetric storage.Metric
	for row := 0; row < len(*mA); row++ {
		if (*mA)[row].ID != "" {
			curURL = fmt.Sprintf("http://%s:%s/update/", cl.IP, cl.Port)
			curMetric = (*mA)[row]
			fmt.Printf("DEBUG[A]: For sending. curMetric.ID = %v, curMetric.MType = %v, curMetric.Value = %v, curMetric.Delta = %v.\n",
				curMetric.ID, curMetric.MType, curMetric.Value, curMetric.Delta)
			_, err := cl.sendJsonPostRequest(curURL, curMetric)
			if err != nil {
				fmt.Printf("ERROR[A]: %v.\n", err)
			}
		}
	}
}

// Sends POST request with content type "plain/text"
func (cl Client) sendPlainPostRequest(curURL string) (string, error) {
	request, err := http.NewRequest(http.MethodPost, curURL, nil)
	if err != nil {
		fmt.Printf("ERROR[A]: Error value is %v.\n", err)
		return "", err
	}
	request.Header.Set("Content-Type", "plain/text")
	response, err := cl.Client.Do(request)
	if err != nil {
		fmt.Printf("ERROR[A]: Error value is  %v. Response is  %v \n", err, response)
		return "", err
	}
	defer response.Body.Close()
	fmt.Printf("DEBUG[A]: Response status code: %s.\n", response.Status)
	return string(response.Status), nil
}

// Sends request by POST method with content type "application/json" and metric data in JSON
func (cl Client) sendJsonPostRequest(curURL string, m storage.Metric) (string, error) {
	payloadBuf := new(bytes.Buffer)
	err := json.NewEncoder(payloadBuf).Encode(m)
	if err != nil {
		fmt.Printf("ERROR[A]: Error value is  %v.\n", err)
		return "", err
	}
	request, err := http.NewRequest(http.MethodPost, curURL, payloadBuf)
	fmt.Printf("DEBUG[A]: Request is %v.\n", request)
	if err != nil {
		fmt.Printf("ERROR[A]: %s.\n", err)
		return "", err
	}
	request.Header.Set("Content-Type", "application/json")
	response, err := cl.Client.Do(request)
	if err != nil {
		fmt.Printf("ERROR[A]: Error value is  %v. Response is  %v \n", err, response)
		return "", err
	}
	defer response.Body.Close()
	fmt.Printf("DEBUG[A]: Response status code: %s.\n", response.Status)
	return string(response.Status), nil
}

func getMetricsStorage(mA *storage.MetricsStorage, PollCount *int64, rtm *runtime.MemStats) {
	runtime.ReadMemStats(rtm)
	*PollCount = *PollCount + 1
	RandomValue := float64(rand.Float64())
	(*mA)[0] = storage.MakeMetric("Alloc", "gauge", strconv.FormatUint(rtm.Alloc, 10))
	(*mA)[1] = storage.MakeMetric("BuckHashSys", "gauge", strconv.FormatUint(rtm.BuckHashSys, 10))
	(*mA)[2] = storage.MakeMetric("Frees", "gauge", strconv.FormatUint(rtm.Frees, 10))
	(*mA)[3] = storage.MakeMetric("GCCPUFraction", "gauge", strconv.FormatFloat(rtm.GCCPUFraction, 'G', -1, 64))
	(*mA)[4] = storage.MakeMetric("GCSys", "gauge", strconv.FormatUint(rtm.GCSys, 10))
	(*mA)[5] = storage.MakeMetric("HeapAlloc", "gauge", strconv.FormatUint(rtm.HeapAlloc, 10))
	(*mA)[6] = storage.MakeMetric("HeapIdle", "gauge", strconv.FormatUint(rtm.HeapIdle, 10))
	(*mA)[7] = storage.MakeMetric("HeapInuse", "gauge", strconv.FormatUint(rtm.HeapInuse, 10))
	(*mA)[8] = storage.MakeMetric("HeapObjects", "gauge", strconv.FormatUint(rtm.HeapObjects, 10))
	(*mA)[9] = storage.MakeMetric("HeapReleased", "gauge", strconv.FormatUint(rtm.HeapReleased, 10))
	(*mA)[10] = storage.MakeMetric("HeapSys", "gauge", strconv.FormatUint(rtm.HeapSys, 10))
	(*mA)[11] = storage.MakeMetric("LastGC", "gauge", strconv.FormatUint(rtm.LastGC, 10))
	(*mA)[12] = storage.MakeMetric("Lookups", "gauge", strconv.FormatUint(rtm.Lookups, 10))
	(*mA)[13] = storage.MakeMetric("MCacheInuse", "gauge", strconv.FormatUint(rtm.MCacheInuse, 10))
	(*mA)[14] = storage.MakeMetric("MCacheSys", "gauge", strconv.FormatUint(rtm.MCacheSys, 10))
	(*mA)[15] = storage.MakeMetric("MSpanInuse", "gauge", strconv.FormatUint(rtm.MSpanInuse, 10))
	(*mA)[16] = storage.MakeMetric("MSpanSys", "gauge", strconv.FormatUint(rtm.MSpanSys, 10))
	(*mA)[17] = storage.MakeMetric("Mallocs", "gauge", strconv.FormatUint(rtm.Mallocs, 10))
	(*mA)[18] = storage.MakeMetric("NextGC", "gauge", strconv.FormatUint(rtm.NextGC, 10))
	(*mA)[19] = storage.MakeMetric("NumForcedGC", "gauge", strconv.FormatUint(uint64(rtm.NumForcedGC), 10))
	(*mA)[20] = storage.MakeMetric("NumGC", "gauge", strconv.FormatUint(uint64(rtm.NumGC), 10))
	(*mA)[21] = storage.MakeMetric("OtherSys", "gauge", strconv.FormatUint(rtm.OtherSys, 10))
	(*mA)[22] = storage.MakeMetric("PollCount", "counter", strconv.FormatInt(int64(*PollCount), 10))
	(*mA)[23] = storage.MakeMetric("PauseTotalNs", "gauge", strconv.FormatUint(rtm.PauseTotalNs, 10))
	(*mA)[24] = storage.MakeMetric("RandomValue", "gauge", strconv.FormatFloat(RandomValue, 'G', -1, 64))
	(*mA)[25] = storage.MakeMetric("StackInuse", "gauge", strconv.FormatUint(rtm.StackInuse, 10))
	(*mA)[26] = storage.MakeMetric("StackSys", "gauge", strconv.FormatUint(rtm.StackSys, 10))
	(*mA)[27] = storage.MakeMetric("Sys", "gauge", strconv.FormatUint(rtm.Sys, 10))
	(*mA)[28] = storage.MakeMetric("TotalAlloc", "gauge", strconv.FormatUint(rtm.TotalAlloc, 10))
}

func main() {
	StartTime := time.Now()
	fmt.Printf("Start time: %s.\n", string(StartTime.String()))
	var CurTime time.Time
	LastPoolTime := time.Now()
	LastReportTime := time.Now()

	var hostPortStr string = ""
	var reportIntervalStr string = ""
	var pollIntervalStr string = ""
	var keyValueStr string = ""
	flag.StringVar(&hostPortStr, "a", "127.0.0.1:8080", "Value for -a (ADDRESS) should be in 'ip:port' format, example: 127.0.0.1:8080")                                                      //(i7)  ADDRESS, через флаг: "-a=<ЗНАЧЕНИЕ>"
	flag.StringVar(&reportIntervalStr, "r", "10", "Value for -r (REPORT_INTERVAL) flag 'r' should be time in second, example: 10")                                                            //(i7)  REPORT_INTERVAL, через флаг: "-r=<ЗНАЧЕНИЕ>"
	flag.StringVar(&pollIntervalStr, "p", "2", "Value for -p (POLL_INTERVAL) flag 'p' should be time in second, example: 2")                                                                  //(i7)  POLL_INTERVAL, через флаг: "-p=<ЗНАЧЕНИЕ>"
	flag.StringVar(&keyValueStr, "k", "", "Key value for HMAC-SHA-256 calculation of hash. Should be hexstring, example: 'dce8b88a0e5943ab3431c6e41293e1e33790162f09020704342b064a92d651d5'") //(i9)  добавьте поддержку аргумента через флаг k=<КЛЮЧ>;
	flag.Parse()

	envHostPortStr, isEnvHostPort := os.LookupEnv("ADDRESS")                     //(i5) ADDRESS (по умолчанию: "127.0.0.1:8080" или "localhost:8080")
	envReportIntervalStr, isEnvReportInterval := os.LookupEnv("REPORT_INTERVAL") //(i5) REPORT_INTERVAL (по умолчанию: 10 секунд)
	envPollIntervalStr, isEnvPollInterval := os.LookupEnv("POLL_INTERVAL")       //(i5) POLL_INTERVAL (по умолчанию: 2 секунды)
	envKeyValueStr, isKeyValue := os.LookupEnv("KEY")                            //(i9) добавьте поддержку аргумента через переменную окружения KEY=<КЛЮЧ>;

	if isEnvHostPort && envHostPortStr != "" {
		hostPortStr = envHostPortStr
	}
	if isEnvReportInterval && envReportIntervalStr != "" {
		reportIntervalStr = envReportIntervalStr
	}
	if isEnvPollInterval && envPollIntervalStr != "" {
		reportIntervalStr = envReportIntervalStr
	}
	hostPortStr = auxiliary.TrimQuotes(hostPortStr)
	serverIPAddress, serverTCPPort, err := net.SplitHostPort(hostPortStr)
	if err != nil {
		fmt.Printf("WARN[A]: Cannot get IP and PORT value from ADDRESS string (%s). Will be used default values (127.0.0.1:8080).\n", hostPortStr)
		serverIPAddress = "127.0.0.1"
		serverTCPPort = "8080"
	}
	if isKeyValue && envKeyValueStr != "" {
		keyValueStr = envKeyValueStr
	}

	var pollInterval time.Duration
	pollValue, err := strconv.Atoi(pollIntervalStr)
	if err == nil {
		pollInterval = time.Duration(pollValue) * time.Second
	}
	var reportInterval time.Duration
	reportValue, err := strconv.Atoi(reportIntervalStr)
	if err == nil {
		reportInterval = time.Duration(reportValue) * time.Second
	}
	fmt.Printf("DEBUG[A]: PollInterval is %s.\n", pollInterval)
	fmt.Printf("DEBUG[A]: ReportInterval is %s.\n", reportInterval)
	baseURL := fmt.Sprintf("http://%s:%s", serverIPAddress, serverTCPPort)
	fmt.Printf("DEBUG[A]: BaseURL is %s.\n", baseURL)

	METRIC_AMOUNT := 29
	var PollCount int64 = 0
	var rtm runtime.MemStats
	var agentMetricStorage storage.MetricsStorage
	agentMetricStorage = make(storage.MetricsStorage, METRIC_AMOUNT) //agentMetricStorage init
	for i := 0; i < METRIC_AMOUNT; i++ {
		agentMetricStorage = append(agentMetricStorage, storage.NilMetric)
	}
	var cl Client
	cl.IP = serverIPAddress
	cl.Port = serverTCPPort
	cl.Client = &http.Client{}
	cl.Client.Timeout = 1 * time.Second
	transport := &http.Transport{}
	transport.MaxIdleConns = 20
	transport.IdleConnTimeout = 5 * time.Second
	cl.Client.Transport = transport
	for {
		CurTime = time.Now()
		if CurTime.Sub(LastPoolTime) > pollInterval {
			fmt.Printf("INFO[A]: PoolTime: %s.\n", string(LastPoolTime.String()))
			getMetricsStorage(&agentMetricStorage, &PollCount, &rtm)
			LastPoolTime = time.Now()
		}
		if CurTime.Sub(LastReportTime) > reportInterval {
			fmt.Printf("INFO[A]: ReportTime: %s.\n", string(LastReportTime.String()))
			//cl.metricStoragePlainSending(&agentMetricStorage)
			cl.metricStorageJsonSending(&agentMetricStorage)
			LastReportTime = time.Now()
		}
	}
}
