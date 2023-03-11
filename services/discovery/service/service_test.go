package service_test

import (
	"atm-machine/services/discovery/service"
	"fmt"
	"net"
	"net/rpc"
	"testing"
	"time"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func TestDiscovery(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Discovery Suite")
}

var _ = Describe("Discovery", func() {
	var (
		discovery *service.Discovery
		listener  net.Listener
	)

	When("start&stop needed", func() {
		BeforeEach(func() {
			var err error
			discovery = &service.Discovery{}

			listener, err = net.Listen("tcp", ":8081")
			fmt.Printf("%s", listener.Addr().String())

			Expect(err).NotTo(HaveOccurred())
			go rpc.Accept(listener)
			go discovery.Start()
		})

		AfterEach(func() {
			time.Sleep(30 * time.Millisecond)
			discovery.Stop()
			listener.Close()
		})
		It("should register a service", func() {
			req := &service.RegisterRequest{
				Service: &service.Service{
					Name: "test-service",
					Addr: "localhost:8080",
				},
			}
			var res service.RegisterResponse
			err := discovery.Register(req, &res)
			Expect(err).NotTo(HaveOccurred())
			Expect(discovery.Services).To(ContainElement(req.Service))
		})

		It("should list registered services", func() {
			discovery.Services = []*service.Service{
				{Name: "test-service-1", Addr: "localhost:8080"},
				{Name: "test-service-2", Addr: "localhost:8081"},
			}
			var res service.ListResponse
			discovery.List(&service.ListRequest{}, &res)
			Expect(res.Services).To(Equal(discovery.Services))
		})
		Context("when registering a service with invalid input", func() {
			It("should return an error", func() {

				// Try to register a service with an empty name
				req := &service.RegisterRequest{Service: &service.Service{Name: "", Addr: "localhost:8000"}}
				var res service.RegisterResponse
				err := discovery.Register(req, &res)

				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(Equal("invalid service name"))

				// Try to register a service with an empty address
				req = &service.RegisterRequest{Service: &service.Service{Name: "test-service", Addr: ""}}
				err = discovery.Register(req, &res)

				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(Equal("invalid service address"))
			})
		})
		Context("when registering a service that is already registered", func() {
			It("should return an error", func() {
				req := &service.RegisterRequest{Service: &service.Service{Name: "test-service", Addr: "localhost:8000"}}
				var res service.RegisterResponse

				err := discovery.Register(req, &res)
				Expect(err).NotTo(HaveOccurred())

				// Try to register the same service again
				err = discovery.Register(req, &res)

				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(Equal("service already registered"))
			})
		})
		Context("when starting the service with an invalid listening address", func() {
			It("should return an error", func() {
				err := discovery.Start()
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("listen tcp :8081: bind: address already in use"))
			})
		})
	})

})
