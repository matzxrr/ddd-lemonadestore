package customer

import (
	"errors"
	"time"

	"github.com/matzxrr/ddd-lemonadestore/internal/domain/shared"
)

// Customer is the aggregate root for customer management
// WHY: Customer has its own lifecycle and consistency boundary
type Customer struct {
    shared.AggregateRoot
    id           CustomerID
    email        Email
    firstName    string
    lastName     string
    phoneNumber  PhoneNumber
    address      shared.Address
    customerType CustomerType
    loyaltyPoints int
    registeredAt time.Time
    isActive     bool
}

// NewCustomer creates a new customer
// WHERE: Called during customer registration
func NewCustomer(email string, firstName string, lastName string) (*Customer, error) {
    if firstName == "" || lastName == "" {
        return nil, errors.New("first and last name are required")
    }
    
    emailVO, err := NewEmail(email)
    if err != nil {
        return nil, err
    }
    
    customer := &Customer{
        id:            NewCustomerID(),
        email:         emailVO,
        firstName:     firstName,
        lastName:      lastName,
        customerType:  CustomerTypeRegular,
        loyaltyPoints: 0,
        registeredAt:  time.Now(),
        isActive:      true,
    }
    
    // Raise domain event
    customer.Raise(CustomerRegisteredEvent{
        BaseEvent:  shared.NewBaseEvent(),
        CustomerID: string(customer.id),
        Email:      string(emailVO),
        FirstName:  firstName,
        LastName:   lastName,
    })
    
    return customer, nil
}

// UpdateContactInfo updates customer contact details
// WHY: Contact information can change over time
func (c *Customer) UpdateContactInfo(phone string, address shared.Address) error {
    if !c.isActive {
        return errors.New("cannot update inactive customer")
    }
    
    phoneVO, err := NewPhoneNumber(phone)
    if err != nil {
        return err
    }
    
    c.phoneNumber = phoneVO
    c.address = address
    
    // Raise domain event
    c.Raise(CustomerContactUpdatedEvent{
        BaseEvent:   shared.NewBaseEvent(),
        CustomerID:  string(c.id),
        PhoneNumber: string(phoneVO),
        Address:     address,
    })
    
    return nil
}

// AddLoyaltyPoints increases customer loyalty points
// WHERE: Called when orders are completed
func (c *Customer) AddLoyaltyPoints(points int) error {
    if points <= 0 {
        return errors.New("points must be positive")
    }
    
    if !c.isActive {
        return errors.New("cannot add points to inactive customer")
    }
    
    c.loyaltyPoints += points
    
    // Check for tier upgrade
    oldType := c.customerType
    c.updateCustomerType()
    
    if oldType != c.customerType {
        c.Raise(CustomerTierUpgradedEvent{
            BaseEvent:    shared.NewBaseEvent(),
            CustomerID:   string(c.id),
            OldTier:      string(oldType),
            NewTier:      string(c.customerType),
            TotalPoints:  c.loyaltyPoints,
        })
    }
    
    return nil
}

// RedeemPoints uses loyalty points
// WHY: Customers can exchange points for discounts
func (c *Customer) RedeemPoints(points int) error {
    if points <= 0 {
        return errors.New("points must be positive")
    }
    
    if c.loyaltyPoints < points {
        return errors.New("insufficient loyalty points")
    }
    
    c.loyaltyPoints -= points
    
    c.Raise(PointsRedeemedEvent{
        BaseEvent:        shared.NewBaseEvent(),
        CustomerID:       string(c.id),
        PointsRedeemed:   points,
        RemainingPoints:  c.loyaltyPoints,
    })
    
    return nil
}

// updateCustomerType updates tier based on points
// WHAT: Business rule for customer tier progression
func (c *Customer) updateCustomerType() {
    switch {
    case c.loyaltyPoints >= 1000:
        c.customerType = CustomerTypeVIP
    case c.loyaltyPoints >= 500:
        c.customerType = CustomerTypePremium
    default:
        c.customerType = CustomerTypeRegular
    }
}

// Deactivate marks customer as inactive
func (c *Customer) Deactivate() {
    c.isActive = false
    
    c.Raise(CustomerDeactivatedEvent{
        BaseEvent:  shared.NewBaseEvent(),
        CustomerID: string(c.id),
    })
}

// GetDiscountRate returns discount based on customer type
// WHY: Different customer tiers get different benefits
func (c *Customer) GetDiscountRate() float64 {
    switch c.customerType {
    case CustomerTypeVIP:
        return 0.20 // 20% discount
    case CustomerTypePremium:
        return 0.10 // 10% discount
    default:
        return 0.0
    }
}

// Getters
func (c *Customer) ID() CustomerID           { return c.id }
func (c *Customer) Email() Email             { return c.email }
func (c *Customer) FirstName() string        { return c.firstName }
func (c *Customer) LastName() string         { return c.lastName }
func (c *Customer) FullName() string         { return c.firstName + " " + c.lastName }
func (c *Customer) PhoneNumber() PhoneNumber { return c.phoneNumber }
func (c *Customer) Address() shared.Address  { return c.address }
func (c *Customer) Type() CustomerType       { return c.customerType }
func (c *Customer) LoyaltyPoints() int       { return c.loyaltyPoints }
func (c *Customer) IsActive() bool           { return c.isActive }
