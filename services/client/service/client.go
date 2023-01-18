package service

import (
	"atm-machine/services/gateway/service"
	"context"
	"net/rpc"
)

type Client struct {
	Addr string
}

func (c *Client) Authenticate(ctx context.Context, cardNumber string, pin string) (string, error) {
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
		return "", err
	}
	println("authenticated", res.Token)
	return res.Token, nil
}

func (c *Client) GetBalance(ctx context.Context, cardNumber, token string) (int, error) {
	conn, err := rpc.Dial("tcp", c.Addr)
	if err != nil {
		return 0, err
	}
	defer conn.Close()

	req := service.GetBalanceRequest{
		CardNumber: cardNumber,

		Token: token,
	}
	res := service.GetBalanceResponse{}
	if err := conn.Call("Gateway.GetBalance", &req, &res); err != nil {
		return 0, err
	}
	return res.Balance, nil
}

func (c *Client) Withdraw(ctx context.Context, cardNumber, token string, amount int) (int, error) {
	conn, err := rpc.Dial("tcp", c.Addr)
	if err != nil {
		return 0, err
	}
	defer conn.Close()

	req := service.WithdrawRequest{
		Token:      token,
		CardNumber: cardNumber,

		Amount: amount,
	}
	res := service.WithdrawResponse{}
	if err := conn.Call("Gateway.Withdraw", &req, &res); err != nil {
		return 0, err
	}
	return res.Balance, nil
}
