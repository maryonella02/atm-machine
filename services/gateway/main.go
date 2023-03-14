package main

import (
	"atm-machine/cache"
	"atm-machine/services/gateway/service"
	"fmt"
	"log"
	"net"
	"net/rpc"
	"time"
)

func main() {
	lruCache := cache.NewLRUCache(1000 * time.Hour)

	// Create a Gateway instance and register it as an RPC server.
	gateway := &service.Gateway{
		DiscoveryAddr: "atm-machine:8091",
		Cache:         lruCache,
	}
	rpc.Register(gateway)

	listener, err := net.Listen("tcp", ":8090")
	fmt.Printf("%s", listener.Addr().String())

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
