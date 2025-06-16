package commands

import (
	"context"

	"github.com/matzxrr/ddd-lemonadestore/internal/application/interfaces"
	"github.com/matzxrr/ddd-lemonadestore/internal/domain/shared"
	"github.com/matzxrr/ddd-lemonadestore/internal/domain/store"
)

// UpdatePriceCommand represents request to update product price
type UpdatePriceCommand struct {
    StoreID   string
    ProductID string
    NewPrice  float64
    Currency  string
}

// UpdatePriceHandler handles price updates
type UpdatePriceHandler struct {
    storeRepo     store.StoreRepository
    eventPublisher interfaces.EventPublisher
}

func NewUpdatePriceHandler(
    storeRepo store.StoreRepository,
    eventPublisher interfaces.EventPublisher,
) *UpdatePriceHandler {
    return &UpdatePriceHandler{
        storeRepo:     storeRepo,
        eventPublisher: eventPublisher,
    }
}

func (h *UpdatePriceHandler) Handle(ctx context.Context, cmd UpdatePriceCommand) error {
    // Convert price to domain value object
    priceCents := int64(cmd.NewPrice * 100)
    newPrice, err := shared.NewMoney(priceCents, cmd.Currency)
    if err != nil {
        return err
    }
    
    // Load store aggregate
    storeAgg, err := h.storeRepo.FindByID(store.StoreID(cmd.StoreID))
    if err != nil {
        return err
    }
    
    // Get product and update price
    product, err := storeAgg.GetProduct(store.ProductID(cmd.ProductID))
    if err != nil {
        return err
    }
    
    err = product.UpdatePrice(newPrice)
    if err != nil {
        return err
    }
    
    // Save changes
    err = h.storeRepo.Save(storeAgg)
    if err != nil {
        return err
    }
    
    // Publish events
    events := storeAgg.PullEvents()
    if len(events) > 0 {
        h.eventPublisher.Publish(ctx, events...)
    }
    
    return nil
}
