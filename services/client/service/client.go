package service

import (
	"atm-machine/services/money_operations/service"
	"context"
	"net/rpc"
)

type Client struct {
	Addr string
}

func (c *Client) GetBalance(ctx context.Context) (int, error) {
	conn, err := rpc.Dial("tcp", c.Addr)
	if err != nil {
		return 0, err
	}
	defer conn.Close()

	req := service.GetBalanceRequest{}
	res := service.GetBalanceResponse{}
	if err := conn.Call("MoneyOperations.GetBalance", &req, &res); err != nil {
		return 0, err
	}
	return res.Balance, nil
}

func (c *Client) Withdraw(ctx context.Context, amount int) (int, error) {
	conn, err := rpc.Dial("tcp", c.Addr)
	if err != nil {
		return 0, err
	}
	defer conn.Close()

	req := service.WithdrawRequest{
		Amount: amount,
	}
	res := service.WithdrawResponse{}
	if err := conn.Call("MoneyOperations.Withdraw", &req, &res); err != nil {
		return 0, err
	}
	return res.Balance, nil
}
