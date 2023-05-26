package main

import (
	"github.com/DmitrySkalnenkov/alerting/internal/handlers"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"log"
	"net/http"
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
	//TODO: Change HandleFunc to Get or Post func
	r.HandleFunc("/", handlers.GetAllMetricsHandler)
	r.HandleFunc("/update/*", handlers.UpdateHandler)
	r.Post("/value/*", handlers.ValueHandler)
	//r.HandleFunc("/update/gauge/*", handlers.GaugeHandlerAPI1)
	//r.HandleFunc("/update/counter/*", handlers.CounterHandler)
	r.Get("/value/gauge/{MetricName}", handlers.GetGaugeHandlerAPI1)
	//r.HandleFunc("/value/gauge/{MetricName}", handlers.GetGaugeHandler)
	r.Get("/value/counter/{MetricName}", handlers.GetCounterHandlerAPI1)

	log.Fatal(http.ListenAndServe("127.0.0.1:8080", r))

}
