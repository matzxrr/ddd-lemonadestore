package customer

// CustomerRepository defines persistence operations for Customer aggregate
type CustomerRepository interface {
    Save(customer *Customer) error
    FindByID(id CustomerID) (*Customer, error)
    FindByEmail(email Email) (*Customer, error)
    FindByType(customerType CustomerType) ([]*Customer, error)
}
