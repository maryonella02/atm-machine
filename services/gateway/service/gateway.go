package service

import (
	"atm-machine/cache"
	"atm-machine/services/discovery/service"
	moneyService "atm-machine/services/money_operations/service"
	"errors"
	"log"
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
	log.Printf("[Withdraw] CardNumber=%s, Amount=%d, Token=%s\n", req.CardNumber, req.Amount, req.Token)
	if !g.isValidToken(req.CardNumber, req.Token) {
		log.Printf("[Withdraw] Invalid token for CardNumber=%s, Token=%s\n", req.CardNumber, req.Token)
		return errors.New("invalid token")
	}
	log.Println("[Withdraw] Calling MoneyOperations.Withdraw")
	client, err := g.getMoneyOperationsClient()
	if err != nil {
		return err
	}
	err = client.Call("MoneyOperations.Withdraw", req, res)
	if err != nil {
		log.Printf("[Withdraw] Error calling MoneyOperations.Withdraw: %v\n", err)
		return err
	}
	log.Printf("[Withdraw] MoneyOperations.Withdraw response Balance=%d\n", res.Balance)
	g.Cache.Set(req.Token, res.Balance)
	return nil
}

func (g *Gateway) GetBalance(req *GetBalanceRequest, res *GetBalanceResponse) error {
	log.Printf("[GetBalance] CardNumber=%s, Token=%s\n", req.CardNumber, req.Token)
	if !g.isValidToken(req.CardNumber, req.Token) {
		log.Printf("[GetBalance] Invalid token for CardNumber=%s, Token=%s\n", req.CardNumber, req.Token)
		return errors.New("invalid token")
	}
	val, ok := g.Cache.Get(req.Token)
	balance := val.(int)
	if ok {
		log.Printf("[GetBalance] Returning cached balance=%d for Token=%s\n", balance, req.Token)
		res.Balance = balance
		return nil
	}

	log.Println("[GetBalance] Calling MoneyOperations.GetBalance")
	client, err := g.getMoneyOperationsClient()
	if err != nil {
		return err
	}
	err = client.Call("MoneyOperations.GetBalance", req, res)
	if err != nil {
		log.Printf("[GetBalance] Error calling MoneyOperations.GetBalance: %v\n", err)
		return err
	}
	log.Printf("[GetBalance] MoneyOperations.GetBalance response Balance=%d\n", res.Balance)
	g.Cache.Set(req.Token, res.Balance)
	return nil
}
func (g *Gateway) Authenticate(req *AuthenticateRequest, res *AuthenticateResponse) error {
	log.Printf("Authenticating card %s with pin %s", req.CardNumber, req.Pin)

	// Check the card number and pin against a database or another service.
	// In this example, we're just hardcoding a card number and pin for demonstration purposes.
	if req.CardNumber == "1234567890" && req.Pin == "1234" {
		res.Token = "abcdefghijklmnopqrstuvwxyz"
		g.Cache.Set(req.CardNumber, res.Token)
		g.Cache.Set(res.Token, 1000)

		log.Printf("Card %s successfully authenticated, generated token: %s", req.CardNumber, res.Token)
		return nil
	}

	log.Printf("Failed to authenticate card %s with pin %s", req.CardNumber, req.Pin)
	return errors.New("invalid card number or pin")
}

// isValidToken is a helper function that checks if the provided token is valid for the given card number.
func (g *Gateway) isValidToken(cardNumber, token string) bool {
	log.Printf("Validating token %s for card %s", token, cardNumber)

	val, ok := g.Cache.Get(cardNumber)
	if !ok {
		log.Printf("Card %s not found in cache", cardNumber)
		return false
	}
	cachedToken := val.(string)
	log.Printf("Cached token for card %s: %s", cardNumber, cachedToken)

	if token == cachedToken {
		log.Printf("Token %s for card %s is valid", token, cardNumber)
		return true
	}

	log.Printf("Token %s for card %s is invalid", token, cardNumber)
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
	log.Printf("Trying to dial discovery service at %s", g.DiscoveryAddr)
	client, err := rpc.Dial("tcp", g.DiscoveryAddr)
	if err != nil {
		log.Printf("Failed to dial discovery service at %s: %v", g.DiscoveryAddr, err)
		return nil, err
	}
	var res = &service.ListResponse{}
	var req service.ListRequest
	err = client.Call("Discovery.List", &req, res)
	if err != nil {
		log.Printf("Failed to call Discovery.List on %s: %v", g.DiscoveryAddr, err)
		return nil, err
	}

	if len(res.Services) == 0 {
		log.Printf("No money operations service found")
		return nil, errors.New("money operations service not found")
	}

	log.Printf("Dialing money operations service at %s", res.Services[0].Addr)
	client, err = rpc.Dial("tcp", "money-operations:"+res.Services[0].Addr)
	if err != nil {
		log.Printf("Failed to dial money operations service at %s: %v", res.Services[0].Addr, err)
		return nil, err
	}

	// Call the status endpoint to check the status of the money operations service
	var status StatusResponse
	err = client.Call("MoneyOperations.ServiceStatus", moneyService.ServiceStatusRequest{}, &status)
	if err != nil {
		log.Printf("Failed to call MoneyOperations.ServiceStatus on money operations service: %v", err)
		return nil, err
	}

	if status.Status != "Running" {
		log.Printf("Money operations service status is not Running")
		return nil, errors.New("money operations service is not available")
	}

	log.Printf("Connected to money operations service")
	return client, nil

}
