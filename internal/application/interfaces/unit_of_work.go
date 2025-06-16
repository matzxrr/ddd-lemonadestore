package interfaces

import (
	"context"

	"github.com/matzxrr/ddd-lemonadestore/internal/domain/customer"
	"github.com/matzxrr/ddd-lemonadestore/internal/domain/order"
	"github.com/matzxrr/ddd-lemonadestore/internal/domain/store"
)

// UnitOfWork manages transactions across repositories
// WHY: Ensures consistency when multiple aggregates need to be updated
// WHAT: Provides transactional boundary for application operations
type UnitOfWork interface {
    Begin(ctx context.Context) error
    Commit() error
    Rollback() error
    
    // Repository accessors
    StoreRepository() store.StoreRepository
    OrderRepository() order.OrderRepository
    CustomerRepository() customer.CustomerRepository
}
