package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httputil"
	"net/url"
	"strings"
	"time"
)

var services map[string][]Service

func main() {

	services = make(map[string][]Service)

	checkServicesHealth()

	http.HandleFunc("/", home)
	http.HandleFunc("/register/", register)
	http.HandleFunc("/service/address/", serviceAddress)
	http.HandleFunc("/service/", proxyService)
	http.ListenAndServe(":8500", nil)

}

type Service struct {
	Id   string `json:"id"`
	Addr string `json:"address"`
	Port string `json:"port"`
}

func proxyService(w http.ResponseWriter, req *http.Request) {
	urlpath := req.URL.Path[len("/service/"):]
	path := strings.Split(urlpath, "/")[1:]
	serviceId := strings.Split(urlpath, "/")[0]

	service := services[serviceId][0]

	serviceAddr := fmt.Sprintf("http://%v:%v", service.Addr, service.Port)
	if path != nil {
		serviceAddr = fmt.Sprintf("%v/%v", serviceAddr, strings.Join(path, "/"))
	}

	remote, err := url.Parse(serviceAddr)
	if err != nil {
		internalServerError(w, err)
	}

	fmt.Printf("forwarding to %v\n", serviceAddr)

	proxy := httputil.NewSingleHostReverseProxy(remote)
	w.Header().Set("X-Ben", "Rad")
	proxy.ServeHTTP(w, req)
}

func serviceAddress(w http.ResponseWriter, req *http.Request) {
	if req.Method != "GET" {
		methodNotAllowed(w)
		return
	}
	serviceId := req.URL.Path[len("/service/address/"):]
	for _, service := range services[serviceId] {
		fmt.Fprintf(w, "%v:%v\n", service.Addr, service.Port)
	}

}

func register(w http.ResponseWriter, req *http.Request) {
	if req.Method != "POST" {
		methodNotAllowed(w)
		return
	}

	bodyBytes, err := ioutil.ReadAll(req.Body)
	defer req.Body.Close()

	if err != nil {

		return
	}

	ct := req.Header.Get("Content-Type")

	if ct != "application/json" {
		w.WriteHeader(http.StatusUnsupportedMediaType)
		w.Write([]byte(fmt.Sprintf("need content-type 'application/json', but got '%s'", ct)))
		return
	}

	var service Service

	err = json.Unmarshal(bodyBytes, &service)

	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
		return
	}

	if service.Addr == "" {
		result := strings.Split(req.RemoteAddr, ":")
		service.Addr = result[0]
		service.Port = result[1]

	}

	newSlice := append(services[service.Id], service)
	services[service.Id] = newSlice
	fmt.Printf("%v registered with address %v:%v\n", service.Id, service.Addr, service.Port)

	w.WriteHeader(http.StatusCreated)
}

func home(w http.ResponseWriter, req *http.Request) {
	fmt.Fprintf(w, "/resgister to register a new service\n")
	fmt.Fprintf(w, "/service/address to get a service address\n")
	fmt.Fprintf(w, "Services registered\n")

	for serviceId := range services {
		for _, service := range services[serviceId] {
			fmt.Fprintf(w, "%v: %v:%v\n", service.Id, service.Addr, service.Port)
		}
	}
}

func checkServicesHealth() {
	ticker := time.NewTicker(60 * time.Second)
	quit := make(chan struct{})

	go func() {
		for {
			select {
			case <-ticker.C:
				for sname, s := range services {
					for i, service := range s {
						if !isServiceUp(service) {
							if len(s) == 1 {
								services[sname] = []Service{}
								continue
							}
							s[i] = s[len(s)-1]
							s = s[:len(s)-1]
						}
					}
				}
			case <-quit:
				ticker.Stop()
				return
			}
		}
	}()
}

func isServiceUp(service Service) bool {

	fullAddr := fmt.Sprintf("http://%v:%v/health", service.Addr, service.Port)

	resp, err := http.Get(fullAddr)

	if err != nil {
		fmt.Errorf(err.Error())
	}

	return resp != nil && resp.StatusCode == http.StatusOK
}

func methodNotAllowed(w http.ResponseWriter) {
	w.WriteHeader(http.StatusMethodNotAllowed)
	w.Write([]byte("method not allowed"))
}

func internalServerError(w http.ResponseWriter, err error) {
	w.WriteHeader(http.StatusInternalServerError)
	w.Write([]byte(err.Error()))
}
