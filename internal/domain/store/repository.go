package store

// StoreRepository defines persistence operations for Store aggregate
// WHY: Domain defines the interface, infrastructure implements it
// WHERE: Used by application layer to persist and retrieve stores
type StoreRepository interface {
    Save(store *Store) error
    FindByID(id StoreID) (*Store, error)
    FindAll() ([]*Store, error)
}
