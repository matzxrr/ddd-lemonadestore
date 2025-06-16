package store

import "errors"

// Domain-specific errors
// WHY: Domain errors express business rule violations
var (
    ErrStoreNotFound       = errors.New("store not found")
    ErrProductNotFound     = errors.New("product not found")
    ErrInsufficientStock   = errors.New("insufficient stock")
    ErrInvalidPrice        = errors.New("invalid price")
    ErrDuplicateProduct    = errors.New("duplicate product name")
)
