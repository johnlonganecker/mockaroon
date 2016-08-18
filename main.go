package main

import (
	"fmt"
	"io/ioutil"
	"net/http"

	"gopkg.in/yaml.v2"

	"github.com/gorilla/mux"
)

type Endpoint struct {
	Paths   []string            `yaml:"paths"`
	Status  int                 `yaml:"status"`
	Headers []map[string]string `yaml:"headers"`
	Methods []string            `yaml:"methods"`
	Body    string              `yaml:"body"`
}

type Config struct {
	Serve     bool       `yaml:"serveFiles",`
	Port      string     `yaml:"port",`
	Endpoints []Endpoint `yaml:"endpoints",`
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

	port := "8080"

	config := Config{}

	// load in config file
	err := config.LoadConfigFile("sample-config.yml")
	if err != nil {
		fmt.Println(err)
	}

	if config.Port != "" {
		port = config.Port
	}
	port = ":" + port

	// create mux router
	muxRouter := mux.NewRouter()

	for _, endpoint := range config.Endpoints {
		for _, path := range endpoint.Paths {
			muxRouter.HandleFunc(path, endpoint.HandleHTTP).Methods(endpoint.Methods...)
			fmt.Println("adding route " + port + path)
		}
	}

	if config.Serve == true {
		muxRouter.PathPrefix("/").Handler(http.FileServer(http.Dir("./")))
	}
	fmt.Println(config.Serve)

	http.ListenAndServe(port, muxRouter)
}
