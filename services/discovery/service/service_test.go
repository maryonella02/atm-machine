package service_test

import (
	"atm-machine/services/discovery/service"
	"context"
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
			Expect(err).NotTo(HaveOccurred())
			go rpc.Accept(listener)
			go discovery.Start()
		})

		AfterEach(func() {
			time.Sleep(1 * time.Second)
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
			err := discovery.Register(context.Background(), req, &res)
			Expect(err).NotTo(HaveOccurred())
			Expect(discovery.Services).To(ContainElement(req.Service))
		})

		It("should list registered services", func() {
			discovery.Services = []*service.Service{
				{Name: "test-service-1", Addr: "localhost:8080"},
				{Name: "test-service-2", Addr: "localhost:8081"},
			}
			var res service.ListResponse
			err := discovery.List(context.Background(), &service.ListRequest{}, &res)
			Expect(err).NotTo(HaveOccurred())
			Expect(res.Services).To(Equal(discovery.Services))
		})
		Context("when registering a service with invalid input", func() {
			It("should return an error", func() {

				// Try to register a service with an empty name
				req := &service.RegisterRequest{Service: &service.Service{Name: "", Addr: "localhost:8000"}}
				var res service.RegisterResponse
				err := discovery.Register(context.Background(), req, &res)

				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(Equal("invalid service name"))

				// Try to register a service with an empty address
				req = &service.RegisterRequest{Service: &service.Service{Name: "test-service", Addr: ""}}
				err = discovery.Register(context.Background(), req, &res)

				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(Equal("invalid service address"))
			})
		})
		Context("when registering a service that is already registered", func() {
			It("should return an error", func() {
				req := &service.RegisterRequest{Service: &service.Service{Name: "test-service", Addr: "localhost:8000"}}
				var res service.RegisterResponse

				err := discovery.Register(context.Background(), req, &res)
				Expect(err).NotTo(HaveOccurred())

				// Try to register the same service again
				err = discovery.Register(context.Background(), req, &res)

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

	When("start&stop not needed", func() {
		Context("when the service is not running", func() {
			It("should return an error when trying to register a service", func() {
				discovery := &service.Discovery{}
				req := &service.RegisterRequest{Service: &service.Service{Name: "test-service", Addr: "localhost:8000"}}
				var res service.RegisterResponse

				ctx, cancel := context.WithTimeout(context.Background(), time.Nanosecond)
				defer cancel()

				err := discovery.Register(ctx, req, &res)

				Expect(err).To(HaveOccurred())
				Expect(err).To(Equal(context.DeadlineExceeded))
			})
			It("should return an error when trying to list services", func() {
				discovery := &service.Discovery{}
				req := &service.ListRequest{}
				var res service.ListResponse

				ctx, cancel := context.WithTimeout(context.Background(), time.Nanosecond)
				defer cancel()

				err := discovery.List(ctx, req, &res)

				Expect(err).To(HaveOccurred())
				Expect(err).To(Equal(context.DeadlineExceeded))
			})
		})
	})

	//Context("when multiple requests are made concurrently", func() {
	//	It("should handle the requests correctly", func() {
	//		// Start the Discovery service
	//		discovery := &service.Discovery{}
	//		go discovery.Start()
	//		defer discovery.Stop()
	//		// Create a wait group to track the goroutines
	//		wg := &sync.WaitGroup{}
	//		wg.Add(10)
	//		// Make 10 concurrent requests to register services
	//		for i := 0; i < 10; i++ {
	//			go func(i int) {
	//				defer wg.Done()
	//				req := &service.RegisterRequest{Service: &service.Service{Name: fmt.Sprintf("test-service-%d", i), Addr: "localhost:8000"}}
	//				var res service.RegisterResponse
	//				err := discovery.Register(context.Background(), req, &res)
	//				Expect(err).NotTo(HaveOccurred())
	//			}(i)
	//		}
	//		// Wait for all the requests to complete
	//		wg.Wait()
	//		// Check that all the services were registered correctly
	//		req := &service.ListRequest{}
	//		var res service.ListResponse
	//		err := discovery.List(context.Background(), req, &res)
	//		Expect(err).NotTo(HaveOccurred())
	//		Expect(res.Services).To(HaveLen(10))
	//	})
	//})
	//Context("when the maximum number of services are registered", func() {
	//	It("should handle the requests correctly", func() {
	//		// Start the Discovery service
	//		discovery := &service.Discovery{}
	//		go discovery.Start()
	//		defer discovery.Stop()
	//		// Create a wait group to track the goroutines
	//		wg := &sync.WaitGroup{}
	//		wg.Add(1000)
	//		// Make 1000 requests to register services
	//		for i := 0; i < 1000; i++ {
	//			go func(i int) {
	//				defer wg.Done()
	//				req := &service.RegisterRequest{Service: &service.Service{Name: fmt.Sprintf("test-service-%d", i), Addr: "localhost:8000"}}
	//				var res service.RegisterResponse
	//				err := discovery.Register(context.Background(), req, &res)
	//				Expect(err).NotTo(HaveOccurred())
	//			}(i)
	//		}
	//		// Wait for all the requests to complete
	//		wg.Wait()
	//		// Try to register one more service
	//		req := &service.RegisterRequest{Service: &service.Service{Name: "test-service-1000", Addr: "localhost:8000"}}
	//		var res service.RegisterResponse
	//		err := discovery.Register(context.Background(), req, &res)
	//		Expect(err).To(HaveOccurred())
	//	})
	//})

})
