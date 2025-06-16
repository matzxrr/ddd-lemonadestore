package queries

import (
	"context"

	"github.com/matzxrr/ddd-lemonadestore/internal/application/dtos"
	"github.com/matzxrr/ddd-lemonadestore/internal/domain/customer"
	"github.com/matzxrr/ddd-lemonadestore/internal/domain/order"
)

// ListOrdersQuery represents request for customer orders
type ListOrdersQuery struct {
    CustomerID string
}

// ListOrdersHandler handles order listing
type ListOrdersHandler struct {
    orderRepo order.OrderRepository
}

func NewListOrdersHandler(orderRepo order.OrderRepository) *ListOrdersHandler {
    return &ListOrdersHandler{orderRepo: orderRepo}
}

func (h *ListOrdersHandler) Handle(ctx context.Context, query ListOrdersQuery) ([]*dtos.OrderDTO, error) {
    // Find orders
    orders, err := h.orderRepo.FindByCustomer(customer.CustomerID(query.CustomerID))
    if err != nil {
        return nil, err
    }
    
    // Convert to DTOs
    result := make([]*dtos.OrderDTO, len(orders))
    for i, orderAgg := range orders {
        items := make([]dtos.OrderItemDTO, len(orderAgg.Items()))
        for j, item := range orderAgg.Items() {
            items[j] = dtos.OrderItemDTO{
                ID:        item.ID(),
                ProductID: string(item.ProductID()),
                Name:      item.Name(),
                Quantity:  item.Quantity(),
                UnitPrice: float64(item.UnitPrice().Amount()) / 100,
                Total:     float64(item.Total().Amount()) / 100,
            }
        }
        
        result[i] = &dtos.OrderDTO{
            ID:          string(orderAgg.ID()),
            CustomerID:  string(orderAgg.CustomerID()),
            StoreID:     string(orderAgg.StoreID()),
            Status:      string(orderAgg.Status()),
            TotalAmount: float64(orderAgg.TotalAmount().Amount()) / 100,
            Currency:    orderAgg.TotalAmount().Currency(),
            Items:       items,
            PlacedAt:    orderAgg.PlacedAt(),
        }
    }
    
    return result, nil
}
