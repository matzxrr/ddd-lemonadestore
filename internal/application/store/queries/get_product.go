package queries

import (
	"context"

	"github.com/matzxrr/ddd-lemonadestore/internal/application/dtos"
	"github.com/matzxrr/ddd-lemonadestore/internal/domain/store"
)

// GetProductQuery represents request for product details
type GetProductQuery struct {
    StoreID   string
    ProductID string
}

// GetProductHandler handles product queries
type GetProductHandler struct {
    storeRepo store.StoreRepository
}

func NewGetProductHandler(storeRepo store.StoreRepository) *GetProductHandler {
    return &GetProductHandler{storeRepo: storeRepo}
}

func (h *GetProductHandler) Handle(ctx context.Context, query GetProductQuery) (*dtos.ProductDTO, error) {
    // Load store
    storeAgg, err := h.storeRepo.FindByID(store.StoreID(query.StoreID))
    if err != nil {
        return nil, err
    }
    
    // Get product
    product, err := storeAgg.GetProduct(store.ProductID(query.ProductID))
    if err != nil {
        return nil, err
    }
    
    // Get quantity
    qty, _ := storeAgg.GetAvailableQuantity(product.ID())
    
    // Convert to DTO
    return &dtos.ProductDTO{
        ID:          string(product.ID()),
        Name:        string(product.Name()),
        Description: product.Description(),
        Price:       float64(product.Price().Amount()) / 100,
        Currency:    product.Price().Currency(),
        IsActive:    product.IsActive(),
        Quantity:    qty,
    }, nil
}
