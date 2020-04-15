package main

import (
	"fmt"
	"net/http"
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

