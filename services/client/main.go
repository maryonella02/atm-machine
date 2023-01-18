package main

import (
	"atm-machine/services/client/service"
	"context"
	"fmt"
	"log"
)

// var (
//
//	addr = flag.String("addr", ":8080", "address of gateway")
//
// )
func main() {
	c := &service.Client{
		Addr: "localhost:50052",
	}

	ctx := context.Background()
	balance, err := c.GetBalance(ctx)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Balance:", balance)

	amount := 50
	balance, err = c.Withdraw(ctx, amount)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("New balance:", balance)
}
