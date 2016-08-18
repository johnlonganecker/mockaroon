package main

import (
	"fmt"
	"io/ioutil"
	"net/http"

	"gopkg.in/yaml.v2"

	"github.com/gorilla/mux"
)

type Endpoint struct {
	Paths   []string            `json:"paths"`
	Status  int                 `json:"status"`
	Headers []map[string]string `json:"headers"`
	Methods []string            `json:"methods"`
	Body    string              `json:"body"`
}

type Config struct {
	Port      string     `json:"Port"`
	Endpoints []Endpoint `json:"Endpoint"`
}

func (e Endpoint) HandleHTTP(w http.ResponseWriter, req *http.Request) {
	for _, header := range e.Headers {
		for key, value := range header {
			w.Header().Set(key, value)
		}
	}
	w.WriteHeader(e.Status)
	w.Write([]byte(e.Body))
}

func (c *Config) LoadConfigFile(filepath string) error {

	var data []byte

	data, err := LoadFile(filepath)
	if err != nil {
		return err
	}

	if err := Unmarshal(c, data); err != nil {
		return err
	}

	return nil
}

func LoadFile(filepath string) ([]byte, error) {

	data, err := ioutil.ReadFile(filepath)
	if err != nil {
		return []byte{}, err
	}

	return data, nil
}

func Unmarshal(c *Config, data []byte) error {

	// unmarshal yaml
	err := yaml.Unmarshal(data, c)
	if err != nil {
		return err
	}

	return nil
}

func main() {

	port := ":8080"

	config := Config{}

	err := config.LoadConfigFile("sample-config.json")
	if err != nil {
		fmt.Println(err)
	}

	// load in config file

	// create Endpoints
	//endpoints := config.Endpoints

	//endpoints = append(endpoints, Endpoint{
	//Paths: []string{"/bob", "/joe", "/moe/joe"},
	//Headers: map[string]string{
	//"Content-Type": "application/json",
	//},
	//Methods: []string{"GET", "POST"},
	//Body:    "{\"ok\": 10}",
	//})

	//endpoints = append(endpoints, Endpoint{
	//Paths: []string{"/okok"},
	//Headers: map[string]string{
	//"Content-Type": "application/json",
	//},
	//Methods: []string{"GET", "POST"},
	//Body:    "{\"ok\": 10, \"b\": {\"c\": [1,2,3,4,5]}}",
	//})

	// create mux router
	muxRouter := mux.NewRouter()

	for _, endpoint := range config.Endpoints {
		for _, path := range endpoint.Paths {
			muxRouter.HandleFunc(path, endpoint.HandleHTTP).Methods(endpoint.Methods...)
			fmt.Println("adding route " + port + path)
		}
	}

	http.ListenAndServe(port, muxRouter)
}
