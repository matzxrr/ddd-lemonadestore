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
    

    // Iterate through all products and get their inventory
    for productID, product := range storeAgg.Products() {
        qty, _ := storeAgg.GetAvailableQuantity(productID)

        products = append(products, dtos.ProductDTO{
            ID: string(product.ID()),
            Name: string(product.Name()),
            Description: product.Description(),
            Price: float64(product.Price().Amount()) / 100,
            Currency: product.Price().Currency(),
            IsActive: product.IsActive(),
            Quantity: qty,
        })
    }
    
    return products, nil
}
