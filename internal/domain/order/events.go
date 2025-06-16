package order

import (
	"time"

	"github.com/matzxrr/ddd-lemonadestore/internal/domain/shared"
)

// OrderItemSnapshot is an immutable representation of an order item for events
type OrderItemSnapshot struct {
    ProductID string       `json:"product_id"`
    Name      string       `json:"name"`
    Quantity  int          `json:"quantity"`
    UnitPrice shared.Money `json:"unit_price"`
    Total     shared.Money `json:"total"`
}

// OrderCreatedEvent is raised when a new order is created
type OrderCreatedEvent struct {
    shared.BaseEvent
    OrderID    string `json:"order_id"`
    CustomerID string `json:"customer_id"`
    StoreID    string `json:"store_id"`
}

func (e OrderCreatedEvent) EventName() string     { return "order.created" }
func (e OrderCreatedEvent) AggregateID() string   { return e.OrderID }
func (e OrderCreatedEvent) AggregateType() string { return "order" }

// OrderConfirmedEvent contains full order details for downstream processing
type OrderConfirmedEvent struct {
    shared.BaseEvent
    OrderID     string              `json:"order_id"`
    CustomerID  string              `json:"customer_id"`
    StoreID     string              `json:"store_id"`
    TotalAmount shared.Money        `json:"total_amount"`
    Items       []OrderItemSnapshot `json:"items"`
}

func (e OrderConfirmedEvent) EventName() string     { return "order.confirmed" }
func (e OrderConfirmedEvent) AggregateID() string   { return e.OrderID }
func (e OrderConfirmedEvent) AggregateType() string { return "order" }

// OrderCancelledEvent is raised when order is cancelled
type OrderCancelledEvent struct {
    shared.BaseEvent
    OrderID    string `json:"order_id"`
    CustomerID string `json:"customer_id"`
    Reason     string `json:"reason"`
}

func (e OrderCancelledEvent) EventName() string     { return "order.cancelled" }
func (e OrderCancelledEvent) AggregateID() string   { return e.OrderID }
func (e OrderCancelledEvent) AggregateType() string { return "order" }

// OrderPreparationStartedEvent indicates order preparation has begun
type OrderPreparationStartedEvent struct {
    shared.BaseEvent
    OrderID string `json:"order_id"`
}

func (e OrderPreparationStartedEvent) EventName() string     { return "order.preparation_started" }
func (e OrderPreparationStartedEvent) AggregateID() string   { return e.OrderID }
func (e OrderPreparationStartedEvent) AggregateType() string { return "order" }

// OrderReadyEvent notifies customer their order is ready
type OrderReadyEvent struct {
    shared.BaseEvent
    OrderID    string `json:"order_id"`
    CustomerID string `json:"customer_id"`
}

func (e OrderReadyEvent) EventName() string     { return "order.ready" }
func (e OrderReadyEvent) AggregateID() string   { return e.OrderID }
func (e OrderReadyEvent) AggregateType() string { return "order" }

// OrderCompletedEvent marks successful order completion
type OrderCompletedEvent struct {
    shared.BaseEvent
    OrderID     string    `json:"order_id"`
    CustomerID  string    `json:"customer_id"`
    CompletedAt time.Time `json:"completed_at"`
}

func (e OrderCompletedEvent) EventName() string     { return "order.completed" }
func (e OrderCompletedEvent) AggregateID() string   { return e.OrderID }
func (e OrderCompletedEvent) AggregateType() string { return "order" }
