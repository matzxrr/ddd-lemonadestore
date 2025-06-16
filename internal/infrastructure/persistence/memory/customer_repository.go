package memory

import (
	"errors"
	"sync"

	"github.com/matzxrr/ddd-lemonadestore/internal/domain/customer"
)

// InMemoryCustomerRepository is an in-memory implementation of CustomerRepository
type InMemoryCustomerRepository struct {
    mu         sync.RWMutex
    customers  map[customer.CustomerID]*customer.Customer
    emailIndex map[customer.Email]customer.CustomerID
}

// NewInMemoryCustomerRepository creates a new in-memory customer repository
func NewInMemoryCustomerRepository() *InMemoryCustomerRepository {
    return &InMemoryCustomerRepository{
        customers:  make(map[customer.CustomerID]*customer.Customer),
        emailIndex: make(map[customer.Email]customer.CustomerID),
    }
}

// Save persists a customer aggregate
func (r *InMemoryCustomerRepository) Save(customerAgg *customer.Customer) error {
    r.mu.Lock()
    defer r.mu.Unlock()
    
    // Save customer
    r.customers[customerAgg.ID()] = customerAgg
    
    // Update email index
    r.emailIndex[customerAgg.Email()] = customerAgg.ID()
    
    return nil
}

// FindByID retrieves a customer by ID
func (r *InMemoryCustomerRepository) FindByID(id customer.CustomerID) (*customer.Customer, error) {
    r.mu.RLock()
    defer r.mu.RUnlock()
    
    customerAgg, exists := r.customers[id]
    if !exists {
        return nil, errors.New("customer not found")
    }
    
    return customerAgg, nil
}

// FindByEmail retrieves a customer by email
// WHY: Email is a natural key for customers
func (r *InMemoryCustomerRepository) FindByEmail(email customer.Email) (*customer.Customer, error) {
    r.mu.RLock()
    defer r.mu.RUnlock()
    
    customerID, exists := r.emailIndex[email]
    if !exists {
        return nil, errors.New("customer not found")
    }
    
    return r.customers[customerID], nil
}

// FindByType retrieves customers by type
func (r *InMemoryCustomerRepository) FindByType(customerType customer.CustomerType) ([]*customer.Customer, error) {
    r.mu.RLock()
    defer r.mu.RUnlock()
    
    customers := make([]*customer.Customer, 0)
    for _, customerAgg := range r.customers {
        if customerAgg.Type() == customerType {
            customers = append(customers, customerAgg)
        }
    }
    
    return customers, nil
}
