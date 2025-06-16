package commands

import (
	"context"
	"errors"

	"github.com/matzxrr/ddd-lemonadestore/internal/application/dtos"
	"github.com/matzxrr/ddd-lemonadestore/internal/application/interfaces"
	"github.com/matzxrr/ddd-lemonadestore/internal/domain/customer"
	"github.com/matzxrr/ddd-lemonadestore/internal/domain/order"
	"github.com/matzxrr/ddd-lemonadestore/internal/domain/store"
)

// CreateOrderCommand represents request to create an order
type CreateOrderCommand struct {
    CustomerID string
    StoreID    string
    Items      []OrderItemRequest
}

// OrderItemRequest represents item in order request
type OrderItemRequest struct {
    ProductID string
    Quantity  int
}

// CreateOrderHandler handles order creation
// WHY: Orchestrates complex order creation across multiple aggregates
type CreateOrderHandler struct {
    uow            interfaces.UnitOfWork
    eventPublisher interfaces.EventPublisher
}

func NewCreateOrderHandler(
    uow interfaces.UnitOfWork,
    eventPublisher interfaces.EventPublisher,
) *CreateOrderHandler {
    return &CreateOrderHandler{
        uow:            uow,
        eventPublisher: eventPublisher,
    }
}

// Handle creates a new order
// WHAT: Complex orchestration involving store, order, and customer aggregates
func (h *CreateOrderHandler) Handle(ctx context.Context, cmd CreateOrderCommand) (*dtos.OrderDTO, error) {
    // Start transaction
    err := h.uow.Begin(ctx)
    if err != nil {
        return nil, err
    }
    defer func() {
        if err != nil {
            h.uow.Rollback()
        }
    }()
    
    // 1. Validate customer exists and is active
    customerAgg, err := h.uow.CustomerRepository().FindByID(customer.CustomerID(cmd.CustomerID))
    if err != nil {
        return nil, err
    }
    if !customerAgg.IsActive() {
        return nil, errors.New("customer is not active")
    }
    
    // 2. Load store and validate products
    storeAgg, err := h.uow.StoreRepository().FindByID(store.StoreID(cmd.StoreID))
    if err != nil {
        return nil, err
    }
    
    // 3. Create order aggregate
    orderAgg := order.NewOrder(
        customer.CustomerID(cmd.CustomerID),
        store.StoreID(cmd.StoreID),
    )
    
    // 4. Add items and reserve inventory
    for _, item := range cmd.Items {
        // Get product details
        product, err := storeAgg.GetProduct(store.ProductID(item.ProductID))
        if err != nil {
            return nil, err
        }
        
        if !product.IsActive() {
            return nil, errors.New("product is not available")
        }
        
        // Check inventory
        available, err := storeAgg.GetAvailableQuantity(product.ID())
        if err != nil {
            return nil, err
        }
        
        if available < item.Quantity {
            return nil, errors.New("insufficient inventory for product: " + string(product.Name()))
        }
        
        // Reserve inventory
        err = storeAgg.ReserveInventory(product.ID(), item.Quantity)
        if err != nil {
            return nil, err
        }
        
        // Add to order
        err = orderAgg.AddItem(
            product.ID(),
            string(product.Name()),
            item.Quantity,
            product.Price(),
        )
        if err != nil {
            return nil, err
        }
    }
    
    // 5. Confirm order (in real app, would process payment first)
    err = orderAgg.Confirm()
    if err != nil {
        return nil, err
    }
    
    // 6. Save all changes
    err = h.uow.OrderRepository().Save(orderAgg)
    if err != nil {
        return nil, err
    }
    
    err = h.uow.StoreRepository().Save(storeAgg)
    if err != nil {
        return nil, err
    }
    
    // 7. Commit transaction
    err = h.uow.Commit()
    if err != nil {
        return nil, err
    }
    
    // 8. Publish events (after commit)
    allEvents := append(orderAgg.PullEvents(), storeAgg.PullEvents()...)
    if len(allEvents) > 0 {
        h.eventPublisher.Publish(ctx, allEvents...)
    }
    
    // 9. Return DTO
    return h.toOrderDTO(orderAgg), nil
}

// toOrderDTO converts domain order to DTO
func (h *CreateOrderHandler) toOrderDTO(orderAgg *order.Order) *dtos.OrderDTO {
    items := make([]dtos.OrderItemDTO, len(orderAgg.Items()))
    for i, item := range orderAgg.Items() {
        items[i] = dtos.OrderItemDTO{
            ID:        item.ID(),
            ProductID: string(item.ProductID()),
            Name:      item.Name(),
            Quantity:  item.Quantity(),
            UnitPrice: float64(item.UnitPrice().Amount()) / 100,
            Total:     float64(item.Total().Amount()) / 100,
        }
    }
    
    return &dtos.OrderDTO{
        ID:          string(orderAgg.ID()),
        CustomerID:  string(orderAgg.CustomerID()),
        StoreID:     string(orderAgg.StoreID()),
        Status:      string(orderAgg.Status()),
        TotalAmount: float64(orderAgg.TotalAmount().Amount()) / 100,
        Currency:    orderAgg.TotalAmount().Currency(),
        Items:       items,
        PlacedAt:    orderAgg.PlacedAt(),
    }
}
