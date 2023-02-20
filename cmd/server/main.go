package server

import (
	//"alerting/internal"
	"alerting/internal/handlers"
	//"alerting/internal/storage"
	//	"fmt"
	"io"
	"net/http"
	//"regexp"
	//"strconv"
	//"strings"
)

/*func contains(s []string, str string) bool {
	for _, v := range s {
		if v == str {
			return true
		}
	}

	return false
}*/

func main() {
	server := &http.Server{
		Addr: "127.0.0.1:8080",
	}
	handlers.Mstorage.Gauges = make(map[string]float64)
	handlers.Mstorage.Counters = make(map[string]int64)

	hg := handlers.GaugesHandler
	hc := handlers.CountersHandler

	hni := func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, http.StatusText(http.StatusNotImplemented), http.StatusNotImplemented)
		io.WriteString(w, "Hello from not implemented handler.\n")

	}
	http.Handle("/", http.NotFoundHandler())
	http.HandleFunc("/update/", hni)
	http.HandleFunc("/update/gauge/", hg)
	http.HandleFunc("/update/counter/", hc)
	server.ListenAndServe()
}
