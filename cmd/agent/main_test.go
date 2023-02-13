package main

import (
	"fmt"
	"net"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"
)

//func (cl Client) sendRequest(curURL string) (string, error) {
func TestSendRequest(t *testing.T) {

	tests := []struct {
		name  string
		input string
		want  string
	}{
		{
			"Positive test",
			"http://127.0.0.1:8080/update/type/231231",
			"200 OK",
		},
		{
			"Wrong IP",
			"http://111.1.1.1:8080/update/type/231231",
			"400 OK",
		},
		{
			"Wrong TCP prot",
			"http://111.1.1.1:9999/update/type/231231",
			"400 OK",
		},
	}

	l, err := net.Listen("tcp", "127.0.0.1:8080")
	if err != nil {
		t.Errorf("Error creating of test server: %s", err)
	}
	s := httptest.NewUnstartedServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
	}))
	s.Listener.Close()
	s.Listener = l
	defer s.Close()
	s.Start()

	var cl Client
	cl.IP = "127.0.0.1"
	cl.Port = "8080"
	cl.Client = s.Client()

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			res, err := cl.sendRequest(tt.input)
			if (res != tt.want) || (err != nil) {
				t.Errorf("Request is %s, want is %s, but response is %s and error is %s", tt.input, tt.want, string(res), err)
			}
		})
	}
}

//func (cl Client) metricSending(mA *[29][3]string) {
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
	defer s.Close()
	url, err := url.Parse(s.URL)
	if err != nil {
		t.Error(err)
	}
	var cl Client
	cl.IP = url.Hostname()
	cl.Port = url.Port()
	//cl.Client = &http.Client{}
	cl.Client = s.Client()

	if err != nil {
		t.Errorf("Error %s", err)
	}
	cl.metricSending(&mArray)
}

//func getMetrics(mArray *[29][3]string, PollCount *int64, rtm *runtime.MemStats) {
