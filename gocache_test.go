package gocache_test

import (
	"os"
	"testing"
	"time"

	"github.com/morkid/gocache"
)

func TestInMemory(t *testing.T) {
	config := gocache.InMemoryCacheConfig{
		ExpiresIn: 10 * time.Second,
	}

	adapter := *gocache.NewInMemoryCache(config)
	adapter.Set("foo", "bar")

	if adapter.IsValid("foo") {
		value, err := adapter.Get("foo")
		if nil != err {
			t.Error(err)
		} else if value != "bar" {
			t.Error("value not equals to bar")
		}
		adapter.Clear("foo")
	}
}

func TestDisk(t *testing.T) {
	config := gocache.DiskCacheConfig{
		Directory: os.TempDir(),
		ExpiresIn: 10 * time.Second,
	}

	adapter := *gocache.NewDiskCache(config)
	adapter.Set("foo", "bar")

	if adapter.IsValid("foo") {
		value, err := adapter.Get("foo")
		if nil != err {
			t.Error(err)
		} else if value != "bar" {
			t.Error("value not equals to bar")
		}
		adapter.Clear("foo")
	}
}
