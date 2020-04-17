package main

import (
	"io/ioutil"
	"log"
)

var (
	jsonMetricsFileName = "/tmp/metrics"
)

func writeJsonToFile(metricsString []byte) error {
	err := ioutil.WriteFile(jsonMetricsFileName, metricsString, 0666)
	if err != nil {
		log.Println("Writing json to file failed!!")
	}
	return err
}


func readJsonFromFile() (string, error) {
	metricsBytes, err := ioutil.ReadFile(jsonMetricsFileName)
	if err != nil {
		log.Println("read json from file failed!!")
		return "", err
	} else {
		return string(metricsBytes), nil
	}
}