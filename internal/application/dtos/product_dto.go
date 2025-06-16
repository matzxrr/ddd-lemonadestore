package dtos

// ProductDTO represents product data for application layer
// WHY: Decouples domain objects from external representation
// WHERE: Used in application services to transfer data between layers
type ProductDTO struct {
    ID          string  `json:"id"`
    Name        string  `json:"name"`
    Description string  `json:"description"`
    Price       float64 `json:"price"`
    Currency    string  `json:"currency"`
    IsActive    bool    `json:"is_active"`
    Quantity    int     `json:"quantity"`
}
