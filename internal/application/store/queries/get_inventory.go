package queries

import (
	"context"

	"github.com/matzxrr/ddd-lemonadestore/internal/application/dtos"
	"github.com/matzxrr/ddd-lemonadestore/internal/domain/store"
)

// GetInventoryQuery represents request for inventory information
type GetInventoryQuery struct {
    StoreID string
}

// GetInventoryHandler handles inventory queries
// WHY: Read operations don't need full domain logic
type GetInventoryHandler struct {
    storeRepo store.StoreRepository
}

func NewGetInventoryHandler(storeRepo store.StoreRepository) *GetInventoryHandler {
    return &GetInventoryHandler{storeRepo: storeRepo}
}

// Handle returns inventory for all products in store
func (h *GetInventoryHandler) Handle(ctx context.Context, query GetInventoryQuery) ([]dtos.ProductDTO, error) {
    // Load store
    storeAgg, err := h.storeRepo.FindByID(store.StoreID(query.StoreID))
    if err != nil {
        return nil, err
    }
    
    // Convert to DTOs
    var products []dtos.ProductDTO
    
    // In real implementation, would iterate through all products
    // For demo, returning empty slice
    return products, nil
}
