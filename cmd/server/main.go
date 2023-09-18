package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/DmitrySkalnenkov/alerting/internal/handlers"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

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
	hostportStr := "127.0.0.1:8080"
	if hostportStr != "" {
		hostportStr = os.Getenv("ADDRESS")
	}
	fmt.Println(hostportStr)
	s := &http.Server{
		Addr:         hostportStr,
		ReadTimeout:  30 * time.Second,
		WriteTimeout: 30 * time.Second,
	}
	s.Handler = r
	log.Fatal(s.ListenAndServe())
}
