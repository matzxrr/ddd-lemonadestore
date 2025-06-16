package store

import (
	"errors"

	"github.com/matzxrr/ddd-lemonadestore/internal/domain/shared"
)

// Product is an entity representing a lemonade product
// WHY: Products have identity and lifecycle within the store
// WHAT: Encapsulates product information and business rules
type Product struct {
    id          ProductID
    name        ProductName
    description string
    price       shared.Money
    isActive    bool
}

// NewProduct creates a new product with validation
// WHERE: Used by Store aggregate when adding new products
func NewProduct(name string, description string, price shared.Money) (*Product, error) {
    productName, err := NewProductName(name)
    if err != nil {
        return nil, err
    }
    
    if price.Amount() <= 0 {
        return nil, errors.New("product price must be greater than zero")
    }
    
    return &Product{
        id:          NewProductID(),
        name:        productName,
        description: description,
        price:       price,
        isActive:    true,
    }, nil
}

// UpdatePrice changes the product price
// WHY: Business rule - price changes must be tracked and validated
func (p *Product) UpdatePrice(newPrice shared.Money) error {
    if newPrice.Amount() <= 0 {
        return errors.New("price must be greater than zero")
    }
    if newPrice.Currency() != p.price.Currency() {
        return errors.New("cannot change currency of existing product")
    }
    p.price = newPrice
    return nil
}

// Deactivate marks product as unavailable
// WHAT: Soft delete - we don't remove products, just deactivate them
func (p *Product) Deactivate() {
    p.isActive = false
}

// Getters for encapsulation
func (p *Product) ID() ProductID         { return p.id }
func (p *Product) Name() ProductName     { return p.name }
func (p *Product) Description() string   { return p.description }
func (p *Product) Price() shared.Money   { return p.price }
func (p *Product) IsActive() bool        { return p.isActive }
