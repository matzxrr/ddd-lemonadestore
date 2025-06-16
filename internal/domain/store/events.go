package store

import "github.com/matzxrr/ddd-lemonadestore/internal/domain/shared"

// StoreCreatedEvent is raised when a new store is created
// WHY: Other parts of the system need to know when stores are created
type StoreCreatedEvent struct {
	shared.BaseEvent
	StoreID   string         `json:"store_id"`
	StoreName string         `json:"store_name"`
	Location  shared.Address `json:"location"`
}

func (e StoreCreatedEvent) EventName() string     { return "store.created" }
func (e StoreCreatedEvent) AggregateID() string   { return e.StoreID }
func (e StoreCreatedEvent) AggregateType() string { return "store" }

// ProductAddedEvent is raised when a product is added to store
// WHERE: Used by read models to update product catalogs
type ProductAddedEvent struct {
	shared.BaseEvent
	StoreID     string       `json:"store_id"`
	ProductID   string       `json:"product_id"`
	ProductName string       `json:"product_name"`
	Price       shared.Money `json:"price"`
}

func (e ProductAddedEvent) EventName() string     { return "product.added" }
func (e ProductAddedEvent) AggregateID() string   { return e.StoreID }
func (e ProductAddedEvent) AggregateType() string { return "store" }

// InventoryAddedEvent tracks inventory increases
type InventoryAddedEvent struct {
	shared.BaseEvent
	StoreID       string `json:"store_id"`
	ProductID     string `json:"product_id"`
	QuantityAdded int    `json:"quantity_added"`
	NewTotal      int    `json:"new_total"`
}

func (e InventoryAddedEvent) EventName() string     { return "inventory.added" }
func (e InventoryAddedEvent) AggregateID() string   { return e.StoreID }
func (e InventoryAddedEvent) AggregateType() string { return "store" }

// InventoryReservedEvent tracks inventory reservations
type InventoryReservedEvent struct {
	shared.BaseEvent
	StoreID          string `json:"store_id"`
	ProductID        string `json:"product_id"`
	QuantityReserved int    `json:"quantity_reserved"`
	RemainingQty     int    `json:"remaining_qty"`
}

func (e InventoryReservedEvent) EventName() string     { return "inventory.reserved" }
func (e InventoryReservedEvent) AggregateID() string   { return e.StoreID }
func (e InventoryReservedEvent) AggregateType() string { return "store" }
