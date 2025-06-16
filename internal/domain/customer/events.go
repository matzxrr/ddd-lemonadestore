package customer

import "github.com/matzxrr/ddd-lemonadestore/internal/domain/shared"

// CustomerRegisteredEvent is raised when a new customer registers
type CustomerRegisteredEvent struct {
    shared.BaseEvent
    CustomerID string `json:"customer_id"`
    Email      string `json:"email"`
    FirstName  string `json:"first_name"`
    LastName   string `json:"last_name"`
}

func (e CustomerRegisteredEvent) EventName() string     { return "customer.registered" }
func (e CustomerRegisteredEvent) AggregateID() string   { return e.CustomerID }
func (e CustomerRegisteredEvent) AggregateType() string { return "customer" }

// CustomerContactUpdatedEvent tracks contact info changes
type CustomerContactUpdatedEvent struct {
    shared.BaseEvent
    CustomerID  string         `json:"customer_id"`
    PhoneNumber string         `json:"phone_number"`
    Address     shared.Address `json:"address"`
}

func (e CustomerContactUpdatedEvent) EventName() string     { return "customer.contact_updated" }
func (e CustomerContactUpdatedEvent) AggregateID() string   { return e.CustomerID }
func (e CustomerContactUpdatedEvent) AggregateType() string { return "customer" }

// CustomerTierUpgradedEvent celebrates customer loyalty
type CustomerTierUpgradedEvent struct {
    shared.BaseEvent
    CustomerID  string `json:"customer_id"`
    OldTier     string `json:"old_tier"`
    NewTier     string `json:"new_tier"`
    TotalPoints int    `json:"total_points"`
}

func (e CustomerTierUpgradedEvent) EventName() string     { return "customer.tier_upgraded" }
func (e CustomerTierUpgradedEvent) AggregateID() string   { return e.CustomerID }
func (e CustomerTierUpgradedEvent) AggregateType() string { return "customer" }

// PointsRedeemedEvent tracks loyalty point usage
type PointsRedeemedEvent struct {
    shared.BaseEvent
    CustomerID      string `json:"customer_id"`
    PointsRedeemed  int    `json:"points_redeemed"`
    RemainingPoints int    `json:"remaining_points"`
}

func (e PointsRedeemedEvent) EventName() string     { return "customer.points_redeemed" }
func (e PointsRedeemedEvent) AggregateID() string   { return e.CustomerID }
func (e PointsRedeemedEvent) AggregateType() string { return "customer" }

// CustomerDeactivatedEvent when customer is deactivated
type CustomerDeactivatedEvent struct {
    shared.BaseEvent
    CustomerID string `json:"customer_id"`
}

func (e CustomerDeactivatedEvent) EventName() string     { return "customer.deactivated" }
func (e CustomerDeactivatedEvent) AggregateID() string   { return e.CustomerID }
func (e CustomerDeactivatedEvent) AggregateType() string { return "customer" }
