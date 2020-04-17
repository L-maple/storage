package main

import (
	"context"
	"fmt"
	"github.com/prometheus/client_golang/api"
	v1 "github.com/prometheus/client_golang/api/prometheus/v1"
	"os"
	"regexp"
	"strconv"
	"strings"
	"time"
)


func metricsParser(entry string) (string, string, float64) {
	if entry == "" {
		return "", "", 0.0
	}
	re := regexp.MustCompile(`\{namespace="([\w-_]+)", pod="([\w-_]+)"\} => ([\.0-9]+) @\[([0-9\.]+)\]`)
	params := re.FindSubmatch([]byte(entry))
	namespace := string(params[1])
	podName := string(params[2])
	metric, _ := strconv.ParseFloat(string(params[3]), 64)
	return namespace, podName, metric
}


func cpuQuery(client api.Client, currentTime time.Time) {
	v1api := v1.NewAPI(client)
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	queryQL := "sum(rate(container_cpu_usage_seconds_total{image != ''}[1m])) by (pod, namespace)"
	result, warnings, err := v1api.Query(ctx, queryQL, currentTime)

	if err != nil {
		fmt.Printf("Error querying Prometheus: %v\n", err)
		os.Exit(1)
	}
	if len(warnings) > 0 {
		fmt.Printf("Warnings: %v\n", warnings)
	}
	entries := strings.Split(result.String(), "\n")
	for _, entry := range entries {
		namespace, podName, cpuUtilization := metricsParser(entry)
		if _, ok := jsonMetrics[namespace + podName]; ok {  // this pod has exist in jsonMetrics
			metricsInfo := jsonMetrics[namespace + podName]
			metricsInfo.Infos.CpuUtilization = cpuUtilization
			jsonMetrics[namespace + podName] = metricsInfo
		} else {
			metricsInfo := Metrics {
				Namespace: namespace,
				PodName: podName,
				TimeStamp: currentTime.Unix(),
				Infos : MetricsInfo{
					CpuUtilization:    cpuUtilization,
					MemoryUtilization: 0,
					DiskIO:            0,
					NetworkIO:         0,
				},
			}
			jsonMetrics[namespace + podName] = metricsInfo
		}
	}
}


func memoryQuery(client api.Client, currentTime time.Time) {
	v1api := v1.NewAPI(client)
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	queryQL := "sum(container_memory_usage_bytes{image!=''}) by(pod, namespace)"
	result, warnings, err := v1api.Query(ctx, queryQL, currentTime)
	if err != nil {
		fmt.Printf("Error querying Prometheus: %v\n", err)
		os.Exit(1)
	}
	if len(warnings) > 0 {
		fmt.Printf("Warnings: %v\n", warnings)
	}
	entries := strings.Split(result.String(), "\n")
	for _, entry := range entries {
		namespace, podName, memoryUtilization := metricsParser(entry)
		if _, ok := jsonMetrics[namespace+podName]; ok { // this pod has exist in jsonMetrics
			metricsInfo := jsonMetrics[namespace+podName]
			metricsInfo.Infos.MemoryUtilization = memoryUtilization
			jsonMetrics[namespace+podName] = metricsInfo
		} else {
			metricsInfo := Metrics{
				Namespace:   namespace,
				PodName:     podName,
				TimeStamp:   currentTime.Unix(),
				Infos: MetricsInfo{
					CpuUtilization:    0,
					MemoryUtilization: memoryUtilization,
					DiskIO:            0,
					NetworkIO:         0,
				},
			}
			jsonMetrics[namespace+podName] = metricsInfo
		}
	}
}


func networkQuery(client api.Client, currentTime time.Time) {
	v1api := v1.NewAPI(client)
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	queryQL := "sum(rate(container_network_transmit_bytes_total{image != ''}[1m])) by (pod, namespace)"
	result, warnings, err := v1api.Query(ctx, queryQL, currentTime)

	if err != nil {
		fmt.Printf("Error querying Prometheus: %v\n", err)
		os.Exit(1)
	}
	if len(warnings) > 0 {
		fmt.Printf("Warnings: %v\n", warnings)
	}
	entries := strings.Split(result.String(), "\n")
	for _, entry := range entries {
		namespace, podName, networkIO := metricsParser(entry)
		if _, ok := jsonMetrics[namespace+podName]; ok { // this pod has exist in jsonMetrics
			metricsInfo := jsonMetrics[namespace+podName]
			metricsInfo.Infos.NetworkIO = networkIO
			jsonMetrics[namespace+podName] = metricsInfo
		} else {
			metricsInfo := Metrics{
				Namespace:   namespace,
				PodName:     podName,
				TimeStamp:   currentTime.Unix(),
				Infos: MetricsInfo{
					CpuUtilization:    0,
					MemoryUtilization: 0,
					DiskIO:            0,
					NetworkIO:         networkIO,
				},
			}
			jsonMetrics[namespace+podName] = metricsInfo
		}
	}
}


func diskIOQuery(client api.Client, currentTime time.Time) {
	v1api := v1.NewAPI(client)
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	queryQL := "sum(rate(container_fs_writes_bytes_total{image!=''}[1m])) by (pod, namespace) + sum(rate(container_fs_reads_bytes_total{image!=''}[1m])) by (pod, namespace)"
	result, warnings, err := v1api.Query(ctx, queryQL, currentTime)

	if err != nil {
		fmt.Printf("Error querying Prometheus: %v\n", err)
		os.Exit(1)
	}
	if len(warnings) > 0 {
		fmt.Printf("Warnings: %v\n", warnings)
	}
	entries := strings.Split(result.String(), "\n")
	for _, entry := range entries {
		namespace, podName, diskIO := metricsParser(entry)
		if _, ok := jsonMetrics[namespace+podName]; ok { // this pod has exist in jsonMetrics
			metricsInfo := jsonMetrics[namespace+podName]
			metricsInfo.Infos.DiskIO = diskIO
			jsonMetrics[namespace+podName] = metricsInfo
		} else {
			metricsInfo := Metrics{
				Namespace:   namespace,
				PodName:     podName,
				TimeStamp:   currentTime.Unix(),
				Infos: MetricsInfo{
					CpuUtilization:    0,
					MemoryUtilization: 0,
					DiskIO:            diskIO,
					NetworkIO:         0,
				},
			}
			jsonMetrics[namespace+podName] = metricsInfo
		}
	}
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

