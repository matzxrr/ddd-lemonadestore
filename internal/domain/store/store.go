package store

import (
	"errors"

	"github.com/matzxrr/ddd-lemonadestore/internal/domain/shared"
)

// Store is the aggregate root for store management
// WHY: Store is the consistency boundary for products and inventory
// WHAT: Manages products and inventory as a cohesive unit
type Store struct {
    shared.AggregateRoot // Embed for event functionality
    id        StoreID
    name      string
    location  shared.Address
    products  map[ProductID]*Product
    inventory map[ProductID]Quantity
}

// NewStore creates a new store
// WHERE: Called during store initialization/setup
func NewStore(name string, location shared.Address) (*Store, error) {
    if name == "" {
        return nil, errors.New("store name is required")
    }
    
    store := &Store{
        id:        NewStoreID(),
        name:      name,
        location:  location,
        products:  make(map[ProductID]*Product),
        inventory: make(map[ProductID]Quantity),
    }
    
    // Raise domain event
    store.Raise(StoreCreatedEvent{
        BaseEvent: shared.NewBaseEvent(),
        StoreID:   string(store.id),
        StoreName: name,
        Location:  location,
    })
    
    return store, nil
}

// AddProduct adds a new product to the store
// WHY: Products can only be added through the Store aggregate
// WHAT: Ensures product uniqueness and raises events
func (s *Store) AddProduct(name string, description string, price shared.Money) (*Product, error) {
    product, err := NewProduct(name, description, price)
    if err != nil {
        return nil, err
    }
    
    // Check for duplicate product names
    for _, p := range s.products {
        if p.Name() == product.Name() && p.IsActive() {
            return nil, errors.New("product with this name already exists")
        }
    }
    
    s.products[product.ID()] = product
    s.inventory[product.ID()] = 0 // Initialize with zero inventory
    
    // Raise domain event
    s.Raise(ProductAddedEvent{
        BaseEvent:   shared.NewBaseEvent(),
        StoreID:     string(s.id),
        ProductID:   string(product.ID()),
        ProductName: string(product.Name()),
        Price:       price,
    })
    
    return product, nil
}

// AddInventory increases product quantity
// WHERE: Called when receiving new stock
func (s *Store) AddInventory(productID ProductID, quantity int) error {
    if _, exists := s.products[productID]; !exists {
        return errors.New("product not found")
    }
    
    if !s.products[productID].IsActive() {
        return errors.New("cannot add inventory to inactive product")
    }
    
    qty, err := NewQuantity(quantity)
    if err != nil {
        return err
    }
    
    s.inventory[productID] += qty
    
    // Raise domain event
    s.Raise(InventoryAddedEvent{
        BaseEvent:     shared.NewBaseEvent(),
        StoreID:       string(s.id),
        ProductID:     string(productID),
        QuantityAdded: int(qty),
        NewTotal:      int(s.inventory[productID]),
    })
    
    return nil
}

// ReserveInventory decreases product quantity for an order
// WHY: Ensures we don't oversell products
func (s *Store) ReserveInventory(productID ProductID, quantity int) error {
    currentQty, exists := s.inventory[productID]
    if !exists {
        return errors.New("product not found")
    }
    
    if int(currentQty) < quantity {
        return errors.New("insufficient inventory")
    }
    
    s.inventory[productID] = Quantity(int(currentQty) - quantity)
    
    // Raise domain event
    s.Raise(InventoryReservedEvent{
        BaseEvent:         shared.NewBaseEvent(),
        StoreID:          string(s.id),
        ProductID:        string(productID),
        QuantityReserved: quantity,
        RemainingQty:     int(s.inventory[productID]),
    })
    
    return nil
}

// GetProduct returns a product by ID
func (s *Store) GetProduct(productID ProductID) (*Product, error) {
    product, exists := s.products[productID]
    if !exists {
        return nil, errors.New("product not found")
    }
    return product, nil
}

// GetAvailableQuantity returns current inventory level
func (s *Store) GetAvailableQuantity(productID ProductID) (int, error) {
    qty, exists := s.inventory[productID]
    if !exists {
        return 0, errors.New("product not found")
    }
    return int(qty), nil
}

// Getters
func (s *Store) ID() StoreID              { return s.id }
func (s *Store) Name() string             { return s.name }
func (s *Store) Location() shared.Address { return s.location }
