package main

import (
	"log"
	"net/http"
	"time"
)

type Handler struct {
}

func (h *Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	toSend := []byte("Hello")
	_, err := w.Write(toSend)
	if err != nil {
		log.Printf("error while writing on the body %s", err)
	}
}

func main() {
	// create a server
	myServer := &http.Server{
		// set the server address
		Addr: "127.0.0.1:8080",
		// define some specific configuration
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
		Handler:      &Handler{},
	}
	// launch the server
	log.Fatal(myServer.ListenAndServe())
}
