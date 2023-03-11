package main

import (
	"atm-machine/services/client/service"
	"fmt"
	"log"
)

func main() {
	c := &service.Client{
		Addr: "localhost:8090",
	}

	token, err := c.Authenticate("1234567890", "1234")
	if err != nil {
		log.Fatal(err)
	}

	balance, err := c.GetBalance("1234567890", token)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Balance:", balance, token)

	amount := 50
	balance, err = c.Withdraw("1234567890", token, amount)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("New balance:", balance)
}
