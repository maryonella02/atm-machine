package service_test

import (
	"atm-machine/services/discovery/service"
	"net"
	"net/rpc"
	"testing"

	"github.com/onsi/ginkgo"
	"github.com/onsi/gomega"
)

func TestDiscovery(t *testing.T) {
	gomega.RegisterFailHandler(ginkgo.Fail)
	ginkgo.RunSpecs(t, "Discovery Suite")
}

var _ = ginkgo.Describe("Discovery", func() {
	var (
		discovery *service.Discovery
		listener  net.Listener
	)

	ginkgo.BeforeEach(func() {
		var err error
		discovery = &service.Discovery{}
		listener, err = net.Listen("tcp", ":8081")
		gomega.Expect(err).NotTo(gomega.HaveOccurred())
		go rpc.Accept(listener)
	})

	ginkgo.AfterEach(func() {
		listener.Close()
	})

	ginkgo.It("should register a service", func() {
		req := &service.RegisterRequest{
			Service: &service.Service{
				Name: "test-service",
				Addr: "localhost:8080",
			},
		}
		var res service.RegisterResponse
		err := discovery.Register(req, &res)
		gomega.Expect(err).NotTo(gomega.HaveOccurred())
		gomega.Expect(discovery.Services).To(gomega.ContainElement(req.Service))
	})

	ginkgo.It("should list registered services", func() {
		discovery.Services = []*service.Service{
			{Name: "test-service-1", Addr: "localhost:8080"},
			{Name: "test-service-2", Addr: "localhost:8081"},
		}
		var res service.ListResponse
		err := discovery.List(&service.ListRequest{}, &res)
		gomega.Expect(err).NotTo(gomega.HaveOccurred())
		gomega.Expect(res.Services).To(gomega.Equal(discovery.Services))
	})
})
