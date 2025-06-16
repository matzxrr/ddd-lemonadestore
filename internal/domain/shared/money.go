package shared

import (
    "errors"
    "fmt"
)

// Money is a value object representing monetary amounts
// WHY: Encapsulates money logic, prevents floating point errors, ensures currency consistency
// WHERE: Used throughout the domain wherever money is involved
type Money struct {
    amount   int64  // Store as cents to avoid floating point issues
    currency string
}

// NewMoney creates a new Money value object
// WHAT: Factory function that ensures Money is always created in a valid state
func NewMoney(cents int64, currency string) (Money, error) {
    if cents < 0 {
        return Money{}, errors.New("money amount cannot be negative")
    }
    if currency == "" {
        return Money{}, errors.New("currency is required")
    }
    return Money{amount: cents, currency: currency}, nil
}

// Add performs money addition with currency validation
// WHY: Ensures business rule that you can't add different currencies
func (m Money) Add(other Money) (Money, error) {
    if m.currency != other.currency {
        return Money{}, fmt.Errorf("cannot add different currencies: %s and %s", m.currency, other.currency)
    }
    return Money{
        amount:   m.amount + other.amount,
        currency: m.currency,
    }, nil
}

// Multiply calculates money times a quantity
// WHERE: Used in order calculations when multiplying price by quantity
func (m Money) Multiply(factor int) Money {
    return Money{
        amount:   m.amount * int64(factor),
        currency: m.currency,
    }
}

// Getters for encapsulation
func (m Money) Amount() int64    { return m.amount }
func (m Money) Currency() string { return m.currency }

// String implements Stringer for display
func (m Money) String() string {
    dollars := float64(m.amount) / 100
    return fmt.Sprintf("%.2f %s", dollars, m.currency)
}
