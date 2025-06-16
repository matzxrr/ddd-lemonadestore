package customer

import (
    "errors"
    "regexp"
    "github.com/google/uuid"
)

// CustomerID uniquely identifies a customer
type CustomerID string

func NewCustomerID() CustomerID {
    return CustomerID(uuid.New().String())
}

// Email is a value object for email addresses
// WHY: Emails require validation and are used as natural keys
type Email string

var emailRegex = regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)

func NewEmail(email string) (Email, error) {
    if !emailRegex.MatchString(email) {
        return "", errors.New("invalid email format")
    }
    return Email(email), nil
}

// PhoneNumber with validation
type PhoneNumber string

func NewPhoneNumber(phone string) (PhoneNumber, error) {
    // Simple validation - in real app would be more complex
    if len(phone) < 10 {
        return "", errors.New("invalid phone number")
    }
    return PhoneNumber(phone), nil
}

// CustomerType represents different customer categories
type CustomerType string

const (
    CustomerTypeRegular  CustomerType = "REGULAR"
    CustomerTypePremium  CustomerType = "PREMIUM"
    CustomerTypeVIP      CustomerType = "VIP"
)
