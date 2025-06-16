package order

import "github.com/matzxrr/ddd-lemonadestore/internal/domain/shared"

// OrderPolicy defines business policies for orders
// WHY: Centralizes business rules that involve multiple factors
type OrderPolicy interface {
    CanBeCancelled(order *Order) bool
    GetPreparationTime(order *Order) int // minutes
}

// StandardOrderPolicy implements default business policies
type StandardOrderPolicy struct{}

// CanBeCancelled determines if order can be cancelled
// WHAT: Business rule - orders can only be cancelled before preparation starts
func (p *StandardOrderPolicy) CanBeCancelled(order *Order) bool {
    return order.Status() == OrderStatusPending || 
           order.Status() == OrderStatusConfirmed
}

// GetPreparationTime estimates preparation time based on order complexity
func (p *StandardOrderPolicy) GetPreparationTime(order *Order) int {
    baseTime := 5 // 5 minutes base
    itemTime := len(order.Items()) * 2 // 2 minutes per item
    
    // Large orders take longer
    largeOrderThreshold, _ := shared.NewMoney(5000, "USD") // $50
    largeOrderSpec := NewLargeOrderSpec(largeOrderThreshold)
    if largeOrderSpec.IsSatisfiedBy(order) {
        itemTime += 10 // Extra 10 minutes for large orders
    }
    
    return baseTime + itemTime
}
