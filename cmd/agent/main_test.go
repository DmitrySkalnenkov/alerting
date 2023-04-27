package main

import (
	"github.com/DmitrySkalnenkov/alerting/internal/storage"
	"net"
	"net/http"
	"net/http/httptest"
	"net/url"
	"runtime"
	"strings"
	"testing"
)

func TestSendRequest(t *testing.T) {

	tests := []struct {
		name         string
		input        string
		wantResponse string
		wantMessage  string
	}{
		{
			"Positive test",
			"http://127.0.0.1:8080/update/type/231231",
			"200 OK",
			"",
		},
		{
			"Wrong IP",
			"http://127.0.1.1:8080/update/type/231231",
			"",
			"connect: connection refused",
		},
		{
			"Wrong TCP prot",
			"http://127.0.0.1:9999/update/type/231231",
			"400 OK",
			"connect: connection refused",
		},
	}

	l, err := net.Listen("tcp", "127.0.0.1:8080")
	if err != nil {
		t.Errorf("TEST_ERROR: Test server creating was failed: %s", err)
	}
	s := httptest.NewUnstartedServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
	}))
	//log.Fatal(s.Listener.Close())

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
			if (res != tt.wantResponse) || (err != nil) {
				//fmt.Printf("Error mesage is '%s'\n", err)
				if !strings.Contains(err.Error(), tt.wantMessage) {
					t.Errorf("TEST_ERROR: Request is %s, want is %s, but response is %s and error is %s", tt.input, tt.wantResponse, string(res), err)
				}
			}
		})
	}
}

func TestSendJSONMetric(t *testing.T) {

	type inputs struct {
		url    string
		metric storage.Metrics
	}

	tests := []struct {
		name         string
		input        inputs
		wantResponse string
		wantMessage  string
	}{
		{
			name: "Positive test gauge",
			input: inputs{
				url: "http://127.0.0.1:8080/update/",
				metric: storage.Metrics{
					ID:    "TestMetric1",
					MType: "gauge",
					Value: storage.PointOf(123.321),
				},
			},
			wantResponse: "200 OK",
			wantMessage:  "connect: connection refused",
		},
		{
			name: "Positive test counter",
			input: inputs{
				url: "http://127.0.0.1:8080/update/",
				metric: storage.Metrics{
					ID:    "TestMetric1",
					MType: "counter",
					Delta: storage.PointOf(int64(123)),
				},
			},
			wantResponse: "200 OK",
			wantMessage:  "connect: connection refused",
		},
	}

	l, err := net.Listen("tcp", "127.0.0.1:8080")
	if err != nil {
		t.Errorf("TEST_ERROR: Test server creating was failed: %s", err)
	}

	s := httptest.NewUnstartedServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
	}))
	//log.Fatal(s.Listener.Close())
	s.Listener = l
	defer s.Close()
	s.Start()

	var cl Client
	cl.IP = "127.0.0.1"
	cl.Port = "8080"
	cl.Client = s.Client()

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			res, err := cl.sendJSONMetric(tt.input.url, tt.input.metric)
			if (res != tt.wantResponse) || (err != nil) {
				//fmt.Printf("Error message is '%s'\n", err)
				if !strings.Contains(err.Error(), tt.wantMessage) {
					t.Errorf("TEST_ERROR: Request is %s, want is %s, but response is %s and error is %s", tt.input.url, tt.wantResponse, string(res), err)
				}
			}
		})
	}
}

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
	// Starting of local HTTP server
	s := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
	}))
	defer s.Close()
	urlString, err := url.Parse(s.URL)
	if err != nil {
		t.Error(err)
	}
	var cl Client
	cl.IP = urlString.Hostname()
	cl.Port = urlString.Port()
	cl.Client = s.Client()
	if err != nil {
		t.Errorf("Error %s", err)
	}
	cl.metricSending(&mArray)
}

