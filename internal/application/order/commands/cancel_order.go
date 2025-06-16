package commands

import (
	"context"

	"github.com/matzxrr/ddd-lemonadestore/internal/application/interfaces"
	"github.com/matzxrr/ddd-lemonadestore/internal/domain/order"
)

// CancelOrderCommand represents request to cancel an order
type CancelOrderCommand struct {
    OrderID string
    Reason  string
}

// CancelOrderHandler handles order cancellation
type CancelOrderHandler struct {
    uow            interfaces.UnitOfWork
    eventPublisher interfaces.EventPublisher
}

func NewCancelOrderHandler(
    uow interfaces.UnitOfWork,
    eventPublisher interfaces.EventPublisher,
) *CancelOrderHandler {
    return &CancelOrderHandler{
        uow:            uow,
        eventPublisher: eventPublisher,
    }
}

func (h *CancelOrderHandler) Handle(ctx context.Context, cmd CancelOrderCommand) error {
    // Start transaction
    err := h.uow.Begin(ctx)
    if err != nil {
        return err
    }
    defer func() {
        if err != nil {
            h.uow.Rollback()
        }
    }()
    
    // Load order
    orderAgg, err := h.uow.OrderRepository().FindByID(order.OrderID(cmd.OrderID))
    if err != nil {
        return err
    }
    
    // Cancel order
    err = orderAgg.Cancel(cmd.Reason)
    if err != nil {
        return err
    }
    
    // If order was confirmed, release inventory
    if orderAgg.Status() == order.OrderStatusCancelled {
        storeAgg, err := h.uow.StoreRepository().FindByID(orderAgg.StoreID())
        if err != nil {
            return err
        }
        
        // Return inventory for each item
        for _, item := range orderAgg.Items() {
            // In real app, would have ReleaseInventory method
            err = storeAgg.AddInventory(item.ProductID(), item.Quantity())
            if err != nil {
                // Log error but continue
            }
        }
        
        err = h.uow.StoreRepository().Save(storeAgg)
        if err != nil {
            return err
        }
    }
    
    // Save order
    err = h.uow.OrderRepository().Save(orderAgg)
    if err != nil {
        return err
    }
    
    // Commit
    err = h.uow.Commit()
    if err != nil {
        return err
    }
    
    // Publish events
    events := orderAgg.PullEvents()
    if len(events) > 0 {
        h.eventPublisher.Publish(ctx, events...)
    }
    
    return nil
}
