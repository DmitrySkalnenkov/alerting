package main

import (
	"alerting/internal/handlers"
	//"alerting/internal/storage"
	"io"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func main() {
	r := chi.NewRouter()
	r.Use(middleware.Logger)

	//Mstorage := storage.NewMemStorage()
	//Mstorage.Gauges = make(map[string]float64)
	//Mstorage.Counters = make(map[string]int64)

	hg := handlers.GaugesHandler
	hc := handlers.CountersHandler

	hni := func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, http.StatusText(http.StatusNotImplemented), http.StatusNotImplemented)
		io.WriteString(w, "Hello from not implemented handler.\n")
	}
	hggh := handlers.GetGaugeHandler
	hgch := handlers.GetCounterHandler

	r.HandleFunc("/", r.NotFoundHandler())
	r.HandleFunc("/update/*", hni)
	r.HandleFunc("/update/gauge/*", hg)
	r.HandleFunc("/update/counter/*", hc)
	r.HandleFunc("/value/gauge/{MetricName}", hggh)
	r.HandleFunc("/value/counter/{MetricName}", hgch)

	/*http.Handle("/", http.NotFoundHandler())
	http.HandleFunc("/update/", hni)
	http.HandleFunc("/update/gauge/*", hg)
	http.HandleFunc("/update/counter/", hc)*/

	http.ListenAndServe("127.0.0.1:8080", r)
	//server.ListenAndServe()
}
