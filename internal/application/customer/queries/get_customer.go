package queries

import (
	"context"

	"github.com/matzxrr/ddd-lemonadestore/internal/application/dtos"
	"github.com/matzxrr/ddd-lemonadestore/internal/domain/customer"
)

// GetCustomerQuery represents customer query request
type GetCustomerQuery struct {
    CustomerID string
}

// GetCustomerHandler handles customer queries
type GetCustomerHandler struct {
    customerRepo customer.CustomerRepository
}

func NewGetCustomerHandler(customerRepo customer.CustomerRepository) *GetCustomerHandler {
    return &GetCustomerHandler{customerRepo: customerRepo}
}

func (h *GetCustomerHandler) Handle(ctx context.Context, query GetCustomerQuery) (*dtos.CustomerDTO, error) {
    // Load customer
    customerAgg, err := h.customerRepo.FindByID(customer.CustomerID(query.CustomerID))
    if err != nil {
        return nil, err
    }
    
    // Convert to DTO
    return &dtos.CustomerDTO{
        ID:            string(customerAgg.ID()),
        Email:         string(customerAgg.Email()),
        FirstName:     customerAgg.FirstName(),
        LastName:      customerAgg.LastName(),
        PhoneNumber:   string(customerAgg.PhoneNumber()),
        Type:          string(customerAgg.Type()),
        LoyaltyPoints: customerAgg.LoyaltyPoints(),
        IsActive:      customerAgg.IsActive(),
    }, nil
}
