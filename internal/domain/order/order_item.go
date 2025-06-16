package order

import (
	"errors"

	"github.com/google/uuid"
	"github.com/matzxrr/ddd-lemonadestore/internal/domain/shared"
	"github.com/matzxrr/ddd-lemonadestore/internal/domain/store"
)

// OrderItem is an entity within the Order aggregate
// WHY: Order items have identity within an order (can be modified/removed)
// WHAT: Represents a line item in an order
type OrderItem struct {
    id        string
    productID store.ProductID
    name      string
    quantity  int
    unitPrice shared.Money
}

// NewOrderItem creates a new order item
// WHERE: Created when adding items to an order
func NewOrderItem(productID store.ProductID, name string, quantity int, unitPrice shared.Money) (*OrderItem, error) {
    if quantity <= 0 {
        return nil, errors.New("quantity must be positive")
    }
    
    return &OrderItem{
        id:        uuid.New().String(),
        productID: productID,
        name:      name,
        quantity:  quantity,
        unitPrice: unitPrice,
    }, nil
}

// Total calculates the total price for this item
func (i *OrderItem) Total() shared.Money {
    return i.unitPrice.Multiply(i.quantity)
}

// UpdateQuantity changes the item quantity
// WHY: Business rule - quantity can be adjusted before order confirmation
func (i *OrderItem) UpdateQuantity(newQuantity int) error {
    if newQuantity <= 0 {
        return errors.New("quantity must be positive")
    }
    i.quantity = newQuantity
    return nil
}

// Getters
func (i *OrderItem) ID() string                 { return i.id }
func (i *OrderItem) ProductID() store.ProductID { return i.productID }
func (i *OrderItem) Name() string               { return i.name }
func (i *OrderItem) Quantity() int              { return i.quantity }
func (i *OrderItem) UnitPrice() shared.Money    { return i.unitPrice }
