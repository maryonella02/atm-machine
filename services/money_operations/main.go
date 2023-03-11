package main

import (
	"atm-machine/services/money_operations/service"
	"flag"
	"log"
)

var (
	addr = flag.String("addr", ":8080", "address to listen on")
)

func main() {
	m := &service.MoneyOperations{
		Balance: 1000,
	}
	if err := m.Start(); err != nil {
		log.Fatal(err)
	}
}
