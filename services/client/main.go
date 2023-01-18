package main

import (
	"atm-machine/services/client/service"
	"context"
	"fmt"
	"log"
)

func main() {
	c := &service.Client{
		Addr: "localhost:50052",
	}

	ctx := context.Background()
	token, err := c.Authenticate(ctx, "1234567890", "1234")
	if err != nil {
		log.Fatal(err)
	}

	balance, err := c.GetBalance(ctx, "1234567890", token)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Balance:", balance, token)

	amount := 50
	balance, err = c.Withdraw(ctx, "1234567890", token, amount)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("New balance:", balance)
}
