// mockaroon [--config=path/to/config] [port]

// TODO make formatting match Python -m SimpleHTTPServer
// 127.0.0.1 - - [24/Aug/2016 10:27:39] "GET / HTTP/1.1" 200 -
// fmt.Printf("%s - - [%s/%s/%s %s:%s:%s] \"%s %s %s\" %s -\n", ip, day, month, year, hour, minute, second, method, path, httpVersion, status)

package main

import (
	"errors"
	"fmt"
	"io/ioutil"
	"math/rand"
	"net/http"
	"net/http/httputil"
	"net/url"
	"os"
	"strconv"
	"time"

	"gopkg.in/yaml.v2"

	"github.com/gorilla/mux"

	"github.com/urfave/cli"
)

type Endpoint struct {
	Paths   []string            `yaml:"paths"`
	Status  int                 `yaml:"status"`
	Headers []map[string]string `yaml:"headers"`
	Methods []string            `yaml:"methods"`
	Body    string              `yaml:"body"`
	Latency Latency             `yaml:"latency"`
}

type Latency struct {
	Min int `yaml:"min"`
	Max int `yaml:"max"`
}

type Proxy struct {
	Paths       []string `yaml:"paths"`
	Destination string   `yaml:"destination"`
}

type SSL struct {
	Cert    string `yaml:"cert"`
	Private string `yaml:"private"`
}

type Config struct {
	ServeFiles bool       `yaml:"serveFiles",omitempty`
	Port       int        `yaml:"port"`
	SSL        SSL        `yaml:"ssl"`
	Endpoints  []Endpoint `yaml:"endpoints"`
	Proxies    []Proxy    `yaml:"proxies"`
}

type Context struct {
	Latency Latency
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

func (p Proxy) HandleHttp(w http.ResponseWriter, req *http.Request) {

	url, _ := url.Parse(p.Destination)

	// set the proper host
	req.Host = url.Host

	// make sure we set the proper host for the proxy
	w.Header().Set("Host", url.Host+req.URL.Path)

	// TODO perhaps don't make this on every request
	proxy := httputil.NewSingleHostReverseProxy(url)

	proxy.ServeHTTP(w, req)
}

func CalcLatency(min int, max int) (int, error) {

	if min == 0 && max == 0 {
		return 0, nil
	}

	if max < min {
		return 0, errors.New("Latency Max has to be less then Min")
	}

	timeRange := max - min

	return (rand.Intn(timeRange) + min), nil
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

type Middleware struct {
	handler http.HandlerFunc
	context Context
}

func (m *Middleware) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	latency, err := CalcLatency(m.context.Latency.Min, m.context.Latency.Max)
	if err != nil {
		handleError(err)
	}
	time.Sleep(time.Duration(latency) * time.Millisecond)
	m.handler.ServeHTTP(w, req)
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
	finalOutput := ""

	app := cli.NewApp()
	app.Usage = "A Simple HTTPS Server for local development"
	app.UsageText = "mockaroon [global options] [port]"

	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:        "config, c",
			Usage:       "Load configuration from `FILE`",
			Destination: &configPath,
		},
	}
	app.Version = "0.1.0"
	app.Compiled = time.Now()
	app.Authors = []cli.Author{
		cli.Author{
			Name:  "John Longanecker",
			Email: "johnlonganecker@gmail.com",
		},
	}

	// defaults
	config := Config{
		Port:       8000,
		ServeFiles: true,
	}

	app.Action = func(c *cli.Context) error {

		var err error

		// load in config file
		if configPath != "" {
			err := config.LoadConfigFile(configPath)
			if err != nil {
				handleError(err)
			}
		}

		// command line port overrides all
		tail := c.Args()
		if len(tail) > 0 {
			config.Port, err = strconv.Atoi(tail[0])
			if err != nil {
				handleError(err)
			}
		}

		// validate port
		port := strconv.Itoa(config.Port)
		if config.Port < 1 {
			handleError(errors.New(port + " is not a valid port"))
		}

		// if port is not last param
		if len(tail) > 1 {
			fmt.Println("Warning: port goes at the end, all params after ignored")
		}

		// create mux router
		muxRouter := mux.NewRouter()

		for _, endpoint := range config.Endpoints {
			for _, path := range endpoint.Paths {
				m := Middleware{
					handler: endpoint.HandleHTTP,
					context: Context{Latency{Min: endpoint.Latency.Min, Max: endpoint.Latency.Max}},
				}
				muxRouter.HandleFunc(path, m.ServeHTTP).Methods(endpoint.Methods...)
				finalOutput += "adding route " + path + "\n"
			}
		}

		for _, proxy := range config.Proxies {
			for _, path := range proxy.Paths {
				muxRouter.HandleFunc(path, proxy.HandleHttp)
				finalOutput += "adding proxy route " + path + "\n"
			}
		}

		if config.ServeFiles == true {
			muxRouter.PathPrefix("/").Handler(http.FileServer(http.Dir("./")))
			finalOutput = fmt.Sprintf("Serving static files\n") + finalOutput
		}

		if config.SSL.Cert == "" && config.SSL.Private == "" {
			finalOutput = fmt.Sprintf("Serving HTTP on 0.0.0.0 port %d ...\n", config.Port) + finalOutput
			fmt.Print(finalOutput)
			http.ListenAndServe(":"+port, muxRouter)
		} else {
			finalOutput = fmt.Sprintf("Serving HTTPS on 0.0.0.0 port %d ...\n", config.Port) + finalOutput
			fmt.Print(finalOutput)
			err := http.ListenAndServeTLS(":"+port, config.SSL.Cert, config.SSL.Private, muxRouter)
			if err != nil {
				fmt.Println(err)
			}
		}

		return nil
	}

	app.Run(os.Args)
}
