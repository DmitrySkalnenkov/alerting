package main

import (
	"net/http"
)

// HelloWorld — обработчик запроса.
func HelloWorld(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("<h1>Hello, World</h1>"))
}

func main() {

	server := &http.Server{
		Addr: "127.0.0.1:8080",
	}
	http.HandleFunc("/", HelloWorld)
	server.ListenAndServe()
}
