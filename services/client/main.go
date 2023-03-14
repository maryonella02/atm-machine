package main

import (
	"atm-machine/services/client/service"
	"fmt"
	"log"
	"math/rand"
	"time"
)

func main() {
	c := &service.Client{
		Addr: "gateway:8090",
	}
	time.Sleep(10 * time.Second)

	token, err := c.Authenticate("1234567890", "1234")
	if err != nil {
		log.Fatal(err)
	}

	rand.Seed(time.Now().UnixNano()) // Seed the random number generator with the current time

	for {
		balance, err := c.GetBalance("1234567890", token)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println("Balance:", balance, token)

		// Withdraw a random amount between 1 and 10
		amount := rand.Intn(10) + 1
		balance, err = c.Withdraw("1234567890", token, amount)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println("New balance:", balance)

		// Wait for a random amount of time before withdrawing again
		time.Sleep(time.Duration(rand.Intn(10)) * time.Second)
	}

}
