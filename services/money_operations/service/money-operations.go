package service

import (
	"atm-machine/services/discovery/service"
	"fmt"
	"net"
	"net/rpc"
)

type MoneyOperations struct {
	Balance int
	Port    string
	Status  string
	Stats   map[string]int
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

type ServiceStatusRequest struct{}

type ServiceStatusResponse struct {
	Status string
	Port   string
	Stats  map[string]int
}

func (m *MoneyOperations) ServiceStatus(req *ServiceStatusRequest, res *ServiceStatusResponse) error {
	res.Status = m.Status
	res.Port = m.Port
	res.Stats = m.Stats
	return nil
}

func (m *MoneyOperations) Withdraw(req *WithdrawRequest, res *WithdrawResponse) error {
	if m.Balance < req.Amount {
		return fmt.Errorf("insufficient funds")
	}
	m.Balance -= req.Amount
	res.Balance = m.Balance

	m.Stats["TotalWithdrawals"]++
	return nil
}

func (m *MoneyOperations) GetBalance(req *GetBalanceRequest, res *GetBalanceResponse) error {
	res.Balance = m.Balance
	return nil
}

func (m *MoneyOperations) Start() error {
	rpc.Register(m)
	rpc.HandleHTTP()
	m.Port = "8092"
	m.Status = "Running"
	m.Stats = make(map[string]int)
	m.Stats["TotalWithdrawals"] = 0
	ln, err := net.Listen("tcp", ":"+m.Port)
	fmt.Printf("%s", ln.Addr().String())

	if err != nil {
		return err
	}

	// Create an RPC client to connect to the discovery service
	client, err := rpc.Dial("tcp", "localhost:8091")
	if err != nil {
		return err
	}

	err = client.Call("Discovery.Register", service.RegisterRequest{Service: &service.Service{Name: "MoneyOperations", Addr: m.Port}}, &service.RegisterResponse{})
	if err != nil {
		return err
	}

	rpc.Accept(ln)
	return nil
}
