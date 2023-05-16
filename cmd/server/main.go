package main

import (
	"log"
	"net/http"

	"github.com/DmitrySkalnenkov/alerting/internal/handlers"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func main() {
	r := chi.NewRouter()
	r.Use(middleware.Logger)

	//hni := func(w http.ResponseWriter, r *http.Request) {
	//	http.Error(w, http.StatusText(http.StatusNotImplemented), http.StatusNotImplemented)
	//	_, err := io.WriteString(w, "Hello from not implemented handler.\n")
	//	if err != nil {
	//		log.Fatal(err)
	//	}
	//}
	//TODO: update -> If POST then GaugeHandlerAPI2, if GET then GaugeHandlerAPI1
	r.HandleFunc("/", handlers.GetAllMetricsHandler)
	r.HandleFunc("/update/*", handlers.UpdateHandler)
	//r.HandleFunc("/value/*", handlers.UpdateHandler)
	//r.HandleFunc("/update/gauge/*", handlers.GaugeHandler)
	//r.HandleFunc("/update/counter/*", handlers.CounterHandler)
	r.HandleFunc("/value/gauge/{MetricName}", handlers.GetGaugeHandler)
	r.HandleFunc("/value/counter/{MetricName}", handlers.GetCounterHandler)

	log.Fatal(http.ListenAndServe("127.0.0.1:8080", r))

}
