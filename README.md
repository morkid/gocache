# gocache - simple key value cache adapter for golang
[![Go](https://github.com/morkid/gocache/actions/workflows/go.yml/badge.svg)](https://github.com/morkid/gocache/actions/workflows/go.yml)
[![Build Status](https://travis-ci.com/morkid/gocache.svg?branch=master)](https://travis-ci.com/morkid/gocache)
[![Go Report Card](https://goreportcard.com/badge/github.com/morkid/gocache)](https://goreportcard.com/report/github.com/morkid/gocache)

## Installation
```bash
go get -d github.com/morkid/gocache
```

## In Memory cache example
```go
package main
import (
    "time"
    "fmt"
    "github.com/morkid/gocache"
)

func main() {
    config := gocache.InMemoryCacheConfig{
        ExpiresIn: 10 * time.Second,
    }

    adapter := *gocache.NewInMemoryCache(config)
    adapter.Set("foo", "bar")

    if adapter.IsValid("foo") {
        value, err := adapter.Get("foo")
        if nil != err {
            fmt.Println(err.Error())
        } else if value != "bar" {
            fmt.Println("value not equals to bar")
        }
        adapter.Clear("foo")
    }
}
```

## Disk cache example
```go
package main
import (
    "os"
    "time"
    "fmt"
    "github.com/morkid/gocache"
)

func main() {
    config := gocache.DiskCacheConfig{
        Directory: os.TempDir(),
        ExpiresIn: 10 * time.Second,
    }

    adapter := *gocache.NewDiskCache(config)
    adapter.Set("foo", "bar")

    if adapter.IsValid("foo") {
        value, err := adapter.Get("foo")
        if nil != err {
            fmt.Println(err.Error())
        } else if value != "bar" {
            fmt.Println("value not equals to bar")
        }
        adapter.Clear("foo")
    }
}
```
## Custom cache adapter
You can create your custom cache adapter by implementing the `AdapterInterface`:

```go
type AdapterInterface interface {

	Set(key string, value string) error

	Get(key string) (string, error)
	
    IsValid(key string) bool
	
    Clear(key string) error
	
    ClearPrefix(keyPrefix string) error
	
    ClearAll() error
}
```


## License

Published under the [MIT License](https://github.com/morkid/paginate/blob/master/LICENSE).