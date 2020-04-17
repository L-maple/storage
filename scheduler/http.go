package main

import (
	"fmt"
	"log"
	"net/http"
	"time"
)

type Handler struct {
	response string
}


func (receiver Handler) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	switch req.URL.Path {
	case "/metrics":
		fmt.Fprintf(w, "%s", receiver.response)
	default:
		w.WriteHeader(http.StatusNotFound)
		fmt.Fprintf(w, "No such page: %s\n", req.URL.Path)
	}
}


func main() {
	handler := Handler{response : ""}

	go func() {
		log.Fatal(http.ListenAndServe("localhost:33000", &handler))
	}()

	for {
		if metricsString, err := readJsonFromFile(); err != nil {
			handler.response = ""
			time.Sleep(time.Second * 3)
			continue
		} else {
			handler.response = metricsString
		}

		time.Sleep(time.Millisecond * 500)
	}
}