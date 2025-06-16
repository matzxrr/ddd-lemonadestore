package order

import "github.com/lemonade-store/internal/domain/shared"

// Specification pattern for complex business rules
// WHY: Encapsulates business rules in reusable, composable units
type Specification interface {
    IsSatisfiedBy(candidate interface{}) bool
}

// LargeOrderSpec identifies orders above a threshold
// WHERE: Used for applying discounts or special handling
type LargeOrderSpec struct {
    minAmount shared.Money
}

func NewLargeOrderSpec(minAmount shared.Money) *LargeOrderSpec {
    return &LargeOrderSpec{minAmount: minAmount}
}

func (s *LargeOrderSpec) IsSatisfiedBy(candidate interface{}) bool {
    order, ok := candidate.(*Order)
    if !ok {
        return false
    }
    return order.TotalAmount().Amount() >= s.minAmount.Amount()
}

// RushOrderSpec identifies orders needing quick preparation
type RushOrderSpec struct {
    maxItems int
}

func NewRushOrderSpec(maxItems int) *RushOrderSpec {
    return &RushOrderSpec{maxItems: maxItems}
}

func (s *RushOrderSpec) IsSatisfiedBy(candidate interface{}) bool {
    order, ok := candidate.(*Order)
    if !ok {
        return false
    }
    return len(order.Items()) <= s.maxItems
}

// AndSpec combines two specifications with AND logic
type AndSpec struct {
    spec1 Specification
    spec2 Specification
}

func (s *AndSpec) IsSatisfiedBy(candidate interface{}) bool {
    return s.spec1.IsSatisfiedBy(candidate) && s.spec2.IsSatisfiedBy(candidate)
}
