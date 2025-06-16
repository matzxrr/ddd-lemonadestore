package commands

import (
	"context"
	"errors"

	"github.com/matzxrr/ddd-lemonadestore/internal/application/dtos"
	"github.com/matzxrr/ddd-lemonadestore/internal/application/interfaces"
	"github.com/matzxrr/ddd-lemonadestore/internal/domain/customer"
)

// RegisterCustomerCommand represents customer registration request
type RegisterCustomerCommand struct {
    Email     string
    FirstName string
    LastName  string
}

// RegisterCustomerHandler handles customer registration
type RegisterCustomerHandler struct {
    customerRepo   customer.CustomerRepository
    eventPublisher interfaces.EventPublisher
}

func NewRegisterCustomerHandler(
    customerRepo customer.CustomerRepository,
    eventPublisher interfaces.EventPublisher,
) *RegisterCustomerHandler {
    return &RegisterCustomerHandler{
        customerRepo:   customerRepo,
        eventPublisher: eventPublisher,
    }
}

func (h *RegisterCustomerHandler) Handle(ctx context.Context, cmd RegisterCustomerCommand) (*dtos.CustomerDTO, error) {
    // Check if email already exists
    existingCustomer, _ := h.customerRepo.FindByEmail(customer.Email(cmd.Email))
    if existingCustomer != nil {
        return nil, errors.New("email already registered")
    }
    
    // Create customer
    customerAgg, err := customer.NewCustomer(cmd.Email, cmd.FirstName, cmd.LastName)
    if err != nil {
        return nil, err
    }
    
    // Save
    err = h.customerRepo.Save(customerAgg)
    if err != nil {
        return nil, err
    }
    
    // Publish events
    events := customerAgg.PullEvents()
    if len(events) > 0 {
        h.eventPublisher.Publish(ctx, events...)
    }
    
    // Return DTO
    return &dtos.CustomerDTO{
        ID:            string(customerAgg.ID()),
        Email:         string(customerAgg.Email()),
        FirstName:     customerAgg.FirstName(),
        LastName:      customerAgg.LastName(),
        Type:          string(customerAgg.Type()),
        LoyaltyPoints: customerAgg.LoyaltyPoints(),
        RegisteredAt:  customerAgg.RegisteredAt,
        IsActive:      customerAgg.IsActive(),
    }, nil
}
