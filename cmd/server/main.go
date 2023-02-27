package main

import (
	"alerting/internal/handlers"
	"alerting/internal/storage"
	"io"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func main() {
	/*server := &http.Server{
		Addr: "127.0.0.1:8080",
	}*/

	r := chi.NewRouter()
	r.Use(middleware.Logger)

	storage.Mstorage.Gauges = make(map[string]float64)
	storage.Mstorage.Counters = make(map[string]int64)

	hg := handlers.GaugesHandler
	hc := handlers.CountersHandler

	hni := func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, http.StatusText(http.StatusNotImplemented), http.StatusNotImplemented)
		io.WriteString(w, "Hello from not implemented handler.\n")
	}
	hggh := handlers.GetGaugeHandler

	r.Get("/", r.NotFoundHandler())
	r.Get("/update", hni)
	r.Get("/update/gauge/", hg)
	r.Get("/update/counter/", hc)
	r.Get("/values/gauge/{MetricName}", hggh)

	/*http.Handle("/", http.NotFoundHandler())
	http.HandleFunc("/update/", hni)
	http.HandleFunc("/update/gauge/", hg)
	http.HandleFunc("/update/counter/", hc)*/

	http.ListenAndServe("127.0.0.1:8080", r)
	//server.ListenAndServe()
}
