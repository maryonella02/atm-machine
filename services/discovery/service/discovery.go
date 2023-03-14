package service

import (
	"fmt"
	"log"
	"net"
	"net/rpc"
)

// Service is a service that has been registered with the discovery service.
type Service struct {
	Name string
	Addr string
}

// Discovery is a service that allows other services to register themselves and be discovered by clients.
type Discovery struct {
	Services []*Service
	ln       net.Listener
}

type RegisterRequest struct {
	Service *Service
}

type RegisterResponse struct{}

type ListRequest struct{}

type ListResponse struct {
	Services []*Service
}

func (d *Discovery) Register(req *RegisterRequest, res *RegisterResponse) error {
	if err := ValidateServiceData(req.Service); err != nil {
		return err
	}
	if !serviceExists(req.Service, d.Services) {
		d.Services = append(d.Services, req.Service)
		log.Printf("Registered service name: %s, address: %s\n", req.Service.Name, req.Service.Addr)
		return nil
	}
	return fmt.Errorf("service already registered")
}

// List lists the services registered with the discovery service.
func (d *Discovery) List(req *ListRequest, res *ListResponse) error {
	res.Services = d.Services
	return nil
}

func (d *Discovery) Start() error {
	rpc.Register(d)
	ln, err := net.Listen("tcp", ":8091")
	if err != nil {
		return err
	}

	log.Printf("Discovery service started on %s\n", ln.Addr().String())
	d.ln = ln

	for {
		conn, err := ln.Accept()
		if err != nil {
			log.Println(err)
			continue

		}
		go rpc.ServeConn(conn)
	}
}

func (d *Discovery) Stop() error {
	if d.ln != nil {
		err := d.ln.Close()
		return err
	}
	return nil
}

func serviceExists(service *Service, services []*Service) bool {
	for _, s := range services {
		if s == service {
			return true
		}
	}
	return false
}

func ValidateServiceData(service *Service) error {
	if service.Name == "" {
		return fmt.Errorf("invalid service name")
	}
	if service.Addr == "" {
		return fmt.Errorf("invalid service address")
	}
	return nil
}
