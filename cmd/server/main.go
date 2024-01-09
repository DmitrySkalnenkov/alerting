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
	var hostPortStr string
	var isRestoreBool bool
	var storeIntervalStr string
	var storeFilePathStr string
	var serverKeyValue string
	flag.StringVar(&hostPortStr, "a", "127.0.0.1:8080", "Value for -a (ADDRESS) should be in 'ip:port' format, example: 127.0.0.1:8080") //(i7) ADDRESS, через флаг "-a=<ЗНАЧЕНИЕ>"
	flag.BoolVar(&isRestoreBool, "r", true, "Value for -r (RESTORE)  should be 'true' of 'false'")                                       //(i7) RESTORE, через флаг "-r=<ЗНАЧЕНИЕ>"
	flag.StringVar(&storeIntervalStr, "i", "300s", "Value for -i (STORE_INTERVAL) flag 'r' should be time in second, example: 300")      //(i7) STORE_INTERVAL, через флаг "-i=<ЗНАЧЕНИЕ>"
	flag.StringVar(&storeFilePathStr, "f", "/tmp/devops-metrics-db.json", "Store file path should be "+                                  //(i7) STORE_FILE, через флаг "-f=<ЗНАЧЕНИЕ>"
		"absolute path to file. If STORE_FILE variable is empty string than storing functionality will not be used.")
	flag.StringVar(&serverKeyValue, "k", "", "Server key value for HMAC-SHA-256 calculation. Should be hexstring. Example: 300") //(i9) добавьте поддержку аргумента через флаг k=<КЛЮЧ>;
	flag.Parse()

	envHostPortStr, isEnvHostPort := os.LookupEnv("ADDRESS")                  //(i5) ADDRESS (по умолчанию: "127.0.0.1:8080" или "localhost:8080")
	envStoreIntervalStr, isEnvStoreInterval := os.LookupEnv("STORE_INTERVAL") //(i6) STORE_INTERVAL (по умолчанию 300) - интервал времени в секундах, по истечении которого текущие показания сервера сбрасываются на диск. (значение 0 - делает запись синхронной)
	envRestoreStr, isEnvRestore := os.LookupEnv("RESTORE")                    //(i6) RESTORE по умолчанию (true) - булево значение (true|false), определяющее загружать или нет начальные значения из указанного файла при старте сервера.
	envStoreFilePath, isEnvStoreFilePath := os.LookupEnv("STORE_FILE")        //(i6) STORE_FILE по умолчанию ("/tmp/devops-metrics-db.json") - строка - имя файла, где хранятся значения (пустое значение - отключает функцию записи на диск)
	envServerKeyValue, isServerKeyValue := os.LookupEnv("KEY")                //(i9) добавьте поддержку аргумента через переменную окружения KEY=<КЛЮЧ>;

	if isEnvHostPort && envHostPortStr != "" { //(i7) Во всех случаях иметь значения по умолчанию и реализовать приоритет значений полученных через ENV, перед значениями задаваемые посредством флагов.
		hostPortStr = envHostPortStr
	}
	if isEnvStoreInterval && envStoreIntervalStr != "" { //(i7) --
		storeIntervalStr = envStoreIntervalStr
	}
	if isEnvStoreInterval && envStoreIntervalStr != "" { //(i7) --
		storeIntervalStr = envStoreIntervalStr
	}
	if isEnvRestore && envRestoreStr != "" { //(i7) --
		if envRestoreStr == "false" {
			isRestoreBool = false
		} else {
			isRestoreBool = true
		}
	}
	if isEnvStoreFilePath && envStoreIntervalStr != "" { //(i7) --
		storeFilePathStr = envStoreFilePath
	}
	if isServerKeyValue && envServerKeyValue != "" { //(i7) --
		serverKeyValue = envServerKeyValue
	}

	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Use(middleware.Compress(5))
	r.Get("/", handlers.AllMetricsHandlerGet)                             //(i3) По запросу GET http://<АДРЕС_СЕРВЕРА>/ сервер должен отдавать html-страничку, со списком имен и значений всех известных ему на текущий момент метрик
	r.Post("/update/", handlers.UpdateHandlerJson)                        //(i4) Для передачи метрик на сервер использовать Content-Type: "application/json", в теле запроса описанный выше JSON, передача через: POST update/
	r.Post("/value/", handlers.ValueHandlerJson)                          //(i4) Для получения метрик с сервера использовать Content-Type: "application/json", в теле запроса описанный выше JSON (заполняем только ID и MType), в ответ получаем такой же JSON, но уже с запоsлненными значениями метрик. Запрос через: POST value/
	r.Post("/update/gauge/*", handlers.UpdateGaugeHandlerPlain)           //(i2) Метрики принимаются сервером по протоколу http, методом POST: в формате: "http://<АДРЕС_СЕРВЕРА>/update/<ТИП_МЕТРИКИ>/<ИМЯ_МЕТРИКИ>/<ЗНАЧЕНИЕ_МЕТРИКИ>" application-type: "text/plain"
	r.Post("/update/counter/*", handlers.UpdateCounterHandlerPlain)       //(i2) --
	r.Post("/update/*", handlers.NotImplementedHandler)                   //(i3) При попытке запроса неизвестной серверу метрики, сервер должен возвращать http.StatusNotFound
	r.Post("/value/*", handlers.NotImplementedHandler)                    //(i3) --
	r.Get("/value/gauge/{MetricName}", handlers.ValueGaugeHandlerGet)     //(i3) Сервер должен возвращать текущее значение запрашиваемой метрики в текстовом виде по запросу GET http://<АДРЕС_СЕРВЕРА>/value/<ТИП_МЕТРИКИ>/<ИМЯ_МЕТРИКИ> со статусом http.StatusOK
	r.Get("/value/counter/{MetricName}", handlers.ValueCounterHandlerGet) //(i3) --

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

	if storeFilePathStr != "" {
		go storage.UpdateMetricsInChannel(c)
		go storage.WriteMetricsToFile(storeFilePathStr, c, storeIntervalTime)
	}
	log.Fatal(s.ListenAndServe())
}
