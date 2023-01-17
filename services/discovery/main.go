package main

import (
	"atm-machine/services/discovery/service"
	"log"
)

func main() {
	d := &service.Discovery{}
	if err := d.Start(); err != nil {
		log.Fatal(err)
	}
}
