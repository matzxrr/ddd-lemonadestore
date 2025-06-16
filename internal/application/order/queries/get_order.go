package queries

import (
	"context"

	"github.com/matzxrr/ddd-lemonadestore/internal/application/dtos"
	"github.com/matzxrr/ddd-lemonadestore/internal/domain/order"
)

// GetOrderQuery represents request for order details
type GetOrderQuery struct {
    OrderID string
}

// GetOrderHandler handles order queries
type GetOrderHandler struct {
    orderRepo order.OrderRepository
}

func NewGetOrderHandler(orderRepo order.OrderRepository) *GetOrderHandler {
    return &GetOrderHandler{orderRepo: orderRepo}
}

func (h *GetOrderHandler) Handle(ctx context.Context, query GetOrderQuery) (*dtos.OrderDTO, error) {
    // Load order
    orderAgg, err := h.orderRepo.FindByID(order.OrderID(query.OrderID))
    if err != nil {
        return nil, err
    }
    
    // Convert to DTO
    items := make([]dtos.OrderItemDTO, len(orderAgg.Items()))
    for i, item := range orderAgg.Items() {
        items[i] = dtos.OrderItemDTO{
            ID:        item.ID(),
            ProductID: string(item.ProductID()),
            Name:      item.Name(),
            Quantity:  item.Quantity(),
            UnitPrice: float64(item.UnitPrice().Amount()) / 100,
            Total:     float64(item.Total().Amount()) / 100,
        }
    }
    
    return &dtos.OrderDTO{
        ID:          string(orderAgg.ID()),
        CustomerID:  string(orderAgg.CustomerID()),
        StoreID:     string(orderAgg.StoreID()),
        Status:      string(orderAgg.Status()),
        TotalAmount: float64(orderAgg.TotalAmount().Amount()) / 100,
        Currency:    orderAgg.TotalAmount().Currency(),
        Items:       items,
        PlacedAt:    orderAgg.PlacedAt(),
    }, nil
}
