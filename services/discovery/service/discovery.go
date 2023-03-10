package service

import (
	"context"
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

func (d *Discovery) Register(ctx context.Context, req *RegisterRequest, res *RegisterResponse) error {
	select {
	case <-ctx.Done():
		return ctx.Err()
	default:
		if err := ValidateServiceData(req.Service); err != nil {
			return err
		}
		if !serviceExists(req.Service, d.Services) {
			d.Services = append(d.Services, req.Service)
			return nil
		}
		return fmt.Errorf("service already registered")
	}
}

// List lists the services registered with the discovery service.
func (d *Discovery) List(ctx context.Context, req *ListRequest, res *ListResponse) error {
	select {
	case <-ctx.Done():
		return ctx.Err()
	default:
		res.Services = d.Services
		return nil
	}
}

func (d *Discovery) Start() error {
	rpc.Register(d)
	ln, err := net.Listen("tcp", ":8081")
	if err != nil {
		return err
	}
	d.ln = ln
	go func() {
		for {
			conn, err := ln.Accept()
			if err != nil {
				log.Println(err)
				continue
			}
			go rpc.ServeConn(conn)
		}
	}()
	return nil
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
