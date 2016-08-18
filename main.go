// mockaroon [--config=path/to/config] [port]

package main

import (
	"errors"
	"flag"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strconv"

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
	ServeFiles bool       `yaml:"serveFiles",omitempty`
	Port       string     `yaml:"port"`
	Endpoints  []Endpoint `yaml:"endpoints"`
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

func handleError(err error) {
	fmt.Println(err)
	os.Exit(2)
}

func main() {

	var configPath string

	// defaults
	config := Config{
		Port:       "8000",
		ServeFiles: false,
	}

	// parse flags
	flag.StringVar(&configPath, "config", "", "Path to Config File")
	flag.Parse()

	// load in config file
	if configPath != "" {
		err := config.LoadConfigFile(configPath)
		if err != nil {
			handleError(err)
		}
	}

	// command line port overrides all
	tail := flag.Args()
	if len(tail) > 0 {
		config.Port = tail[0]
	}

	// validate port
	port, err := strconv.Atoi(config.Port)
	if err != nil {
		handleError(err)
	}
	if port < 1 {
		handleError(errors.New(config.Port + " is not a valid port"))
	}

	// if port is not last param
	if len(tail) > 1 {
		fmt.Println("Warning: port goes at the end, all params after ignored")
	}

	// create mux router
	muxRouter := mux.NewRouter()

	for _, endpoint := range config.Endpoints {
		for _, path := range endpoint.Paths {
			muxRouter.HandleFunc(path, endpoint.HandleHTTP).Methods(endpoint.Methods...)
			fmt.Println("adding route " + path)
		}
	}

	if config.ServeFiles == true {
		muxRouter.PathPrefix("/").Handler(http.FileServer(http.Dir("./")))
		fmt.Println("Serving static files")
	}

	fmt.Println("listening on port " + config.Port)
	http.ListenAndServe(":"+config.Port, muxRouter)
}
