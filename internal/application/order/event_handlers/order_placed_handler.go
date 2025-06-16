package eventhandlers

import (
	"context"
	"log"

	"github.com/matzxrr/ddd-lemonadestore/internal/domain/customer"
	"github.com/matzxrr/ddd-lemonadestore/internal/domain/order"
	"github.com/matzxrr/ddd-lemonadestore/internal/domain/shared"
)

// OrderPlacedHandler handles OrderConfirmedEvent
// WHY: Decouples order processing from customer loyalty updates
type OrderPlacedHandler struct {
    customerRepo customer.CustomerRepository
}

func NewOrderPlacedHandler(customerRepo customer.CustomerRepository) *OrderPlacedHandler {
    return &OrderPlacedHandler{customerRepo: customerRepo}
}

// Handle processes the event
// WHERE: Registered with event bus to handle order.confirmed events
func (h *OrderPlacedHandler) Handle(ctx context.Context, event shared.DomainEvent) error {
    // Type assert to specific event
    orderConfirmed, ok := event.(order.OrderConfirmedEvent)
    if !ok {
        return nil // Not our event
    }
    
    // Load customer
    customerAgg, err := h.customerRepo.FindByID(customer.CustomerID(orderConfirmed.CustomerID))
    if err != nil {
        log.Printf("Failed to find customer %s: %v", orderConfirmed.CustomerID, err)
        return err
    }
    
    // Calculate loyalty points (1 point per dollar)
    points := int(orderConfirmed.TotalAmount.Amount() / 100)
    
    // Add points
    err = customerAgg.AddLoyaltyPoints(points)
    if err != nil {
        log.Printf("Failed to add loyalty points: %v", err)
        return err
    }
    
    // Save customer
    err = h.customerRepo.Save(customerAgg)
    if err != nil {
        log.Printf("Failed to save customer: %v", err)
        return err
    }
    
    log.Printf("Added %d loyalty points to customer %s for order %s",
        points, orderConfirmed.CustomerID, orderConfirmed.OrderID)
    
    return nil
}
