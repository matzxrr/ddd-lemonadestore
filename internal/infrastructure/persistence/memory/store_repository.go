package memory

import (
	"sync"

	"github.com/matzxrr/ddd-lemonadestore/internal/domain/store"
)

// InMemoryStoreRepository is an in-memory implementation of StoreRepository
// WHY: For testing and demo purposes without database dependency
// WHERE: Injected into application services in main.go
type InMemoryStoreRepository struct {
    mu     sync.RWMutex
    stores map[store.StoreID]*store.Store
}

// NewInMemoryStoreRepository creates a new in-memory store repository
func NewInMemoryStoreRepository() *InMemoryStoreRepository {
    return &InMemoryStoreRepository{
        stores: make(map[store.StoreID]*store.Store),
    }
}

// Save persists a store aggregate
// WHAT: Stores the entire aggregate in memory
func (r *InMemoryStoreRepository) Save(storeAgg *store.Store) error {
    r.mu.Lock()
    defer r.mu.Unlock()
    
    r.stores[storeAgg.ID()] = storeAgg
    return nil
}

// FindByID retrieves a store by ID
func (r *InMemoryStoreRepository) FindByID(id store.StoreID) (*store.Store, error) {
    r.mu.RLock()
    defer r.mu.RUnlock()
    
    storeAgg, exists := r.stores[id]
    if !exists {
        return nil, store.ErrStoreNotFound
    }
    
    return storeAgg, nil
}

// FindAll returns all stores
func (r *InMemoryStoreRepository) FindAll() ([]*store.Store, error) {
    r.mu.RLock()
    defer r.mu.RUnlock()
    
    stores := make([]*store.Store, 0, len(r.stores))
    for _, storeAgg := range r.stores {
        stores = append(stores, storeAgg)
    }
    
    return stores, nil
}
