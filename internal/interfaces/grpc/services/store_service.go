package services

import (
	"context"

	"github.com/matzxrr/ddd-lemonadestore/internal/application/store/commands"
	"github.com/matzxrr/ddd-lemonadestore/internal/application/store/queries"
	"github.com/matzxrr/ddd-lemonadestore/internal/domain/store"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	pb "github.com/matzxrr/ddd-lemonadestore/internal/interfaces/grpc/pb/store/v1"
)

// StoreService implements the gRPC StoreService
// WHY: Exposes store operations via gRPC protocol
// WHERE: Registered with gRPC server in main.go
type StoreService struct {
    pb.UnimplementedStoreServiceServer
    
    // Command handlers
    addInventoryHandler *commands.AddInventoryHandler
    updatePriceHandler  *commands.UpdatePriceHandler
    
    // Query handlers
    getProductHandler   *queries.GetProductHandler
    getInventoryHandler *queries.GetInventoryHandler
}

// NewStoreService creates a new store service
func NewStoreService(
    addInventory *commands.AddInventoryHandler,
    updatePrice *commands.UpdatePriceHandler,
    getProduct *queries.GetProductHandler,
    getInventory *queries.GetInventoryHandler,
) *StoreService {
    return &StoreService{
        addInventoryHandler: addInventory,
        updatePriceHandler:  updatePrice,
        getProductHandler:   getProduct,
        getInventoryHandler: getInventory,
    }
}

// AddInventory adds inventory to a product
func (s *StoreService) AddInventory(
    ctx context.Context,
    req *pb.AddInventoryRequest,
) (*pb.AddInventoryResponse, error) {
    // Validate request
    if req.StoreId == "" || req.ProductId == "" {
        return nil, status.Error(codes.InvalidArgument, "store_id and product_id are required")
    }
    
    if req.Quantity <= 0 {
        return nil, status.Error(codes.InvalidArgument, "quantity must be positive")
    }
    
    // Create command
    cmd := commands.AddInventoryCommand{
        StoreID:   req.StoreId,
        ProductID: req.ProductId,
        Quantity:  int(req.Quantity),
    }
    
    // Execute command
    err := s.addInventoryHandler.Handle(ctx, cmd)
    if err != nil {
        // Convert domain errors to gRPC status
        return nil, toGRPCError(err)
    }
    
    // For demo, return the added quantity
    // In real app, would query for new total
    return &pb.AddInventoryResponse{
        NewQuantity: req.Quantity,
    }, nil
}

// UpdatePrice updates product price
func (s *StoreService) UpdatePrice(
    ctx context.Context,
    req *pb.UpdatePriceRequest,
) (*pb.UpdatePriceResponse, error) {
    // Validate request
    if req.NewPrice <= 0 {
        return nil, status.Error(codes.InvalidArgument, "price must be positive")
    }
    
    // Create command
    cmd := commands.UpdatePriceCommand{
        StoreID:   req.StoreId,
        ProductID: req.ProductId,
        NewPrice:  req.NewPrice,
        Currency:  req.Currency,
    }
    
    // Execute command
    err := s.updatePriceHandler.Handle(ctx, cmd)
    if err != nil {
        return nil, toGRPCError(err)
    }
    
    return &pb.UpdatePriceResponse{
        Success: true,
    }, nil
}

// GetProduct retrieves product details
func (s *StoreService) GetProduct(
    ctx context.Context,
    req *pb.GetProductRequest,
) (*pb.GetProductResponse, error) {
    // Create query
    query := queries.GetProductQuery{
        StoreID:   req.StoreId,
        ProductID: req.ProductId,
    }
    
    // Execute query
    productDTO, err := s.getProductHandler.Handle(ctx, query)
    if err != nil {
        return nil, toGRPCError(err)
    }
    
    // Convert to protobuf
    return &pb.GetProductResponse{
        Product: &pb.Product{
            Id:          productDTO.ID,
            Name:        productDTO.Name,
            Description: productDTO.Description,
            Price:       productDTO.Price,
            Currency:    productDTO.Currency,
            IsActive:    productDTO.IsActive,
        },
    }, nil
}

// Helper function to convert domain errors to gRPC status
func toGRPCError(err error) error {
    // Map domain errors to appropriate gRPC codes
    switch err {
    case store.ErrStoreNotFound:
        return status.Error(codes.NotFound, "store not found")
    case store.ErrProductNotFound:
        return status.Error(codes.NotFound, "product not found")
    case store.ErrInsufficientStock:
        return status.Error(codes.FailedPrecondition, "insufficient stock")
    default:
        return status.Error(codes.Internal, err.Error())
    }
}
