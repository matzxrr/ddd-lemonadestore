package commands

import (
	"context"

	"github.com/matzxrr/ddd-lemonadestore/internal/application/dtos"
	"github.com/matzxrr/ddd-lemonadestore/internal/application/interfaces"
	"github.com/matzxrr/ddd-lemonadestore/internal/domain/customer"
	"github.com/matzxrr/ddd-lemonadestore/internal/domain/shared"
)

// UpdateCustomerCommand represents customer update request
type UpdateCustomerCommand struct {
    CustomerID  string
    PhoneNumber string
    Address     dtos.AddressDTO
}

// UpdateCustomerHandler handles customer updates
type UpdateCustomerHandler struct {
    customerRepo   customer.CustomerRepository
    eventPublisher interfaces.EventPublisher
}

func NewUpdateCustomerHandler(
    customerRepo customer.CustomerRepository,
    eventPublisher interfaces.EventPublisher,
) *UpdateCustomerHandler {
    return &UpdateCustomerHandler{
        customerRepo:   customerRepo,
        eventPublisher: eventPublisher,
    }
}

func (h *UpdateCustomerHandler) Handle(ctx context.Context, cmd UpdateCustomerCommand) error {
    // Load customer
    customerAgg, err := h.customerRepo.FindByID(customer.CustomerID(cmd.CustomerID))
    if err != nil {
        return err
    }
    
    // Create address value object
    address, err := shared.NewAddress(
        cmd.Address.Street,
        cmd.Address.City,
        cmd.Address.State,
        cmd.Address.ZipCode,
        cmd.Address.Country,
    )
    if err != nil {
        return err
    }
    
    // Update contact info
    err = customerAgg.UpdateContactInfo(cmd.PhoneNumber, address)
    if err != nil {
        return err
    }
    
    // Save
    err = h.customerRepo.Save(customerAgg)
    if err != nil {
        return err
    }
    
    // Publish events
    events := customerAgg.PullEvents()
    if len(events) > 0 {
        h.eventPublisher.Publish(ctx, events...)
    }
    
    return nil
}
