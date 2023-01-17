package service

import (
	"context"
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
	d.Services = append(d.Services, req.Service)
	return nil
}

// List lists the services registered with the discovery service.
func (d *Discovery) List(ctx context.Context, req *ListRequest, res *ListResponse) error {
	res.Services = d.Services
	return nil
}

func (d *Discovery) Start() error {
	rpc.Register(d)
	ln, err := net.Listen("tcp", ":8081")
	if err != nil {
		return err
	}
	rpc.Accept(ln)
	return nil

}
