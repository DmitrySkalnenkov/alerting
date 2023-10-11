package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"

	//"strconv"
	"time"

	"github.com/DmitrySkalnenkov/alerting/internal/handlers"
	"github.com/DmitrySkalnenkov/alerting/internal/storage"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func GetEnvVariable(envVarName string, defaultValue string) string {
	envVarValue := defaultValue
	if os.Getenv(envVarName) != "" {
		envVarValue = os.Getenv(envVarName)
	}
	fmt.Printf("DEBUG: Variable '%s' has value '%s'.\n", envVarName, envVarValue)
	return envVarValue
}

func main() {
	//time.Sleep(500 * time.Millisecond)
	r := chi.NewRouter()
	r.Use(middleware.Logger)
	//hni := func(w http.ResponseWriter, r *http.Request) {
	//	http.Error(w, http.StatusText(http.StatusNotImplemented), http.StatusNotImplemented)
	//	_, err := io.WriteString(w, "Hello from not implemented handler.\n")
	//	if err != nil {
	//		log.Fatal(err)
	//	}
	//}
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

	hostportStr := GetEnvVariable("ADDRESS", "127.0.0.1:8080")
	storeIntervalStr := GetEnvVariable("STORE_INTERVAL", "300")
	storeIntervalStr += "s"
	storeFilePath := GetEnvVariable("STORE_FILE", "/tmp/devops-metrics-db.json")
	isRestoreStr := GetEnvVariable("RESTORE", "true")
	storeIntervalTime, err := time.ParseDuration(storeIntervalStr)
	if err != nil {
		log.Fatal(err)
	}

	s := &http.Server{
		Addr:         hostportStr,
		ReadTimeout:  30 * time.Second,
		WriteTimeout: 30 * time.Second,
	}
	s.Handler = r

	var c chan storage.MetricsStorage = make(chan storage.MetricsStorage)
	if strings.ToLower(isRestoreStr) == "true" {
		storage.RestoreMetricsFromFile(storeFilePath, storage.MetStorage)
	}

	go storage.UpdateMetricsInChannel(c)
	go storage.WriteMetricsToFile(storeFilePath, c, storeIntervalTime)

	log.Fatal(s.ListenAndServe())
	var input string
	fmt.Scanln(&input)
}
