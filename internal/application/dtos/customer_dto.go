package dtos

import "time"

// CustomerDTO represents customer data for application layer
type CustomerDTO struct {
    ID            string    `json:"id"`
    Email         string    `json:"email"`
    FirstName     string    `json:"first_name"`
    LastName      string    `json:"last_name"`
    PhoneNumber   string    `json:"phone_number,omitempty"`
    Type          string    `json:"type"`
    LoyaltyPoints int       `json:"loyalty_points"`
    RegisteredAt  time.Time `json:"registered_at"`
    IsActive      bool      `json:"is_active"`
}

// AddressDTO represents address data
type AddressDTO struct {
    Street  string `json:"street"`
    City    string `json:"city"`
    State   string `json:"state"`
    ZipCode string `json:"zip_code"`
    Country string `json:"country"`
}
