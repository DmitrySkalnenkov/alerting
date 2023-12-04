package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	//"github.com/DmitrySkalnenkov/alerting/internal/auxiliary"
	"github.com/DmitrySkalnenkov/alerting/internal/handlers"
	"github.com/DmitrySkalnenkov/alerting/internal/storage"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func main() {
	//    ADDRESS, через флаг "-a=<ЗНАЧЕНИЕ>"
	//    RESTORE, через флаг "-r=<ЗНАЧЕНИЕ>"
	//    STORE_INTERVAL, через флаг "-i=<ЗНАЧЕНИЕ>"
	//    STORE_FILE, через флаг "-f=<ЗНАЧЕНИЕ>"
	var hostPort string
	var isRestoreBool bool
	var storeIntervalStr string
	var storeFilePathStr string
	flag.StringVar(&hostPort, "a", "localhost:8080", "ADDRESS should be in 'ip:port' format")
	flag.BoolVar(&isRestoreBool, "r", true, "RESTORE variable or flag 'r' should be 'true' of 'false'")
	flag.StringVar(&storeIntervalStr, "i", "", "RESTORE variable or flag 'r' should be 'true' of 'false'.")
	flag.StringVar(&storeFilePathStr, "f", "/tmp/devops-metrics-db.json", "Store file path should be "+
		"absolute path to file. If STORE_FILE variable is empty string than storing functionality will not be used.")
	flag.Parse()

	//   ADDRESS (по умолчанию: "127.0.0.1:8080" или "localhost:8080")
	//   STORE_INTERVAL (по умолчанию 300) - интервал времени в секундах, по истечении которого текущие показания сервера сбрасываются на диск. (значение 0 - делает запись синхронной)
	//   STORE_FILE по умолчанию ("/tmp/devops-metrics-db.json") - строка - имя файла, где хранятся значения (пустое значение - отключает функцию записи на диск)
	//   RESTORE по умолчанию (true) - булево значение (true|false), определяющее загружать или нет начальные значения из указанного файла при старте сервера.
	envHostPortStr, isEnvHostPort := os.LookupEnv("ADDRESS")
	envStoreIntervalStr, isEnvStoreInterval := os.LookupEnv("STORE_INTERVAL")
	envRestoreStr, isEnvRestore := os.LookupEnv("RESTORE")
	envStoreFilePath, isEnvStoreFilePath := os.LookupEnv("STORE_FILE")
	if isEnvHostPort && envHostPortStr != "" {
		hostPort = envHostPortStr
	}
	if isEnvStoreInterval && envStoreIntervalStr != "" {
		storeIntervalStr = envStoreIntervalStr
	}
	if isEnvRestore && envRestoreStr != "" {
		if envRestoreStr == "false" {
			isRestoreBool = false
		} else {
			isRestoreBool = true
		}
	}
	if isEnvStoreFilePath && envStoreIntervalStr != "" {
		storeFilePathStr = envStoreFilePath
	}

	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Use(middleware.Compress(5))
	r.HandleFunc("/", handlers.GetAllMetricsHandler)
	r.Post("/update/", handlers.UpdateHandler)
	r.Post("/value/", handlers.ValueHandler)
	r.Get("/update/gauge/*", handlers.GaugeHandlerAPI1)
	r.Get("/update/counter/*", handlers.CounterHandlerAPI1)
	r.Post("/update/gauge/*", handlers.GaugeHandlerAPI1)
	r.Post("/update/counter/*", handlers.CounterHandlerAPI1)
	r.Post("/update/*", handlers.NotImplementedHandler)
	r.Post("/value/gauge/{MetricName}", handlers.GetGaugeHandlerAPI1)
	r.Post("/value/counter/{MetricName}", handlers.GetCounterHandlerAPI1)
	r.Get("/value/gauge/{MetricName}", handlers.GetGaugeHandlerAPI1)
	r.Get("/value/counter/{MetricName}", handlers.GetCounterHandlerAPI1)

	//hostPort := auxiliary.GetEnvVariable("ADDRESS", "localhost:8080")
	//hostPort := auxiliary.GetParamValue("ADDRESS", "a", "localhost:8080", "ADDRESS should be in 'ip:port' format")
	//hostPort = auxiliary.TrimQuotes(hostPort)

	//storeIntervalStr := auxiliary.GetEnvVariable("STORE_INTERVAL", "300s")
	//storeIntervalStr += "s"

	//storeFilePath := ""

	//isRestoreStr := auxiliary.GetParamValue("RESTORE", "r", "true", "RESTORE variable or flag 'r' should be 'true' of 'false'.")

	/*if !(isStoreFilePath && valueStoreFilePath == "") {
		//storeFilePath = auxiliary.GetEnvVariable("STORE_FILE", "/tmp/devops-metrics-db.json")
		storeFilePath = auxiliary.GetParamValue("STORE_FILE", "f", "/tmp/devops-metrics-db.json", "Store file path should be absolute path "+
			"to file. If STORE_FILE variable is empty string than storing functionality will not be used.")
	}*/
	//isRestoreStr := auxiliary.GetEnvVariable("RESTORE", "true")

	storeIntervalTime := 0 * time.Second
	if storeIntervalStr != "0" {
		var err error
		storeIntervalTime, err = time.ParseDuration(storeIntervalStr)
		if err != nil {
			fmt.Printf("ERROR: Cannot conver STORE_INTERVAL value (%s) to time. Will be used 0 value. \n",
				storeIntervalStr)
			storeIntervalTime = 0 * time.Second
			return
		}
	}

	s := &http.Server{
		Addr:         hostPort,
		ReadTimeout:  30 * time.Second,
		WriteTimeout: 30 * time.Second,
	}
	s.Handler = r
	//s.Handler = handlers.GzipHandle(r)

	var c = make(chan storage.MetricsStorage)
	if isRestoreBool {
		storage.RestoreMetricsFromFile(storeFilePathStr, storage.MetStorage)
	}
	_ = storeIntervalTime
	if storeFilePathStr != "" { // Disabling of storing metrics into the file
		go storage.UpdateMetricsInChannel(c)
		go storage.WriteMetricsToFile(storeFilePathStr, c, storeIntervalTime)
	}
	log.Fatal(s.ListenAndServe())
}
