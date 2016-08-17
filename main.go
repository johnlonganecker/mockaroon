package main

import (
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
)

type endpoint struct {
	Hosts   []string          `json:"hosts"`
	Paths   []string          `json:"paths"`
	Status  int               `json:"status"`
	Headers map[string]string `json:"headers"`
	Methods []string          `json:"methods"`
	Body    string            `json:"body"`
}

func (e endpoint) HandleHTTP(w http.ResponseWriter, req *http.Request) {
	for key, value := range e.Headers {
		w.Header().Set(key, value)
	}
	w.WriteHeader(e.Status)
	w.Write([]byte(e.Body))
}

func main() {

	port := ":8080"

	// load in config file

	// create endpoints
	endpoints := make([]endpoint, 1)

	endpoints = append(endpoints, endpoint{
		Paths: []string{"/bob", "/joe", "/moe/joe"},
		Headers: map[string]string{
			"Content-Type": "application/json",
		},
		Methods: []string{"GET", "POST"},
		Body:    "{\"ok\": 10}",
	})

	endpoints = append(endpoints, endpoint{
		Paths: []string{"/okok"},
		Headers: map[string]string{
			"Content-Type": "application/json",
		},
		Methods: []string{"GET", "POST"},
		Body:    "{\"ok\": 10, \"b\": {\"c\": [1,2,3,4,5]}}",
	})

	// create mux router
	muxRouter := mux.NewRouter()

	for _, endpoint := range endpoints {
		for _, path := range endpoint.Paths {
			muxRouter.HandleFunc(path, endpoint.HandleHTTP).Methods(endpoint.Methods...)
			fmt.Println("adding route " + host + port + path)
		}
	}

	http.ListenAndServe(port, muxRouter)
}
