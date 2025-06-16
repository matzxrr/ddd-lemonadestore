package memory

import (
	"errors"
	"sync"

	"github.com/matzxrr/ddd-lemonadestore/internal/domain/customer"
	"github.com/matzxrr/ddd-lemonadestore/internal/domain/order"
)

// InMemoryOrderRepository is an in-memory implementation of OrderRepository
type InMemoryOrderRepository struct {
    mu     sync.RWMutex
    orders map[order.OrderID]*order.Order
    // Secondary indexes for queries
    customerIndex map[customer.CustomerID][]order.OrderID
}

// NewInMemoryOrderRepository creates a new in-memory order repository
func NewInMemoryOrderRepository() *InMemoryOrderRepository {
    return &InMemoryOrderRepository{
        orders:        make(map[order.OrderID]*order.Order),
        customerIndex: make(map[customer.CustomerID][]order.OrderID),
    }
}

// Save persists an order aggregate
func (r *InMemoryOrderRepository) Save(orderAgg *order.Order) error {
    r.mu.Lock()
    defer r.mu.Unlock()
    
    // Save order
    r.orders[orderAgg.ID()] = orderAgg
    
    // Update customer index
    customerID := orderAgg.CustomerID()
    if _, exists := r.customerIndex[customerID]; !exists {
        r.customerIndex[customerID] = []order.OrderID{}
    }
    
    // Check if order already in index
    found := false
    for _, id := range r.customerIndex[customerID] {
        if id == orderAgg.ID() {
            found = true
            break
        }
    }
    
    if !found {
        r.customerIndex[customerID] = append(r.customerIndex[customerID], orderAgg.ID())
    }
    
    return nil
}

// FindByID retrieves an order by ID
func (r *InMemoryOrderRepository) FindByID(id order.OrderID) (*order.Order, error) {
    r.mu.RLock()
    defer r.mu.RUnlock()
    
    orderAgg, exists := r.orders[id]
    if !exists {
        return nil, errors.New("order not found")
    }
    
    return orderAgg, nil
}

// FindByCustomer retrieves orders for a customer
// WHY: Supports customer order history queries
func (r *InMemoryOrderRepository) FindByCustomer(customerID customer.CustomerID) ([]*order.Order, error) {
    r.mu.RLock()
    defer r.mu.RUnlock()
    
    orderIDs, exists := r.customerIndex[customerID]
    if !exists {
        return []*order.Order{}, nil
    }
    
    orders := make([]*order.Order, 0, len(orderIDs))
    for _, orderID := range orderIDs {
        if orderAgg, exists := r.orders[orderID]; exists {
            orders = append(orders, orderAgg)
        }
    }
    
    return orders, nil
}

// FindByStatus retrieves orders by status
func (r *InMemoryOrderRepository) FindByStatus(status order.OrderStatus) ([]*order.Order, error) {
    r.mu.RLock()
    defer r.mu.RUnlock()
    
    orders := make([]*order.Order, 0)
    for _, orderAgg := range r.orders {
        if orderAgg.Status() == status {
            orders = append(orders, orderAgg)
        }
    }
    
    return orders, nil
}
