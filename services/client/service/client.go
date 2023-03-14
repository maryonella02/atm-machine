package service

import (
	"atm-machine/services/gateway/service"
	"log"
	"net/rpc"
)

type Client struct {
	Addr string
}

func (c *Client) Authenticate(cardNumber string, pin string) (string, error) {
	conn, err := rpc.Dial("tcp", c.Addr)
	if err != nil {
		return "", err
	}
	defer conn.Close()

	req := service.AuthenticateRequest{
		CardNumber: cardNumber,
		Pin:        pin,
	}
	res := service.AuthenticateResponse{}
	if err := conn.Call("Gateway.Authenticate", &req, &res); err != nil {
		log.Printf("Error in Authenticate RPC call: %v", err)
		return "", err
	}
	log.Printf("Authenticated: cardNumber=%v, token=%v", cardNumber, res.Token)
	return res.Token, nil
}

func (c *Client) GetBalance(cardNumber, token string) (int, error) {
	conn, err := rpc.Dial("tcp", c.Addr)
	if err != nil {
		return 0, err
	}
	defer conn.Close()

	req := service.GetBalanceRequest{
		CardNumber: cardNumber,
		Token:      token,
	}
	res := service.GetBalanceResponse{}
	if err := conn.Call("Gateway.GetBalance", &req, &res); err != nil {
		log.Printf("Error in GetBalance RPC call: %v", err)
		return 0, err
	}
	log.Printf("Got balance: cardNumber=%v, balance=%v", cardNumber, res.Balance)
	return res.Balance, nil
}

func (c *Client) Withdraw(cardNumber, token string, amount int) (int, error) {
	conn, err := rpc.Dial("tcp", c.Addr)
	if err != nil {
		return 0, err
	}
	defer conn.Close()

	req := service.WithdrawRequest{
		Token:      token,
		CardNumber: cardNumber,
		Amount:     amount,
	}
	res := service.WithdrawResponse{}
	if err := conn.Call("Gateway.Withdraw", &req, &res); err != nil {
		log.Printf("Error in Withdraw RPC call: %v", err)
		return 0, err
	}
	log.Printf("Withdrew amount: cardNumber=%v, amount=%v, newBalance=%v", cardNumber, amount, res.Balance)
	return res.Balance, nil
}
