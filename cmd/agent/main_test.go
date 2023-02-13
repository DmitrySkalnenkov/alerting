package main

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"
)

//func sendRequest(mA *[29][3]string, ip string, port int, cl *http.Client) {
func TestMetricSending(t *testing.T) {

	var mArray [29][3]string
	//1
	mArray[0][0] = "Alloc"
	mArray[0][1] = "gauge"
	mArray[0][2] = "1231.0"
	//2
	mArray[1][0] = "PollCount"
	mArray[1][1] = "counter"
	mArray[1][2] = "123123"
	//3
	mArray[2][0] = "Frees"
	mArray[2][1] = "gauge"
	mArray[2][2] = "3123.0"
	// Start local HTTP server

	s := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		fmt.Fprintln(rw, "Hello, client")
	}))

	url, err := url.Parse(s.URL)
	if err != nil {
		t.Error(err)
	}
	var cl Client
	cl.IP = url.Hostname()
	cl.Port = url.Port()
	cl.Client = &http.Client{}
	//c := s.Client()

	if err != nil {
		t.Errorf("Error %s", err)
	}
	cl.metricSending(&mArray)
	defer s.Close()
}

//func getMetrics(mArray *[29][3]string, PollCount *int64, rtm *runtime.MemStats) {
