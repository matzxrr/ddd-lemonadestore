package memory

import (
	"context"
	"errors"
	"sync"

	"github.com/matzxrr/ddd-lemonadestore/internal/application/interfaces"
	"github.com/matzxrr/ddd-lemonadestore/internal/domain/customer"
	"github.com/matzxrr/ddd-lemonadestore/internal/domain/order"
	"github.com/matzxrr/ddd-lemonadestore/internal/domain/store"
)

// InMemoryUnitOfWork implements unit of work pattern for in-memory repositories
// WHY: Ensures consistency when updating multiple aggregates
// WHAT: Provides transaction-like behavior for in-memory storage
type InMemoryUnitOfWork struct {
	mu           sync.Mutex
	inProgress   bool
	storeRepo    store.StoreRepository
	orderRepo    order.OrderRepository
	customerRepo customer.CustomerRepository
}

// NewInMemoryUnitOfWork creates a new unit of work
func NewInMemoryUnitOfWork(
	storeRepo store.StoreRepository,
	orderRepo order.OrderRepository,
	customerRepo customer.CustomerRepository,
) *InMemoryUnitOfWork {
	return &InMemoryUnitOfWork{
		storeRepo:    storeRepo,
		orderRepo:    orderRepo,
		customerRepo: customerRepo,
	}
}

// Begin starts a new unit of work
// WHERE: Called at the beginning of application commands
func (uow *InMemoryUnitOfWork) Begin(ctx context.Context) error {
	uow.mu.Lock()
	defer uow.mu.Unlock()

	if uow.inProgress {
		return errors.New("unit of work already in progress")
	}

	uow.inProgress = true
	return nil
}

// Commit commits the unit of work
// WHAT: In real implementation would commit database transaction
func (uow *InMemoryUnitOfWork) Commit() error {
	uow.mu.Lock()
	defer uow.mu.Unlock()

	if !uow.inProgress {
		return errors.New("no unit of work in progress")
	}

	uow.inProgress = false
	return nil
}

// Rollback rolls back the unit of work
func (uow *InMemoryUnitOfWork) Rollback() error {
	uow.mu.Lock()
	defer uow.mu.Unlock()

	uow.inProgress = false
	return nil
}

// Repository accessors
func (uow *InMemoryUnitOfWork) StoreRepository() store.StoreRepository {
	return uow.storeRepo
}

func (uow *InMemoryUnitOfWork) OrderRepository() order.OrderRepository {
	return uow.orderRepo
}

func (uow *InMemoryUnitOfWork) CustomerRepository() customer.CustomerRepository {
	return uow.customerRepo
}

// Ensure it implements the interface
var _ interfaces.UnitOfWork = (*InMemoryUnitOfWork)(nil)
