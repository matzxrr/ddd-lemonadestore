package services

import (
	"context"

	"github.com/matzxrr/ddd-lemonadestore/internal/application/customer/commands"
	"github.com/matzxrr/ddd-lemonadestore/internal/application/customer/queries"
	"github.com/matzxrr/ddd-lemonadestore/internal/application/dtos"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"
	pb "github.com/matzxrr/ddd-lemonadestore/internal/interfaces/grpc/pb/customer/v1"
)

// CustomerService implements the gRPC CustomerService
type CustomerService struct {
    pb.UnimplementedCustomerServiceServer
    
    // Command handlers
    registerCustomerHandler *commands.RegisterCustomerHandler
    updateCustomerHandler   *commands.UpdateCustomerHandler
    
    // Query handlers
    getCustomerHandler *queries.GetCustomerHandler
}

// NewCustomerService creates a new customer service
func NewCustomerService(
    registerCustomer *commands.RegisterCustomerHandler,
    updateCustomer *commands.UpdateCustomerHandler,
    getCustomer *queries.GetCustomerHandler,
) *CustomerService {
    return &CustomerService{
        registerCustomerHandler: registerCustomer,
        updateCustomerHandler:   updateCustomer,
        getCustomerHandler:      getCustomer,
    }
}

// RegisterCustomer registers a new customer
func (s *CustomerService) RegisterCustomer(
    ctx context.Context,
    req *pb.RegisterCustomerRequest,
) (*pb.RegisterCustomerResponse, error) {
    // Validate request
    if req.Email == "" || req.FirstName == "" || req.LastName == "" {
        return nil, status.Error(codes.InvalidArgument, "email, first_name, and last_name are required")
    }
    
    // Create command
    cmd := commands.RegisterCustomerCommand{
        Email:     req.Email,
        FirstName: req.FirstName,
        LastName:  req.LastName,
    }
    
    // Execute command
    customerDTO, err := s.registerCustomerHandler.Handle(ctx, cmd)
    if err != nil {
        return nil, toGRPCError(err)
    }
    
    return &pb.RegisterCustomerResponse{
        CustomerId: customerDTO.ID,
        Type:       customerDTO.Type,
    }, nil
}

// UpdateCustomerContact updates customer contact information
func (s *CustomerService) UpdateCustomerContact(
    ctx context.Context,
    req *pb.UpdateCustomerContactRequest,
) (*pb.UpdateCustomerContactResponse, error) {
    // Validate address
    if req.Address == nil {
        return nil, status.Error(codes.InvalidArgument, "address is required")
    }
    
    // Create command
    cmd := commands.UpdateCustomerCommand{
        CustomerID:  req.CustomerId,
        PhoneNumber: req.PhoneNumber,
        Address: dtos.AddressDTO{
            Street:  req.Address.Street,
            City:    req.Address.City,
            State:   req.Address.State,
            ZipCode: req.Address.ZipCode,
            Country: req.Address.Country,
        },
    }
    
    // Execute command
    err := s.updateCustomerHandler.Handle(ctx, cmd)
    if err != nil {
        return nil, toGRPCError(err)
    }
    
    return &pb.UpdateCustomerContactResponse{
        Success: true,
    }, nil
}

// GetCustomer retrieves customer details
func (s *CustomerService) GetCustomer(
    ctx context.Context,
    req *pb.GetCustomerRequest,
) (*pb.GetCustomerResponse, error) {
    // Create query
    query := queries.GetCustomerQuery{
        CustomerID: req.CustomerId,
    }
    
    // Execute query
    customerDTO, err := s.getCustomerHandler.Handle(ctx, query)
    if err != nil {
        return nil, toGRPCError(err)
    }
    
    // Convert to protobuf
    return &pb.GetCustomerResponse{
        Customer: &pb.Customer{
            Id:            customerDTO.ID,
            Email:         customerDTO.Email,
            FirstName:     customerDTO.FirstName,
            LastName:      customerDTO.LastName,
            PhoneNumber:   customerDTO.PhoneNumber,
            Type:          customerDTO.Type,
            LoyaltyPoints: int32(customerDTO.LoyaltyPoints),
            IsActive:      customerDTO.IsActive,
            RegisteredAt:  timestamppb.New(customerDTO.RegisteredAt),
        },
    }, nil
}
