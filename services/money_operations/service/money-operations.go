package service

import (
	"fmt"
	"net"
	"net/rpc"
)

type MoneyOperations struct {
	Balance int
}

type WithdrawRequest struct {
	Amount int
}

type WithdrawResponse struct {
	Balance int
}

type GetBalanceRequest struct{}

type GetBalanceResponse struct {
	Balance int
}

func (m *MoneyOperations) Withdraw(req *WithdrawRequest, res *WithdrawResponse) error {
	if m.Balance < req.Amount {
		return fmt.Errorf("insufficient funds")
	}
	m.Balance -= req.Amount
	res.Balance = m.Balance
	return nil
}

func (m *MoneyOperations) GetBalance(req *GetBalanceRequest, res *GetBalanceResponse) error {
	res.Balance = m.Balance
	return nil
}

func (m *MoneyOperations) Start() error {
	rpc.Register(m)
	ln, err := net.Listen("tcp", ":8080")
	if err != nil {
		return err
	}
	rpc.Accept(ln)
	return nil
}
