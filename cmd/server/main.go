package main

import (
	"alerting/internal/handlers"
	"io"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func main() {
	r := chi.NewRouter()
	r.Use(middleware.Logger)

	hni := func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, http.StatusText(http.StatusNotImplemented), http.StatusNotImplemented)
		io.WriteString(w, "Hello from not implemented handler.\n")
	}

	r.HandleFunc("/", handlers.GetAllMetricsHandler)
	r.HandleFunc("/update/*", hni)
	r.HandleFunc("/update/gauge/*", handlers.GaugesHandler)
	r.HandleFunc("/update/counter/*", handlers.CountersHandler)
	r.HandleFunc("/value/gauge/{MetricName}", handlers.GetGaugeHandler)
	r.HandleFunc("/value/counter/{MetricName}", handlers.GetCounterHandler)

	http.ListenAndServe("127.0.0.1:8080", r)
}
