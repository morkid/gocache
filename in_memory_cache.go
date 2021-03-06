package gocache

import (
	"errors"
	"strings"
	"time"
)

// InMemoryCacheConfig struct
type InMemoryCacheConfig struct {
	ExpiresIn time.Duration
}

type inMemoryCache struct {
	keyValues       map[string]string
	keyRegistration map[string]time.Time
	ExpiresIn       time.Duration
}

// NewInMemoryCache new instance of inMemoryCache
func NewInMemoryCache(config InMemoryCacheConfig) *AdapterInterface {
	if config.ExpiresIn <= 0 {
		config.ExpiresIn = 3600 * time.Second
	}

	var adapter AdapterInterface = &inMemoryCache{
		ExpiresIn: config.ExpiresIn,
	}

	return &adapter
}

func (n inMemoryCache) Get(key string) (string, error) {
	if v, ok := n.keyValues[key]; ok && n.keyValues[key] != "" {
		return v, nil
	}
	return "", errors.New("Cache not found")
}

func (n *inMemoryCache) Set(key string, value string) error {
	if _, ok := n.keyValues[key]; !ok {
		n.keyValues = map[string]string{}
	}
	if _, ok := n.keyRegistration[key]; !ok {
		n.keyRegistration = map[string]time.Time{}
	}
	n.keyRegistration[key] = time.Now()
	n.keyValues[key] = value
	return nil
}

func (n inMemoryCache) IsValid(key string) bool {
	if _, ok := n.keyValues[key]; ok && n.keyValues[key] != "" {
		if _, ok := n.keyRegistration[key]; ok {
			now := time.Now()
			diff := now.Sub(n.keyRegistration[key])
			if diff > n.ExpiresIn {
				return false
			}
		}
		return true
	}
	return false
}

func (n *inMemoryCache) Clear(key string) error {
	if n.IsValid(key) {
		delete(n.keyValues, key)
		delete(n.keyRegistration, key)
	}
	return nil
}

func (n *inMemoryCache) ClearPrefix(keyPrefix string) error {
	for v := range n.keyValues {
		if strings.HasPrefix(v, keyPrefix) {
			delete(n.keyValues, v)
			delete(n.keyRegistration, v)
		}
	}
	return nil
}

func (n *inMemoryCache) ClearAll() error {
	for v := range n.keyValues {
		delete(n.keyValues, v)
		delete(n.keyRegistration, v)
	}
	return nil
}