func TestGetMetrics(t *testing.T) {
	var rtm runtime.MemStats
	var pollcount int64
	var metrics [29][3]string
	var wmArray [29][3]string
	//1
	wmArray[0][0] = "Alloc"
	wmArray[0][1] = "gauge"
	wmArray[0][2] = ""
	//2
	wmArray[1][0] = "BuckHashSys"
	wmArray[1][1] = "gauge"
	wmArray[1][2] = ""
	//3
	wmArray[2][0] = "Frees"
	wmArray[2][1] = "gauge"
	wmArray[2][2] = ""
	//4
	wmArray[3][0] = "GCCPUFraction"
	wmArray[3][1] = "gauge"
	wmArray[3][2] = ""
	//5
	wmArray[4][0] = "GCSys"
	wmArray[4][1] = "gauge"
	wmArray[4][2] = ""
	//6
	wmArray[5][0] = "HeapAlloc"
	wmArray[5][1] = "gauge"
	wmArray[5][2] = ""
	//7
	wmArray[6][0] = "HeapIdle"
	wmArray[6][1] = "gauge"
	wmArray[6][2] = ""
	//8
	wmArray[7][0] = "HeapInuse"
	wmArray[7][1] = "gauge"
	wmArray[7][2] = ""
	//9
	wmArray[8][0] = "HeapObjects"
	wmArray[8][1] = "gauge"
	wmArray[8][2] = ""
	//10
	wmArray[9][0] = "HeapReleased"
	wmArray[9][1] = "gauge"
	wmArray[9][2] = ""
	//11
	wmArray[10][0] = "HeapSys"
	wmArray[10][1] = "gauge"
	wmArray[10][2] = ""
	//12
	wmArray[11][0] = "LastGC"
	wmArray[11][1] = "gauge"
	wmArray[11][2] = ""
	//13
	wmArray[12][0] = "Lookups"
	wmArray[12][1] = "gauge"
	wmArray[12][2] = ""
	//14
	wmArray[13][0] = "MCacheInuse"
	wmArray[13][1] = "gauge"
	wmArray[13][2] = ""
	//15
	wmArray[14][0] = "MCacheSys"
	wmArray[14][1] = "gauge"
	wmArray[14][2] = ""
	//16
	wmArray[15][0] = "MSpanInuse"
	wmArray[15][1] = "gauge"
	wmArray[15][2] = ""
	//17
	wmArray[16][0] = "MSpanSys"
	wmArray[16][1] = "gauge"
	wmArray[16][2] = ""
	//18
	wmArray[17][0] = "Mallocs"
	wmArray[17][1] = "gauge"
	wmArray[17][2] = ""
	//19
	wmArray[18][0] = "NextGC"
	wmArray[18][1] = "gauge"
	wmArray[18][2] = ""
	//20
	wmArray[19][0] = "NumForcedGC"
	wmArray[19][1] = "gauge"
	wmArray[19][2] = ""
	//21
	wmArray[20][0] = "NumGC"
	wmArray[20][1] = "gauge"
	wmArray[20][2] = ""
	//22
	wmArray[21][0] = "OtherSys"
	wmArray[21][1] = "gauge"
	wmArray[21][2] = ""
	//23
	wmArray[22][0] = "PollCount"
	wmArray[22][1] = "counter"
	wmArray[22][2] = ""
	//24
	wmArray[23][0] = "PauseTotalNs"
	wmArray[23][1] = "gauge"
	wmArray[23][2] = ""
	//25
	wmArray[24][0] = "RandomValue"
	wmArray[24][1] = "gauge"
	wmArray[24][2] = ""
	//26
	wmArray[25][0] = "StackInuse"
	wmArray[25][1] = "gauge"
	wmArray[25][2] = ""
	//27
	wmArray[26][0] = "StackSys"
	wmArray[26][1] = "gauge"
	wmArray[26][2] = ""
	//28
	wmArray[27][0] = "Sys"
	wmArray[27][1] = "gauge"
	wmArray[27][2] = ""
	//29
	wmArray[28][0] = "TotalAlloc"
	wmArray[28][1] = "gauge"
	wmArray[28][2] = ""

	getMetrics(&metrics, &pollcount, &rtm)
	//if len(metrics) != len(wmArray) {
	//	t.Errorf("TEST_ERROR: Wrong length of metric array. Length is %d.\n", len(metrics))
	//}
	//fmt.Printf("Metrics %v: \n", metrics)
	for i := 1; i < 29; i++ {
		if metrics[i][0] != wmArray[i][0] || metrics[i][1] != wmArray[i][1] {
			t.Errorf("TEST_ERROR: Wrong name or type of metric. Actual name and type is %s - %s, want name and type is %s - %s.\n",
				metrics[i][0], metrics[i][1], wmArray[i][0], wmArray[i][1])
		}
		if metrics[i][2] == wmArray[i][2] {
			t.Errorf("TEST_ERROR: No metric value for %s. Value is %s.\n", metrics[i][0], metrics[i][2])
		}
	}

}

func TestMetricSendingAPI2(t *testing.T) {
	type fields struct {
		IP     string
		Port   string
		Client *http.Client
	}
	type args struct {
		mA *[29][3]string
	}

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

	// Starting of local HTTP server
	s := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		//fmt.Fprintln(rw, "Hello, client")
	}))
	defer s.Close()
	urlString, err := url.Parse(s.URL)
	if err != nil {
		t.Error(err)
	}
	tests := []struct {
		name   string
		fields fields
		args   args
	}{
		{name: `positive test #1`,
			fields: fields{
				IP:     urlString.Hostname(),
				Port:   urlString.Port(),
				Client: s.Client(),
			},
			args: args{
				&mArray,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cl := Client{
				IP:     tt.fields.IP,
				Port:   tt.fields.Port,
				Client: tt.fields.Client,
			}
			cl.metricSendingAPI2(tt.args.mA)
		})
	}
}
