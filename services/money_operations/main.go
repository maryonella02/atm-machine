package main

import (
	"atm-machine/services/money_operations/service"
	"database/sql"
	"flag"
	"fmt"
	"log"

	_ "github.com/go-sql-driver/mysql"
)

var (
	addr       = flag.String("addr", ":8080", "address to listen on")
	dbUser     = flag.String("dbuser", "root", "database user")
	dbPassword = flag.String("dbpassword", "password", "database password")
	dbHost     = flag.String("dbhost", "localhost", "database host")
	dbPort     = flag.String("dbport", "3306", "database port")
	dbName     = flag.String("dbname", "atm", "database name")
)

func main() {
	// Open database connection
	db, err := sql.Open("mysql", fmt.Sprintf("%s:%s@tcp(%s:%s)/%s", *dbUser, *dbPassword, *dbHost, *dbPort, *dbName))
	if err != nil {
		log.Fatalf("failed to connect to database: %v", err)
	}
	defer db.Close()

	// Create new MoneyOperations service
	m := &service.MoneyOperations{
		DB: db, // Add database connection to service
	}

	// Start RPC server
	if err := m.Start(); err != nil {
		log.Fatal(err)
	}
}
