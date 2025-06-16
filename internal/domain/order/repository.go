package order

import "github.com/matzxrr/ddd-lemonadestore/internal/domain/customer"

// OrderRepository defines persistence operations for Order aggregate
type OrderRepository interface {
    Save(order *Order) error
    FindByID(id OrderID) (*Order, error)
    FindByCustomer(customerID customer.CustomerID) ([]*Order, error)
    FindByStatus(status OrderStatus) ([]*Order, error)
}
