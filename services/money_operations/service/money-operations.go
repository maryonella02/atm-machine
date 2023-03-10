package service

import (
	"database/sql"
	"fmt"
	"net"
	"net/rpc"

	_ "github.com/go-sql-driver/mysql"
)

type MoneyOperations struct {
	DB     *sql.DB
	Port   string
	Status string
	Stats  map[string]int
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
	// start a transaction
	tx, err := m.DB.Begin()
	if err != nil {
		return err
	}

	// check if the balance is sufficient
	var balance int
	err = tx.QueryRow("SELECT balance FROM account").Scan(&balance)
	if err != nil {
		tx.Rollback()
		return err
	}
	if balance < req.Amount {
		tx.Rollback()
		return fmt.Errorf("insufficient funds")
	}

	// update the balance
	_, err = tx.Exec("UPDATE account SET balance = balance - ?", req.Amount)
	if err != nil {
		tx.Rollback()
		return err
	}

	// commit the transaction
	err = tx.Commit()
	if err != nil {
		return err
	}

	// update stats
	m.Stats["TotalWithdrawals"]++

	// return the new balance
	res.Balance = balance - req.Amount
	return nil
}

func (m *MoneyOperations) GetBalance(req *GetBalanceRequest, res *GetBalanceResponse) error {
	// retrieve the balance from the database
	err := m.DB.QueryRow("SELECT balance FROM account").Scan(&res.Balance)
	if err != nil {
		return err
	}
	return nil
}

func (m *MoneyOperations) Start() error {
	// initialize the database connection
	db, err := sql.Open("mysql", "user:password@tcp(host:port)/database")
	if err != nil {
		return err
	}
	m.DB = db

	// create tables if they don't exist
	_, err = db.Exec("CREATE TABLE IF NOT EXISTS account (balance INT)")
	if err != nil {
		return err
	}

	// register the RPC methods
	rpc.Register(m)
	rpc.HandleHTTP()

	// start the RPC server
	m.Port = "8080"
	m.Status = "Running"
	m.Stats = make(map[string]int)
	m.Stats["TotalWithdrawals"] = 0
	ln, err := net.Listen("tcp", ":"+m.Port)
	if err != nil {
		return err
	}
	rpc.Accept(ln)

	return nil
}
