package service

import (
	"atm-machine/services/discovery/service"
	"fmt"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"net"
	"net/http"
	"net/rpc"
	"time"
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

var totalWithdrawals = prometheus.NewCounter(
	prometheus.CounterOpts{
		Name: "money_operations_total_withdrawals",
		Help: "The total number of withdrawals made",
	},
)

var (
	getBalanceCounter = prometheus.NewCounter(
		prometheus.CounterOpts{
			Name: "money_operations_get_balance_requests_total",
			Help: "Total number of GetBalance requests",
		},
	)

	getBalanceDuration = prometheus.NewHistogram(
		prometheus.HistogramOpts{
			Name:    "money_operations_get_balance_request_duration_seconds",
			Help:    "Duration of GetBalance requests in seconds",
			Buckets: []float64{0.01, 0.1, 1, 10},
		},
	)
)

func init() {
	prometheus.MustRegister(getBalanceCounter)
	prometheus.MustRegister(getBalanceDuration)
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

	totalWithdrawals.Inc()
	m.Stats["TotalWithdrawals"]++

	res.Balance = m.Balance - req.Amount
	return nil
}

func (m *MoneyOperations) GetBalance(req *GetBalanceRequest, res *GetBalanceResponse) error {
	start := time.Now()

	res.Balance = m.Balance

	getBalanceCounter.Inc()
	getBalanceDuration.Observe(time.Since(start).Seconds())
	return nil
}

func (m *MoneyOperations) Start() error {

	// register the RPC methods
	rpc.Register(m)
	rpc.HandleHTTP()

	// Register Prometheus metrics
	prometheus.MustRegister(totalWithdrawals)

	// Start Prometheus metrics server
	go func() {
		http.Handle("/metrics", promhttp.Handler())
		http.ListenAndServe(":8080", nil)
	}()

	// start the RPC server
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
