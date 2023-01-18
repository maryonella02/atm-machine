package main

import (
	"atm-machine/cache"
	"errors"
	"log"
	"net"
	"net/rpc"
)

type WithdrawRequest struct {
	CardNumber string
	Amount     int
}

type WithdrawResponse struct {
	Balance int
}

type GetBalanceRequest struct {
	CardNumber string
}

type GetBalanceResponse struct {
	Balance int
}

// Gateway is a service that provides access to the ATM system.
type Gateway struct {
	DiscoveryAddr string
	Cache         *cache.Cache
}

func (g *Gateway) Withdraw(req *WithdrawRequest, res *WithdrawResponse) error {
	// Call the money operations service to perform the withdrawal.
	client, err := g.getMoneyOperationsClient()
	if err != nil {
		return err
	}
	err = client.Call("MoneyOperations.Withdraw", req, res)
	if err != nil {
		return err
	}

	// Update the cache with the new balance.
	g.Cache.Set(req.CardNumber, res.Balance)
	return nil
}

func (g *Gateway) GetBalance(req *GetBalanceRequest, res *GetBalanceResponse) error {
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
	// Call the discovery service to find the address of the money operations service.
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
func main() {
	lruCache := cache.NewLRUCache(1000)

	// Create a Gateway instance and register it as an RPC server.
	gateway := &Gateway{
		DiscoveryAddr: "localhost:50051",
		Cache:         lruCache,
	}
	rpc.Register(gateway)

	listener, err := net.Listen("tcp", ":50052")
	if err != nil {
		log.Fatal(err)
	}
	defer listener.Close()

	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Println(err)
			continue
		}
		go rpc.ServeConn(conn)
	}
}
