package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"
)

type Handler struct {
}

func (h *Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	toSend := []byte("Hello")
	_, err := w.Write(toSend)
	log.Printf("Message: %s", toSend)
	if err != nil {
		log.Printf("error while writing on the body %s", err)
	}
}

func main() {

	// create a server
	myServer := &http.Server{
		// set the server address
		Addr: "127.0.0.1:8081",
		// define some specific configuration
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
		Handler:      &Handler{},
	}
	register()
	// launch the server
	log.Fatal(myServer.ListenAndServe())

}

func GetStatement() string {
	return "no money in the system yet"
}

type Service struct {
	Id   string `json:"id"`
	Addr string `json:"address"`
	Port string `json:"port"`
}

func register() {

	content, err := json.Marshal(Service{Id: "1"})
	resp, err := http.Post("http://localhost:8500/register/", "application/json", bytes.NewBuffer(content))

	if err != nil {
		fmt.Errorf(err.Error())
	}
	fmt.Printf("Response status : %s \n", resp.Status)
	fmt.Printf("Body : %s \n ", resp.Body)
}
