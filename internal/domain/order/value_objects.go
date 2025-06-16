package order

import (
    "github.com/google/uuid"
)

// OrderID uniquely identifies an order
type OrderID string

func NewOrderID() OrderID {
    return OrderID(uuid.New().String())
}

// OrderStatus represents the state of an order
// WHY: Orders have a lifecycle with specific valid transitions
type OrderStatus string

const (
    OrderStatusPending   OrderStatus = "PENDING"
    OrderStatusConfirmed OrderStatus = "CONFIRMED"
    OrderStatusPreparing OrderStatus = "PREPARING"
    OrderStatusReady     OrderStatus = "READY"
    OrderStatusCompleted OrderStatus = "COMPLETED"
    OrderStatusCancelled OrderStatus = "CANCELLED"
)

// IsValidTransition checks if status transition is allowed
// WHAT: Implements business rules for order state machine
func (s OrderStatus) IsValidTransition(newStatus OrderStatus) bool {
    validTransitions := map[OrderStatus][]OrderStatus{
        OrderStatusPending:   {OrderStatusConfirmed, OrderStatusCancelled},
        OrderStatusConfirmed: {OrderStatusPreparing, OrderStatusCancelled},
        OrderStatusPreparing: {OrderStatusReady, OrderStatusCancelled},
        OrderStatusReady:     {OrderStatusCompleted},
        OrderStatusCompleted: {},
        OrderStatusCancelled: {},
    }
    
    allowed := validTransitions[s]
    for _, status := range allowed {
        if status == newStatus {
            return true
        }
    }
    return false
}
