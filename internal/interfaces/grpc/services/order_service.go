package services

import (
	"context"

	"github.com/matzxrr/ddd-lemonadestore/internal/application/order/commands"
	"github.com/matzxrr/ddd-lemonadestore/internal/application/order/queries"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"
	pb "github.com/matzxrr/ddd-lemonadestore/internal/interfaces/grpc/pb/order/v1"
)

// OrderService implements the gRPC OrderService
type OrderService struct {
    pb.UnimplementedOrderServiceServer
    
    // Command handlers
    createOrderHandler *commands.CreateOrderHandler
    cancelOrderHandler *commands.CancelOrderHandler
    
    // Query handlers
    getOrderHandler    *queries.GetOrderHandler
    listOrdersHandler  *queries.ListOrdersHandler
}

// NewOrderService creates a new order service
func NewOrderService(
    createOrder *commands.CreateOrderHandler,
    cancelOrder *commands.CancelOrderHandler,
    getOrder *queries.GetOrderHandler,
    listOrders *queries.ListOrdersHandler,
) *OrderService {
    return &OrderService{
        createOrderHandler: createOrder,
        cancelOrderHandler: cancelOrder,
        getOrderHandler:    getOrder,
        listOrdersHandler:  listOrders,
    }
}

// CreateOrder creates a new order
// WHY: Main entry point for customer purchases
func (s *OrderService) CreateOrder(
    ctx context.Context,
    req *pb.CreateOrderRequest,
) (*pb.CreateOrderResponse, error) {
    // Validate request
    if req.CustomerId == "" || req.StoreId == "" {
        return nil, status.Error(codes.InvalidArgument, "customer_id and store_id are required")
    }
    
    if len(req.Items) == 0 {
        return nil, status.Error(codes.InvalidArgument, "order must have at least one item")
    }
    
    // Convert items
    items := make([]commands.OrderItemRequest, len(req.Items))
    for i, item := range req.Items {
        if item.Quantity <= 0 {
            return nil, status.Error(codes.InvalidArgument, "item quantity must be positive")
        }
        items[i] = commands.OrderItemRequest{
            ProductID: item.ProductId,
            Quantity:  int(item.Quantity),
        }
    }
    
    // Create command
    cmd := commands.CreateOrderCommand{
        CustomerID: req.CustomerId,
        StoreID:    req.StoreId,
        Items:      items,
    }
    
    // Execute command
    orderDTO, err := s.createOrderHandler.Handle(ctx, cmd)
    if err != nil {
        return nil, toGRPCError(err)
    }
    
    return &pb.CreateOrderResponse{
        OrderId:     orderDTO.ID,
        TotalAmount: orderDTO.TotalAmount,
        Currency:    orderDTO.Currency,
    }, nil
}

// CancelOrder cancels an existing order
func (s *OrderService) CancelOrder(
    ctx context.Context,
    req *pb.CancelOrderRequest,
) (*pb.CancelOrderResponse, error) {
    // Create command
    cmd := commands.CancelOrderCommand{
        OrderID: req.OrderId,
        Reason:  req.Reason,
    }
    
    // Execute command
    err := s.cancelOrderHandler.Handle(ctx, cmd)
    if err != nil {
        return nil, toGRPCError(err)
    }
    
    return &pb.CancelOrderResponse{
        Success: true,
    }, nil
}

// GetOrder retrieves order details
func (s *OrderService) GetOrder(
    ctx context.Context,
    req *pb.GetOrderRequest,
) (*pb.GetOrderResponse, error) {
    // Create query
    query := queries.GetOrderQuery{
        OrderID: req.OrderId,
    }
    
    // Execute query
    orderDTO, err := s.getOrderHandler.Handle(ctx, query)
    if err != nil {
        return nil, toGRPCError(err)
    }
    
    // Convert to protobuf
    items := make([]*pb.OrderItemDetail, len(orderDTO.Items))
    for i, item := range orderDTO.Items {
        items[i] = &pb.OrderItemDetail{
            Id:        item.ID,
            ProductId: item.ProductID,
            Name:      item.Name,
            Quantity:  int32(item.Quantity),
            UnitPrice: item.UnitPrice,
            Total:     item.Total,
        }
    }
    
    return &pb.GetOrderResponse{
        Order: &pb.Order{
            Id:          orderDTO.ID,
            CustomerId:  orderDTO.CustomerID,
            StoreId:     orderDTO.StoreID,
            Status:      orderDTO.Status,
            TotalAmount: orderDTO.TotalAmount,
            Currency:    orderDTO.Currency,
            Items:       items,
            PlacedAt:    timestamppb.New(orderDTO.PlacedAt),
        },
    }, nil
}

// ListCustomerOrders lists orders for a customer
func (s *OrderService) ListCustomerOrders(
    ctx context.Context,
    req *pb.ListCustomerOrdersRequest,
) (*pb.ListCustomerOrdersResponse, error) {
    // Create query
    query := queries.ListOrdersQuery{
        CustomerID: req.CustomerId,
    }
    
    // Execute query
    orders, err := s.listOrdersHandler.Handle(ctx, query)
    if err != nil {
        return nil, toGRPCError(err)
    }
    
    // Convert to protobuf
    pbOrders := make([]*pb.Order, len(orders))
    for i, order := range orders {
        items := make([]*pb.OrderItemDetail, len(order.Items))
        for j, item := range order.Items {
            items[j] = &pb.OrderItemDetail{
                Id:        item.ID,
                ProductId: item.ProductID,
                Name:      item.Name,
                Quantity:  int32(item.Quantity),
                UnitPrice: item.UnitPrice,
                Total:     item.Total,
            }
        }
        
        pbOrders[i] = &pb.Order{
            Id:          order.ID,
            CustomerId:  order.CustomerID,
            StoreId:     order.StoreID,
            Status:      order.Status,
            TotalAmount: order.TotalAmount,
            Currency:    order.Currency,
            Items:       items,
            PlacedAt:    timestamppb.New(order.PlacedAt),
        }
    }
    
    return &pb.ListCustomerOrdersResponse{
        Orders: pbOrders,
    }, nil
}
