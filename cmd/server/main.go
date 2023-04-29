package main

import (
	"io"
	"log"
	"net/http"

	"github.com/DmitrySkalnenkov/alerting/internal/handlers"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func main() {
	r := chi.NewRouter()
	r.Use(middleware.Logger)

	hni := func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, http.StatusText(http.StatusNotImplemented), http.StatusNotImplemented)
		_, err := io.WriteString(w, "Hello from not implemented handler.\n")
		if err != nil {
			log.Fatal(err)
		}
	}

	r.HandleFunc("/", handlers.GetAllMetricsHandlerAPI2)
	r.HandleFunc("/update/*", hni)
	r.HandleFunc("/update/gauge/*", handlers.GaugesHandlerAPI2)
	r.HandleFunc("/update/counter/*", handlers.CounterHandlerAPI2)
	r.HandleFunc("/value/gauge/{MetricName}", handlers.GetGaugeHandlerAPI2)
	r.HandleFunc("/value/counter/{MetricName}", handlers.GetCounterHandlerAPI2)

	log.Fatal(http.ListenAndServe("127.0.0.1:8080", r))

}
