package service

import (
	"atm-machine/cache"
	"errors"
	"fmt"
	"net/rpc"
)

type WithdrawRequest struct {
	CardNumber string
	Amount     int
	Token      string
}

type WithdrawResponse struct {
	Balance int
}

type GetBalanceRequest struct {
	CardNumber string
	Token      string
}

type GetBalanceResponse struct {
	Balance int
}
type AuthenticateRequest struct {
	CardNumber string
	Pin        string
}

type AuthenticateResponse struct {
	Token string
}

// Gateway is a service that provides access to the ATM system.
type Gateway struct {
	DiscoveryAddr string
	Cache         *cache.Cache
}

func (g *Gateway) Withdraw(req *WithdrawRequest, res *WithdrawResponse) error {
	if !g.isValidToken(req.CardNumber, req.Token) {
		return errors.New("invalid token")
	}
	client, err := g.getMoneyOperationsClient()
	if err != nil {
		return err
	}
	err = client.Call("MoneyOperations.Withdraw", req, res)
	if err != nil {
		return err
	}
	g.Cache.Set(req.CardNumber, res.Balance)
	return nil
}

func (g *Gateway) GetBalance(req *GetBalanceRequest, res *GetBalanceResponse) error {
	fmt.Println(req.CardNumber, req.Token)
	if !g.isValidToken(req.CardNumber, req.Token) {
		return errors.New("invalid token")
	}
	val, ok := g.Cache.Get(req.CardNumber)
	balance := val.(int)
	if ok {
		res.Balance = balance
		return nil
	}

	client, err := g.getMoneyOperationsClient()
	if err != nil {
		return err
	}
	err = client.Call("MoneyOperations.GetBalance", req, res)
	if err != nil {
		return err
	}

	g.Cache.Set(req.CardNumber, res.Balance)
	return nil
}
func (g *Gateway) Authenticate(req *AuthenticateRequest, res *AuthenticateResponse) error {
	// Check the card number and pin against a database or another service.
	// In this example, we're just hardcoding a card number and pin for demonstration purposes.
	if req.CardNumber == "1234567890" && req.Pin == "1234" {
		res.Token = "abcdefghijklmnopqrstuvwxyz"
		g.Cache.Set(req.CardNumber, res.Token)
		return nil
	}
	return errors.New("invalid card number or pin")
}

// isValidToken is a helper function that checks if the provided token is valid for the given card number.
func (g *Gateway) isValidToken(cardNumber, token string) bool {
	fmt.Println("validation", token, cardNumber)

	val, ok := g.Cache.Get(cardNumber)
	if !ok {
		return false
	}
	cachedToken := val
	fmt.Println(token, cachedToken)
	if token == cachedToken {
		return true
	}
	return false
}

// ListResponse is a response to a request to list the available services.
type ListResponse struct {
	Addrs []string
}

type StatusResponse struct {
	Status string
	Port   string
	Stats  map[string]int
}

// getMoneyOperationsClient is a helper function that calls the discovery service to find the
// address of the money operations service and then returns an RPC client that can be used to call
// methods on that service.
func (g *Gateway) getMoneyOperationsClient() (*rpc.Client, error) {
	client, err := rpc.Dial("tcp", g.DiscoveryAddr)
	if err != nil {
		return nil, err
	}
	var res ListResponse
	err = client.Call("Discovery.List", &struct{}{}, &res)
	if err != nil {
		return nil, err
	}
	if len(res.Addrs) == 0 {
		return nil, errors.New("money operations service not found")
	}

	client, err = rpc.Dial("tcp", res.Addrs[0])
	if err != nil {
		return nil, err
	}

	// Call the status endpoint to check the status of the money operations service
	var status StatusResponse
	err = client.Call("MoneyOperations.Status", &struct{}{}, &status)
	if err != nil {
		return nil, err
	}
	if status.Status != "OK" {
		return nil, errors.New("money operations service is not available")
	}

	return client, nil
}
