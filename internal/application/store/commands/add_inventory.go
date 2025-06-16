package commands

import (
	"context"
	"errors"

	"github.com/matzxrr/ddd-lemonadestore/internal/application/interfaces"
	"github.com/matzxrr/ddd-lemonadestore/internal/domain/store"
)

// AddInventoryCommand represents request to add inventory
// WHY: Separates external request format from domain logic
type AddInventoryCommand struct {
    StoreID   string
    ProductID string
    Quantity  int
}

// AddInventoryHandler handles inventory addition
// WHERE: Called from presentation layer when new stock arrives
type AddInventoryHandler struct {
    storeRepo     store.StoreRepository
    eventPublisher interfaces.EventPublisher
}

func NewAddInventoryHandler(
    storeRepo store.StoreRepository,
    eventPublisher interfaces.EventPublisher,
) *AddInventoryHandler {
    return &AddInventoryHandler{
        storeRepo:     storeRepo,
        eventPublisher: eventPublisher,
    }
}

// Handle executes the add inventory use case
// WHAT: Orchestrates loading aggregate, executing domain logic, and persisting
func (h *AddInventoryHandler) Handle(ctx context.Context, cmd AddInventoryCommand) error {
    // 1. Validate command
    if cmd.Quantity <= 0 {
        return errors.New("quantity must be positive")
    }
    
    // 2. Load aggregate
    storeAgg, err := h.storeRepo.FindByID(store.StoreID(cmd.StoreID))
    if err != nil {
        return err
    }
    
    // 3. Execute domain logic
    err = storeAgg.AddInventory(store.ProductID(cmd.ProductID), cmd.Quantity)
    if err != nil {
        return err
    }
    
    // 4. Persist changes
    err = h.storeRepo.Save(storeAgg)
    if err != nil {
        return err
    }
    
    // 5. Publish domain events
    events := storeAgg.PullEvents()
    if len(events) > 0 {
        err = h.eventPublisher.Publish(ctx, events...)
        if err != nil {
            // Log error but don't fail the operation
            // Events can be published asynchronously
        }
    }
    
    return nil
}
