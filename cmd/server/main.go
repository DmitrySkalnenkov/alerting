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
	//    добавьте поддержку аргумента через флаг k=<КЛЮЧ>;
	var hostPortStr string
	var isRestoreBool bool
	var storeIntervalStr string
	var storeFilePathStr string
	var serverKeyValue string
	flag.StringVar(&hostPortStr, "a", "127.0.0.1:8080", "Value for -a (ADDRESS) should be in 'ip:port' format, example: 127.0.0.1:8080")
	flag.BoolVar(&isRestoreBool, "r", true, "Value for -r (RESTORE)  should be 'true' of 'false'")
	flag.StringVar(&storeIntervalStr, "i", "300s", "Value for -i (STORE_INTERVAL) flag 'r' should be time in second, example: 300")
	//flag.StringVar(&storeFilePathStr, "f", "/tmp/devops-metrics-db.json", "Store file path should be "+
	//	"absolute path to file. If STORE_FILE variable is empty string than storing functionality will not be used.")
	flag.StringVar(&storeFilePathStr, "f", "", "Store file path should be "+
		"absolute path to file. If STORE_FILE variable is empty string than storing functionality will not be used.")
	flag.StringVar(&serverKeyValue, "k", "", "Server key value for HMAC-SHA-256 calculation. Should be hexstring. Example: 300")
	flag.Parse()

	//   ADDRESS (по умолчанию: "127.0.0.1:8080" или "localhost:8080")
	//   STORE_INTERVAL (по умолчанию 300) - интервал времени в секундах, по истечении которого текущие показания сервера сбрасываются на диск. (значение 0 - делает запись синхронной)
	//   STORE_FILE по умолчанию ("/tmp/devops-metrics-db.json") - строка - имя файла, где хранятся значения (пустое значение - отключает функцию записи на диск)
	//   RESTORE по умолчанию (true) - булево значение (true|false), определяющее загружать или нет начальные значения из указанного файла при старте сервера.
	//   добавьте поддержку аргумента через переменную окружения KEY=<КЛЮЧ>;
	envHostPortStr, isEnvHostPort := os.LookupEnv("ADDRESS")
	envStoreIntervalStr, isEnvStoreInterval := os.LookupEnv("STORE_INTERVAL")
	envRestoreStr, isEnvRestore := os.LookupEnv("RESTORE")
	envStoreFilePath, isEnvStoreFilePath := os.LookupEnv("STORE_FILE")
	envServerKeyValue, isServerKeyValue := os.LookupEnv("KEY")

	if isEnvHostPort && envHostPortStr != "" {
		hostPortStr = envHostPortStr
	}
	if isEnvStoreInterval && envStoreIntervalStr != "" {
		storeIntervalStr = envStoreIntervalStr
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
	if isServerKeyValue && envServerKeyValue != "" {
		serverKeyValue = envServerKeyValue
	}

	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Use(middleware.Compress(5))
	//r.HandleFunc("/", handlers.GetAllMetricsHandler)
	r.Get("/", handlers.GetAllMetricsHandlerAPI15)
	r.Post("/", handlers.GetAllMetricsHandlerAPI2)
	r.Post("/update/", handlers.UpdateHandler)
	r.Post("/value/", handlers.ValueHandler)
	//r.Get("/update/gauge/*", handlers.GaugeHandlerAPI1)
	//r.Get("/update/counter/*", handlers.CounterHandlerAPI1)
	r.Post("/update/gauge/*", handlers.GaugeHandlerAPI1)
	r.Post("/update/counter/*", handlers.CounterHandlerAPI1)
	r.Post("/update/*", handlers.NotImplementedHandler)
	r.Post("/value/gauge/{MetricName}", handlers.GetGaugeHandlerAPI1)
	r.Post("/value/counter/{MetricName}", handlers.GetCounterHandlerAPI1)
	r.Get("/value/gauge/{MetricName}", handlers.GetGaugeHandlerAPI1)
	r.Get("/value/counter/{MetricName}", handlers.GetCounterHandlerAPI1)

	//hostPort = auxiliary.TrimQuotes(hostPort)
	//storeIntervalStr += "s"

	storeIntervalTime := 0 * time.Second
	if storeIntervalStr != "0" {
		var err error
		storeIntervalTime, err = time.ParseDuration(storeIntervalStr)
		if err != nil {
			fmt.Printf("ERROR[S]: Cannot conver STORE_INTERVAL value (%s) to time. Will be used 0 value. \n",
				storeIntervalStr)
			storeIntervalTime = 0 * time.Second
			return
		}
	}

	s := &http.Server{
		Addr:         hostPortStr,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 5 * time.Second,
		IdleTimeout:  30 * time.Second,
	}
	s.Handler = r
	var c = make(chan storage.MetricsStorage)
	if isRestoreBool {
		storage.RestoreMetricsFromFile(storeFilePathStr, storage.ServerMetStorage)
	}
	//_ = storeIntervalTime
	if storeFilePathStr != "" {
		go storage.UpdateMetricsInChannel(c)
		go storage.WriteMetricsToFile(storeFilePathStr, c, storeIntervalTime)
	}
	log.Fatal(s.ListenAndServe())
}
