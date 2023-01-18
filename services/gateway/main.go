package main

import (
	"atm-machine/cache"
	"atm-machine/services/gateway/service"
	"log"
	"net"
	"net/rpc"
	"time"
)

func main() {
	lruCache := cache.NewLRUCache(1000 * time.Hour)

	// Create a Gateway instance and register it as an RPC server.
	gateway := &service.Gateway{
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
