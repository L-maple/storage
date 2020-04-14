package main

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/prometheus/client_golang/api"
	v1 "github.com/prometheus/client_golang/api/prometheus/v1"
)

var prometheusAddr = "http://localhost:9090"

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

func cpuQuery() {
	client := getPromClient()

	v1api := v1.NewAPI(client)
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	queryQL := "sum(rate(container_cpu_usage_seconds_total{image != ''}[1m])) by (pod, namespace)"
	result, warnings, err := v1api.Query(ctx, queryQL, time.Now())

	if err != nil {
		fmt.Printf("Error querying Prometheus: %v\n", err)
		os.Exit(1)
	}
	if len(warnings) > 0 {
		fmt.Printf("Warnings: %v\n", warnings)
	}
	fmt.Printf("Result:\n%v\n", result)
}

func memoryQuery() {
	client := getPromClient()

	v1api := v1.NewAPI(client)
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	queryQL := "sum(container_memory_usage_bytes{image!=''}) by(pod, namespace)"
	result, warnings, err := v1api.Query(ctx, queryQL, time.Now())
	if err != nil {
		fmt.Printf("Error querying Prometheus: %v\n", err)
		os.Exit(1)
	}
	if len(warnings) > 0 {
		fmt.Printf("Warnings: %v\n", warnings)
	}
	fmt.Printf("Result:\n%v\n", result)
}

func networkQuery() {
	client := getPromClient()

	v1api := v1.NewAPI(client)
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	queryQL := "sum(rate(container_network_transmit_bytes_total{image != ''}[1m])) by (pod, namespace)"
	result, warnings, err := v1api.Query(ctx, queryQL, time.Now())

	if err != nil {
		fmt.Printf("Error querying Prometheus: %v\n", err)
		os.Exit(1)
	}
	if len(warnings) > 0 {
		fmt.Printf("Warnings: %v\n", warnings)
	}
	fmt.Printf("Result:\n%v\n", result)
}

func diskIOQuery() {
	client := getPromClient()

	v1api := v1.NewAPI(client)
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	queryQL := "sum(rate(container_fs_writes_bytes_total{image!=''}[1m])) by (pod, namespace) + sum(rate(container_fs_reads_bytes_total{image!=''}[1m])) by (pod, namespace)"
	result, warnings, err := v1api.Query(ctx, queryQL, time.Now())

	if err != nil {
		fmt.Printf("Error querying Prometheus: %v\n", err)
		os.Exit(1)
	}
	if len(warnings) > 0 {
		fmt.Printf("Warnings: %v\n", warnings)
	}
	fmt.Printf("Result:\n%v\n", result)
}

func queryRange() {
	client := getPromClient()

	v1api := v1.NewAPI(client)
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	r := v1.Range{
		Start: time.Now().Add(-time.Hour),
		End:   time.Now(),
		Step:  time.Minute,
	}
	result, warnings, err := v1api.QueryRange(ctx, "rate(prometheus_tsdb_head_samples_appended_total[5m])", r)
	if err != nil {
		fmt.Printf("Error querying Prometheus: %v\n", err)
		os.Exit(1)
	}
	if len(warnings) > 0 {
		fmt.Printf("Warnings: %v\n", warnings)
	}
	fmt.Printf("Result:\n%v\n", result)
}

func series() {
	client := getPromClient()

	v1api := v1.NewAPI(client)
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	lbls, warnings, err := v1api.Series(ctx, []string{
		"{__name__=~\"scrape_.+\",job=\"node\"}",
		"{__name__=~\"scrape_.+\",job=\"prometheus\"}",
	}, time.Now().Add(-time.Hour), time.Now())
	if err != nil {
		fmt.Printf("Error querying Prometheus: %v\n", err)
		os.Exit(1)
	}
	if len(warnings) > 0 {
		fmt.Printf("Warnings: %v\n", warnings)
	}
	fmt.Println("Result:")
	for _, lbl := range lbls {
		fmt.Println(lbl)
	}
}

func main() {
	networkQuery()
}