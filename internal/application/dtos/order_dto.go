package dtos

import "time"

// OrderDTO represents order data for application layer
type OrderDTO struct {
    ID          string         `json:"id"`
    CustomerID  string         `json:"customer_id"`
    StoreID     string         `json:"store_id"`
    Status      string         `json:"status"`
    TotalAmount float64        `json:"total_amount"`
    Currency    string         `json:"currency"`
    Items       []OrderItemDTO `json:"items"`
    PlacedAt    time.Time      `json:"placed_at"`
}

// OrderItemDTO represents order item data
type OrderItemDTO struct {
    ID        string  `json:"id"`
    ProductID string  `json:"product_id"`
    Name      string  `json:"name"`
    Quantity  int     `json:"quantity"`
    UnitPrice float64 `json:"unit_price"`
    Total     float64 `json:"total"`
}
