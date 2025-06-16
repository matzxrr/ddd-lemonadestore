package store

import (
    "errors"
    "github.com/google/uuid"
)

// StoreID is a value object representing unique store identification
// WHY: Type safety - prevents passing wrong IDs to methods
type StoreID string

// NewStoreID creates a new unique store identifier
func NewStoreID() StoreID {
    return StoreID(uuid.New().String())
}

// ProductID uniquely identifies a product
// WHAT: Value object that ensures type safety for product references
type ProductID string

func NewProductID() ProductID {
    return ProductID(uuid.New().String())
}

// ProductName represents the name of a product with validation
// WHY: Business rule - products must have meaningful names
type ProductName string

func NewProductName(name string) (ProductName, error) {
    if len(name) < 3 || len(name) > 100 {
        return "", errors.New("product name must be between 3 and 100 characters")
    }
    return ProductName(name), nil
}

// Quantity represents product quantity with validation
// WHERE: Used in inventory management and order items
type Quantity int

func NewQuantity(q int) (Quantity, error) {
    if q < 0 {
        return 0, errors.New("quantity cannot be negative")
    }
    return Quantity(q), nil
}
