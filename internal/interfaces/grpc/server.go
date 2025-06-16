package grpc

import (
	"log"
	"net"

	"github.com/matzxrr/ddd-lemonadestore/internal/interfaces/grpc/interceptors"
	"github.com/matzxrr/ddd-lemonadestore/internal/interfaces/grpc/services"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"

	customerPb "github.com/matzxrr/ddd-lemonadestore/internal/interfaces/grpc/pb/customer/v1"
	orderPb "github.com/matzxrr/ddd-lemonadestore/internal/interfaces/grpc/pb/order/v1"
	storePb "github.com/matzxrr/ddd-lemonadestore/internal/interfaces/grpc/pb/store/v1"
)

// Server wraps the gRPC server
// WHY: Encapsulates server configuration and lifecycle
type Server struct {
    grpcServer      *grpc.Server
    storeService    *services.StoreService
    orderService    *services.OrderService
    customerService *services.CustomerService
}

// NewServer creates a new gRPC server
// WHERE: Created in main.go during application startup
func NewServer(
    storeService *services.StoreService,
    orderService *services.OrderService,
    customerService *services.CustomerService,
) *Server {
    // Create gRPC server with interceptors
    opts := []grpc.ServerOption{
        grpc.ChainUnaryInterceptor(
            interceptors.LoggingInterceptor,
            interceptors.ErrorInterceptor,
        ),
    }
    
    grpcServer := grpc.NewServer(opts...)
    
    // Register services
    storePb.RegisterStoreServiceServer(grpcServer, storeService)
    orderPb.RegisterOrderServiceServer(grpcServer, orderService)
    customerPb.RegisterCustomerServiceServer(grpcServer, customerService)
    
    // Enable reflection for development
    // WHAT: Allows tools like grpcurl to discover services
    reflection.Register(grpcServer)
    
    return &Server{
        grpcServer:      grpcServer,
        storeService:    storeService,
        orderService:    orderService,
        customerService: customerService,
    }
}

// Start starts the gRPC server
func (s *Server) Start(address string) error {
    lis, err := net.Listen("tcp", address)
    if err != nil {
        return err
    }
    
    log.Printf("gRPC server starting on %s", address)
    return s.grpcServer.Serve(lis)
}

// Stop gracefully stops the server
func (s *Server) Stop() {
    log.Println("Shutting down gRPC server...")
    s.grpcServer.GracefulStop()
}
