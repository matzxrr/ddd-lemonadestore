package shared

import "errors"

// Address is a value object representing a physical address
// WHY: Addresses have no identity, two identical addresses are the same
// WHAT: Immutable object that validates address components
type Address struct {
    street  string
    city    string
    state   string
    zipCode string
    country string
}

// NewAddress creates a new Address with validation
// WHERE: Used when creating customer addresses or store locations
func NewAddress(street, city, state, zipCode, country string) (Address, error) {
    if street == "" || city == "" || state == "" || zipCode == "" || country == "" {
        return Address{}, errors.New("all address fields are required")
    }
    
    return Address{
        street:  street,
        city:    city,
        state:   state,
        zipCode: zipCode,
        country: country,
    }, nil
}

// Getters for encapsulation
func (a Address) Street() string  { return a.street }
func (a Address) City() string    { return a.city }
func (a Address) State() string   { return a.state }
func (a Address) ZipCode() string { return a.zipCode }
func (a Address) Country() string { return a.country }
