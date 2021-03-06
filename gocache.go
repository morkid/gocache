package gocache

// AdapterInterface interface
type AdapterInterface interface {
	// Set cache with key
	Set(key string, value string) error
	// Get cache by key
	Get(key string) (string, error)
	// IsValid check if cache is valid
	IsValid(key string) bool
	// ClearPrefix clear cache by key
	Clear(key string) error
	// ClearPrefix clear cache by key prefix
	ClearPrefix(keyPrefix string) error
	// Clear all cache
	ClearAll() error
}
