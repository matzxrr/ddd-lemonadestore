package order

import (
	"errors"
	"time"

	"github.com/matzxrr/ddd-lemonadestore/internal/domain/customer"
	"github.com/matzxrr/ddd-lemonadestore/internal/domain/shared"
	"github.com/matzxrr/ddd-lemonadestore/internal/domain/store"
)

// Order is the aggregate root for order management
// WHY: Order maintains consistency for all order-related data
// WHAT: Represents a customer's purchase transaction
type Order struct {
    shared.AggregateRoot
    id          OrderID
    customerID  customer.CustomerID
    storeID     store.StoreID
    items       []*OrderItem
    status      OrderStatus
    totalAmount shared.Money
    placedAt    time.Time
    notes       string
}

// NewOrder creates a new order
// WHERE: Called when customer initiates a purchase
func NewOrder(customerID customer.CustomerID, storeID store.StoreID) *Order {
    order := &Order{
        id:         NewOrderID(),
        customerID: customerID,
        storeID:    storeID,
        items:      make([]*OrderItem, 0),
        status:     OrderStatusPending,
        placedAt:   time.Now(),
    }
    
    // Raise domain event
    order.Raise(OrderCreatedEvent{
        BaseEvent:  shared.NewBaseEvent(),
        OrderID:    string(order.id),
        CustomerID: string(customerID),
        StoreID:    string(storeID),
    })
    
    return order
}

// AddItem adds a product to the order
// WHY: Orders can only be modified through aggregate methods
func (o *Order) AddItem(productID store.ProductID, name string, quantity int, unitPrice shared.Money) error {
    if o.status != OrderStatusPending {
        return errors.New("can only add items to pending orders")
    }
    
    // Check if item already exists
    for _, item := range o.items {
        if item.ProductID() == productID {
            // Update quantity instead of adding duplicate
            return item.UpdateQuantity(item.Quantity() + quantity)
        }
    }
    
    item, err := NewOrderItem(productID, name, quantity, unitPrice)
    if err != nil {
        return err
    }
    
    o.items = append(o.items, item)
    o.recalculateTotal()
    
    return nil
}

// RemoveItem removes a product from the order
func (o *Order) RemoveItem(itemID string) error {
    if o.status != OrderStatusPending {
        return errors.New("can only remove items from pending orders")
    }
    
    for i, item := range o.items {
        if item.ID() == itemID {
            o.items = append(o.items[:i], o.items[i+1:]...)
            o.recalculateTotal()
            return nil
        }
    }
    
    return errors.New("item not found in order")
}

// Confirm moves order to confirmed state
// WHERE: Called after payment is processed
func (o *Order) Confirm() error {
    if !o.status.IsValidTransition(OrderStatusConfirmed) {
        return errors.New("cannot confirm order in current status")
    }
    
    if len(o.items) == 0 {
        return errors.New("cannot confirm empty order")
    }
    
    o.status = OrderStatusConfirmed
    
    // Raise domain event with order snapshot
    o.Raise(OrderConfirmedEvent{
        BaseEvent:   shared.NewBaseEvent(),
        OrderID:     string(o.id),
        CustomerID:  string(o.customerID),
        StoreID:     string(o.storeID),
        TotalAmount: o.totalAmount,
        Items:       o.createItemSnapshots(),
    })
    
    return nil
}

// Cancel cancels the order
// WHY: Orders can be cancelled before completion
func (o *Order) Cancel(reason string) error {
    if !o.status.IsValidTransition(OrderStatusCancelled) {
        return errors.New("cannot cancel order in current status")
    }
    
    o.status = OrderStatusCancelled
    
    // Raise domain event
    o.Raise(OrderCancelledEvent{
        BaseEvent:  shared.NewBaseEvent(),
        OrderID:    string(o.id),
        CustomerID: string(o.customerID),
        Reason:     reason,
    })
    
    return nil
}

// StartPreparing moves order to preparing state
func (o *Order) StartPreparing() error {
    if !o.status.IsValidTransition(OrderStatusPreparing) {
        return errors.New("cannot start preparing order in current status")
    }
    
    o.status = OrderStatusPreparing
    
    o.Raise(OrderPreparationStartedEvent{
        BaseEvent: shared.NewBaseEvent(),
        OrderID:   string(o.id),
    })
    
    return nil
}

// MarkReady indicates order is ready for pickup
func (o *Order) MarkReady() error {
    if !o.status.IsValidTransition(OrderStatusReady) {
        return errors.New("cannot mark order ready in current status")
    }
    
    o.status = OrderStatusReady
    
    o.Raise(OrderReadyEvent{
        BaseEvent:  shared.NewBaseEvent(),
        OrderID:    string(o.id),
        CustomerID: string(o.customerID),
    })
    
    return nil
}

// Complete marks order as completed
func (o *Order) Complete() error {
    if !o.status.IsValidTransition(OrderStatusCompleted) {
        return errors.New("cannot complete order in current status")
    }
    
    o.status = OrderStatusCompleted
    
    o.Raise(OrderCompletedEvent{
        BaseEvent:   shared.NewBaseEvent(),
        OrderID:     string(o.id),
        CustomerID:  string(o.customerID),
        CompletedAt: time.Now(),
    })
    
    return nil
}

// recalculateTotal updates the order total
// WHAT: Private method that maintains total consistency
func (o *Order) recalculateTotal() {
    if len(o.items) == 0 {
        o.totalAmount = shared.Money{} // Zero value
        return
    }
    
    total := o.items[0].Total()
    for i := 1; i < len(o.items); i++ {
        itemTotal := o.items[i].Total()
        total, _ = total.Add(itemTotal)
    }
    o.totalAmount = total
}

// createItemSnapshots creates immutable snapshots for events
func (o *Order) createItemSnapshots() []OrderItemSnapshot {
    snapshots := make([]OrderItemSnapshot, len(o.items))
    for i, item := range o.items {
        snapshots[i] = OrderItemSnapshot{
            ProductID: string(item.ProductID()),
            Name:      item.Name(),
            Quantity:  item.Quantity(),
            UnitPrice: item.UnitPrice(),
            Total:     item.Total(),
        }
    }
    return snapshots
}

// Getters
func (o *Order) ID() OrderID                     { return o.id }
func (o *Order) CustomerID() customer.CustomerID { return o.customerID }
func (o *Order) StoreID() store.StoreID          { return o.storeID }
func (o *Order) Items() []*OrderItem             { return o.items }
func (o *Order) Status() OrderStatus             { return o.status }
func (o *Order) TotalAmount() shared.Money       { return o.totalAmount }
func (o *Order) PlacedAt() time.Time             { return o.placedAt }
