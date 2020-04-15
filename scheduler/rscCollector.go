package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/prometheus/client_golang/api"
)

type MetricsInfo struct {
	CpuUtilization       float64
	MemoryUtilization    float64
	DiskIO    float64
	NetworkIO float64
}

type Metrics struct {
	Namespace       string
	PodName         string
	TimeStamp       int64
	Infos           MetricsInfo
}

// prometheusAddr is the prometheus IP:port
var prometheusAddr = "http://localhost:9090"

// jsonMetrics used for storing metrics infos
// Key is namespace + podName, Value is Metrics
var jsonMetrics = make(map[string]Metrics)


func getPromClient() api.Client {
	client, err := api.NewClient(api.Config{
		Address: prometheusAddr,
	})
	if err != nil {
		fmt.Printf("Error creating client: %v\n", err)
		os.Exit(1)
	}

	return client
}


func query(client api.Client, currentTime time.Time) {
	networkQuery(client, currentTime)
	cpuQuery(client, currentTime)
	diskIOQuery(client, currentTime)
	memoryQuery(client, currentTime)
}

func main() {
	client := getPromClient()
	currentTime := time.Now()

	// TODO: (1)use goroutines to deal with four query
	// TODO: (2)protect the jsonMetrics using locks
	// TODO: (3)provide a http interface for the metrics
	// TODO: (4)understand the queryRange and series
	// TODO: (5)learn the CGroup
	//for {
	query(client, currentTime)

	if metricsString, err := json.Marshal(jsonMetrics); err != nil {
		fmt.Println(err)
	} else {
		handler := Handler{response: string(metricsString)}
		go log.Fatal(http.ListenAndServe("localhost:33000", handler))
	}
}