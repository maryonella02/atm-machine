package main

import (
	"atm-machine/services/money_operations/service"
	"log"
)

func main() {
	m := &service.MoneyOperations{
		Balance: 1000,
	}
	if err := m.Start(); err != nil {
		log.Fatal(err)
	}
}
